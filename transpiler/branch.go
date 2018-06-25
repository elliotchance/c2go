// This file contains functions for transpiling common branching and control
// flow, such as "if", "while", "do" and "for". The more complicated control
// flows like "switch" will be put into their own file of the same or sensible
// name.

package transpiler

import (
	"fmt"
	"go/token"

	goast "go/ast"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"
	"golang.org/x/tools/go/ast/astutil"
)

func transpileIfStmt(n *ast.IfStmt, p *program.Program) (
	*goast.IfStmt, []goast.Stmt, []goast.Stmt, error) {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}
	children := n.Children()

	// There is always 4 or 5 children in an IfStmt. For example:
	//
	//     if (i == 0) {
	//         return 0;
	//     } else {
	//         return 1;
	//     }
	//
	// 1. Not sure what this is for. This gets removed.
	// 2. Not sure what this is for.
	// 3. conditional = BinaryOperator: i == 0
	// 4. body = CompoundStmt: { return 0; }
	// 5. elseBody = CompoundStmt: { return 1; }
	//
	// elseBody will be nil if there is no else clause.

	// On linux I have seen only 4 children for an IfStmt with the same
	// definitions above, but missing the first argument. Since we don't
	// know what the first argument is for anyway we will just remove it on
	// Mac if necessary.
	if len(children) == 5 && children[0] != nil {
		panic("non-nil child 0 in IfStmt")
	}
	if len(children) == 5 {
		children = children[1:]
	}

	// From here on there must be 4 children.
	if len(children) != 4 {
		panic(fmt.Sprintf("Expected 4 children in IfStmt, got %#v", children))
	}

	// Maybe we will discover what the nil value is?
	if children[0] != nil {
		panic("non-nil child 0 in IfStmt")
	}

	// The last parameter must be false because we are transpiling an
	// expression - assignment operators need to be wrapped in closures.
	conditional, conditionalType, newPre, newPost, err := transpileToExpr(children[1], p, false)
	if err != nil {
		return nil, nil, nil, err
	}
	// null in C is false
	if conditionalType == types.NullPointer {
		conditional = util.NewIdent("false")
		conditionalType = "bool"
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	// The condition in Go must always be a bool.
	boolCondition, err := types.CastExpr(p, conditional, conditionalType, "bool")
	p.AddMessage(p.GenerateWarningOrErrorMessage(err, n, boolCondition == nil))

	if boolCondition == nil {
		boolCondition = util.NewNil()
	}

	body, newPre, newPost, err := transpileToBlockStmt(children[2], p)
	if err != nil {
		return nil, nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	if body == nil {
		return nil, nil, nil, fmt.Errorf("Body of If cannot by nil")
	}
	if boolCondition == nil {
		return nil, nil, nil, fmt.Errorf("Bool Condition in If cannot by nil")
	}
	r := &goast.IfStmt{
		Cond: boolCondition,
		Body: body,
	}

	if children[3] != nil {
		elseBody, newPre, newPost, err := transpileToBlockStmt(children[3], p)
		if err != nil {
			return nil, nil, nil, err
		}

		preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

		if elseBody != nil {
			r.Else = elseBody
			if _, ok := children[3].(*ast.IfStmt); ok {
				if len(elseBody.List) == 1 {
					r.Else = elseBody.List[0]
				}
			}
		} else {
			return nil, nil, nil, fmt.Errorf("Body of Else in If cannot be nil")
		}
	}

	return r, preStmts, postStmts, nil
}

func transpileForStmt(n *ast.ForStmt, p *program.Program) (
	f goast.Stmt, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {

	// This `defer` is workaround
	// Please remove after solving all problems
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot tranpile ForStmt: err = %v", err)
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}()

	children := n.Children()

	// There are always 5 children in a ForStmt, for example:
	//
	//     for ( c = 0 ; c < n ; c++ ) {
	//         doSomething();
	//     }
	//
	// 1. initExpression = BinaryStmt: c = 0
	// 2. Not sure what this is for, but it's always nil. There is a panic
	//    below in case we discover what it is used for (pun intended).
	// 3. conditionalExpression = BinaryStmt: c < n
	// 4. stepExpression = BinaryStmt: c++
	// 5. body = CompoundStmt: { CallExpr }

	if len(children) != 5 {
		panic(fmt.Sprintf("Expected 5 children in ForStmt, got %#v", children))
	}

	// TODO: The second child of a ForStmt appears to always be null.
	// Are there any cases where it is used?
	if children[1] != nil {
		panic("non-nil child 1 in ForStmt")
	}

	switch c := children[0].(type) {
	case *ast.BinaryOperator:
		if c.Operator == "," {
			// If we have 2 and more initializations like
			// in operator for
			// for( a = 0, b = 0, c = 0; a < 5; a ++)
			// recursive action to code like that:
			// a = 0;
			// b = 0;
			// for(c = 0 ; a < 5 ; a++)
			before, newPre, newPost, err := transpileToStmt(children[0], p)
			if err != nil {
				return nil, nil, nil, err
			}
			preStmts = append(preStmts, combineStmts(before, newPre, newPost)...)
			children[0] = c.Children()[1]
		}
	case *ast.DeclStmt:
		{
			// If we have 2 and more initializations like
			// in operator for
			// for(int a = 0, b = 0, c = 0; a < 5; a ++)
			newPre, err := transpileToStmts(children[0], p)
			if err != nil {
				return nil, nil, nil, err
			}
			children[0] = nil
			preStmts = append(preStmts, newPre...)
		}
	}

	init, newPre, newPost, err := transpileToStmt(children[0], p)
	if err != nil {
		return nil, nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	// If we have 2 and more increments
	// in operator for
	// for( a = 0; a < 5; a ++, b++, c+=2)
	switch c := children[3].(type) {
	case *ast.BinaryOperator:
		if c.Operator == "," {
			// recursive action to code like that:
			// a = 0;
			// b = 0;
			// for(a = 0 ; a < 5 ; ){
			// 		body
			// 		a++;
			// 		b++;
			//		c+=2;
			// }
			//
			var compound *ast.CompoundStmt
			if children[4] != nil {
				// if body is exist
				if _, ok := children[4].(*ast.CompoundStmt); !ok {
					compound = new(ast.CompoundStmt)
					compound.AddChild(children[4])
				} else {
					compound = children[4].(*ast.CompoundStmt)
				}
			} else {
				// if body is not exist
				compound = new(ast.CompoundStmt)
			}
			compound.ChildNodes = append(compound.Children(), c.Children()[0:len(c.Children())]...)
			children[4] = compound
			children[3] = nil
		}
	}

	var post goast.Stmt
	var transpilate bool
	if v, ok := children[3].(*ast.UnaryOperator); ok {
		if vv, ok := v.Children()[0].(*ast.DeclRefExpr); ok {
			if !types.IsPointer(p, vv.Type) && !types.IsFunction(vv.Type) {
				switch v.Operator {
				case "++":
					// for case:
					// for(...;...;i++)...
					post = &goast.IncDecStmt{
						X:   util.NewIdent(vv.Name),
						Tok: token.INC,
					}
					transpilate = true
				case "--":
					// for case:
					// for(...;...;i--)...
					post = &goast.IncDecStmt{
						X:   util.NewIdent(vv.Name),
						Tok: token.DEC,
					}
					transpilate = true
				}
			}
		}
	}
	if !transpilate {
		post, newPre, newPost, err = transpileToStmt(children[3], p)
		if err != nil {
			return nil, nil, nil, err
		}

		preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)
	}

	// If we have 2 and more conditions
	// in operator for
	// for( a = 0; b = c, b++, a < 5; a ++)
	switch c := children[2].(type) {
	case *ast.BinaryOperator:
		if c.Operator == "," {
			// recursive action to code like that:
			// a = 0;
			// b = 0;
			// for(a = 0 ; ; c+=2){
			// 		b = c;
			// 		b++;
			//		if (!(a < 5))
			// 			break;
			// 		body
			// }
			tempSlice := c.Children()[0 : len(c.Children())-1]

			var condition ast.IfStmt
			condition.AddChild(nil)
			var par ast.ParenExpr
			par.Type = "bool"
			par.AddChild(c.Children()[len(c.Children())-1])
			var unitary ast.UnaryOperator
			unitary.Type = "bool"
			unitary.AddChild(&par)
			unitary.Operator = "!"
			condition.AddChild(&unitary)
			var c ast.CompoundStmt
			c.AddChild(&ast.BreakStmt{})
			condition.AddChild(&c)
			condition.AddChild(nil)

			tempSlice = append(tempSlice, &condition)

			var compound *ast.CompoundStmt
			if children[4] != nil {
				// if body is exist
				compound = children[4].(*ast.CompoundStmt)
			} else {
				// if body is not exist
				compound = new(ast.CompoundStmt)
			}
			compound.ChildNodes = append(tempSlice, compound.Children()...)
			children[4] = compound
			children[2] = nil
		}
	}

	// The condition can be nil. This means an infinite loop and will be
	// rendered in Go as "for {".
	var condition goast.Expr
	if children[2] != nil {
		var conditionType string
		var newPre, newPost []goast.Stmt

		// The last parameter must be false because we are transpiling an
		// expression - assignment operators need to be wrapped in closures.
		condition, conditionType, newPre, newPost, err = atomicOperation(children[2], p)
		if err != nil {
			return nil, nil, nil, err
		}
		// null in C is false
		if conditionType == types.NullPointer {
			condition = util.NewIdent("false")
			conditionType = "bool"
		}

		preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

		condition, err = types.CastExpr(p, condition, conditionType, "bool")
		p.AddMessage(p.GenerateWarningOrErrorMessage(err, n, condition == nil))

		if condition == nil {
			condition = util.NewNil()
		}
	}

	body, newPre, newPost, err := transpileToBlockStmt(children[4], p)
	if err != nil {
		return nil, nil, nil, err
	}
	if body == nil {
		return nil, nil, nil, fmt.Errorf("Body of For cannot be nil")
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	// avoid extra block around FOR
	if len(preStmts) == 0 && len(postStmts) == 0 {
		return &goast.ForStmt{
			Init: init,
			Cond: condition,
			Post: post,
			Body: body,
		}, preStmts, postStmts, nil
	}

	// for avoid dublication of init values for
	// case with 2 for`s
	var block goast.BlockStmt
	var forStmt = goast.ForStmt{
		Init: init,
		Cond: condition,
		Post: post,
		Body: body,
	}
	block.List = combineStmts(&forStmt, preStmts, postStmts)
	block.Lbrace = 1

	return &block, nil, nil, nil
}

// transpileWhileStmt - transpiler for operator While.
// We have only operator FOR in Go, but in C we also have
// operator WHILE. So, we have to convert to operator FOR.
// We choose directly conversion  from AST C code to AST C code, for
// - avoid dublicate of code in realization WHILE and FOR.
// - create only one operator FOR powerful.
// Example of C code with operator WHILE:
//	while(i > 0){
//		printf("While: %d\n",i);
//		i--;
//	}
// AST for that code:
//    |-WhileStmt 0x2530a10 <line:6:2, line:9:2>
//    | |-<<<NULL>>>
//    | |-BinaryOperator 0x25307f0 <line:6:8, col:12> 'int' '>'
//    | | |-ImplicitCastExpr 0x25307d8 <col:8> 'int' <LValueToRValue>
//    | | | `-DeclRefExpr 0x2530790 <col:8> 'int' lvalue Var 0x25306f8 'i' 'int'
//    | | `-IntegerLiteral 0x25307b8 <col:12> 'int' 0
//    | `-CompoundStmt 0x25309e8 <col:14, line:9:2>
//    |   |-CallExpr 0x2530920 <line:7:3, col:25> 'int'
//    |   | |-ImplicitCastExpr 0x2530908 <col:3> 'int (*)(const char *, ...)' <FunctionToPointerDecay>
//    |   | | `-DeclRefExpr 0x2530818 <col:3> 'int (const char *, ...)' Function 0x2523ee8 'printf' 'int (const char *, ...)'
//    |   | |-ImplicitCastExpr 0x2530970 <col:10> 'const char *' <BitCast>
//    |   | | `-ImplicitCastExpr 0x2530958 <col:10> 'char *' <ArrayToPointerDecay>
//    |   | |   `-StringLiteral 0x2530878 <col:10> 'char [11]' lvalue "While: %d\n"
//    |   | `-ImplicitCastExpr 0x2530988 <col:24> 'int' <LValueToRValue>
//    |   |   `-DeclRefExpr 0x25308b0 <col:24> 'int' lvalue Var 0x25306f8 'i' 'int'
//    |   `-UnaryOperator 0x25309c8 <line:8:3, col:4> 'int' postfix '--'
//    |     `-DeclRefExpr 0x25309a0 <col:3> 'int' lvalue Var 0x25306f8 'i' 'int'
//
// Example of C code with operator FOR:
//	for (;i > 0;){
//		printf("For: %d\n",i);
//		i--;
//	}
// AST for that code:
//    |-ForStmt 0x2530d08 <line:11:2, line:14:2>
//    | |-<<<NULL>>>
//    | |-<<<NULL>>>
//    | |-BinaryOperator 0x2530b00 <line:11:8, col:12> 'int' '>'
//    | | |-ImplicitCastExpr 0x2530ae8 <col:8> 'int' <LValueToRValue>
//    | | | `-DeclRefExpr 0x2530aa0 <col:8> 'int' lvalue Var 0x25306f8 'i' 'int'
//    | | `-IntegerLiteral 0x2530ac8 <col:12> 'int' 0
//    | |-<<<NULL>>>
//    | `-CompoundStmt 0x2530ce0 <col:15, line:14:2>
//    |   |-CallExpr 0x2530bf8 <line:12:3, col:23> 'int'
//    |   | |-ImplicitCastExpr 0x2530be0 <col:3> 'int (*)(const char *, ...)' <FunctionToPointerDecay>
//    |   | | `-DeclRefExpr 0x2530b28 <col:3> 'int (const char *, ...)' Function 0x2523ee8 'printf' 'int (const char *, ...)'
//    |   | |-ImplicitCastExpr 0x2530c48 <col:10> 'const char *' <BitCast>
//    |   | | `-ImplicitCastExpr 0x2530c30 <col:10> 'char *' <ArrayToPointerDecay>
//    |   | |   `-StringLiteral 0x2530b88 <col:10> 'char [9]' lvalue "For: %d\n"
//    |   | `-ImplicitCastExpr 0x2530c60 <col:22> 'int' <LValueToRValue>
//    |   |   `-DeclRefExpr 0x2530bb8 <col:22> 'int' lvalue Var 0x25306f8 'i' 'int'
//    |   `-UnaryOperator 0x2530ca0 <line:13:3, col:4> 'int' postfix '--'
//    |     `-DeclRefExpr 0x2530c78 <col:3> 'int' lvalue Var 0x25306f8 'i' 'int'
func transpileWhileStmt(n *ast.WhileStmt, p *program.Program) (
	goast.Stmt, []goast.Stmt, []goast.Stmt, error) {
	var forOperator ast.ForStmt
	forOperator.AddChild(nil)
	forOperator.AddChild(nil)
	forOperator.AddChild(n.Children()[1])
	forOperator.AddChild(nil)
	if n.Children()[2] == nil {
		// added for case if WHILE haven't body, for example:
		// while(0);
		n.Children()[2] = &ast.CompoundStmt{}
	}
	forOperator.AddChild(n.Children()[2])

	return transpileForStmt(&forOperator, p)
}

// transpileDoStmt - transpiler for operator Do...While
// We have only operators FOR and IF in Go, but in C we also have
// operator DO...WHILE. So, we have to convert to operators FOR and IF.
// We choose directly conversion  from AST C code to AST C code, for:
// - avoid dublicate of code in realization DO...WHILE and FOR.
// - create only one powerful operator FOR.
// Example of C code with operator DO...WHILE:
//	do{
//		printf("While: %d\n",i);
//		i--;
//	}while(i > 0);
// AST for that code:
//    |-DoStmt 0x3bb1a68 <line:7:2, line:10:14>
//    | |-CompoundStmt 0x3bb19b8 <line:7:4, line:10:2>
//    | | |-CallExpr 0x3bb18f0 <line:8:3, col:25> 'int'
//    | | | |-ImplicitCastExpr 0x3bb18d8 <col:3> 'int (*)(const char *, ...)' <FunctionToPointerDecay>
//    | | | | `-DeclRefExpr 0x3bb17e0 <col:3> 'int (const char *, ...)' Function 0x3ba4ee8 'printf' 'int (const char *, ...)'
//    | | | |-ImplicitCastExpr 0x3bb1940 <col:10> 'const char *' <BitCast>
//    | | | | `-ImplicitCastExpr 0x3bb1928 <col:10> 'char *' <ArrayToPointerDecay>
//    | | | |   `-StringLiteral 0x3bb1848 <col:10> 'char [11]' lvalue "While: %d\n"
//    | | | `-ImplicitCastExpr 0x3bb1958 <col:24> 'int' <LValueToRValue>
//    | | |   `-DeclRefExpr 0x3bb1880 <col:24> 'int' lvalue Var 0x3bb16f8 'i' 'int'
//    | | `-UnaryOperator 0x3bb1998 <line:9:3, col:4> 'int' postfix '--'
//    | |   `-DeclRefExpr 0x3bb1970 <col:3> 'int' lvalue Var 0x3bb16f8 'i' 'int'
//    | `-BinaryOperator 0x3bb1a40 <line:10:9, col:13> 'int' '>'
//    |   |-ImplicitCastExpr 0x3bb1a28 <col:9> 'int' <LValueToRValue>
//    |   | `-DeclRefExpr 0x3bb19e0 <col:9> 'int' lvalue Var 0x3bb16f8 'i' 'int'
//    |   `-IntegerLiteral 0x3bb1a08 <col:13> 'int' 0
//
// Example of C code with operator FOR:
//	for(;;){
//		printf("For: %d\n",i);
//		i--;
//		if(!(i>0)){
//			break;
//		}
//	}
// AST for that code:
//    |-ForStmt 0x3bb1e08 <line:12:2, line:18:2>
//    | |-<<<NULL>>>
//    | |-<<<NULL>>>
//    | |-<<<NULL>>>
//    | |-<<<NULL>>>
//    | `-CompoundStmt 0x3bb1dd8 <line:12:9, line:18:2>
//    |   |-CallExpr 0x3bb1bc8 <line:13:3, col:23> 'int'
//    |   | |-ImplicitCastExpr 0x3bb1bb0 <col:3> 'int (*)(const char *, ...)' <FunctionToPointerDecay>
//    |   | | `-DeclRefExpr 0x3bb1af8 <col:3> 'int (const char *, ...)' Function 0x3ba4ee8 'printf' 'int (const char *, ...)'
//    |   | |-ImplicitCastExpr 0x3bb1c18 <col:10> 'const char *' <BitCast>
//    |   | | `-ImplicitCastExpr 0x3bb1c00 <col:10> 'char *' <ArrayToPointerDecay>
//    |   | |   `-StringLiteral 0x3bb1b58 <col:10> 'char [9]' lvalue "For: %d\n"
//    |   | `-ImplicitCastExpr 0x3bb1c30 <col:22> 'int' <LValueToRValue>
//    |   |   `-DeclRefExpr 0x3bb1b88 <col:22> 'int' lvalue Var 0x3bb16f8 'i' 'int'
//    |   |-UnaryOperator 0x3bb1c70 <line:14:3, col:4> 'int' postfix '--'
//    |   | `-DeclRefExpr 0x3bb1c48 <col:3> 'int' lvalue Var 0x3bb16f8 'i' 'int'
//    |   `-IfStmt 0x3bb1da8 <line:15:3, line:17:3>
//    |     |-<<<NULL>>>
//    |     |-UnaryOperator 0x3bb1d60 <line:15:6, col:11> 'int' prefix '!'
//    |     | `-ParenExpr 0x3bb1d40 <col:7, col:11> 'int'
//    |     |   `-BinaryOperator 0x3bb1d18 <col:8, col:10> 'int' '>'
//    |     |     |-ImplicitCastExpr 0x3bb1d00 <col:8> 'int' <LValueToRValue>
//    |     |     | `-DeclRefExpr 0x3bb1c90 <col:8> 'int' lvalue Var 0x3bb16f8 'i' 'int'
//    |     |     `-IntegerLiteral 0x3bb1ce0 <col:10> 'int' 0
//    |     |-CompoundStmt 0x3bb1d88 <col:13, line:17:3>
//    |     | `-BreakStmt 0x3bb1d80 <line:16:4>
//    |     `-<<<NULL>>>
func transpileDoStmt(n *ast.DoStmt, p *program.Program) (
	f goast.Stmt, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	var forOperator ast.ForStmt
	forOperator.AddChild(nil)
	forOperator.AddChild(nil)
	forOperator.AddChild(nil)
	forOperator.AddChild(nil)
	var c *ast.CompoundStmt
	if _, ok := n.Children()[0].(*ast.CompoundStmt); ok {
		c = n.Children()[0].(*ast.CompoundStmt)
	} else {
		var newComp ast.CompoundStmt
		newComp.AddChild(n.Children()[0])
		c = &newComp
	}
	ifBreak := createIfWithNotConditionAndBreak(n.Children()[1])
	c.AddChild(&ifBreak)
	forOperator.AddChild(c)
	f, preStmts, postStmts, err = transpileForStmt(&forOperator, p)

	if hasContinueStmt(f) {
		adaptContinueStmt(f, p)
	}
	return
}

type continueDetector struct {
	level       int
	forLevel    []int
	hasContinue bool
}

func (c *continueDetector) Visit(node goast.Node) goast.Visitor {
	if node == nil {
		if len(c.forLevel) > 0 && c.level == c.forLevel[len(c.forLevel)-1] {
			// we are exiting a for loop
			c.forLevel = c.forLevel[:len(c.forLevel)-1]
		}
		c.level--
		return nil
	}
	c.level++
	switch n := node.(type) {
	case *goast.ForStmt:
		if len(c.forLevel) > 0 {
			// Do not look for continue within an inner for loop
			c.level--
			return nil
		}
		// we have found the outer for loop
		c.forLevel = append(c.forLevel, c.level)
	case *goast.RangeStmt:
		// Do not look for continue within the children of this AST node type,
		// since a continue would refer to it and not the outer for loop.
		c.level--
		return nil
	case *goast.BranchStmt:
		// only look for continue within the outer for loop
		if n.Tok == token.CONTINUE && len(c.forLevel) > 0 {
			c.hasContinue = true
		}
	}
	return c
}

func hasContinueStmt(e goast.Stmt) bool {
	var c continueDetector
	goast.Walk(&c, e)
	return c.hasContinue
}

func adaptContinueStmt(e goast.Stmt, p *program.Program) {
	var (
		level                int
		forLevel             []int
		beforeConditionLabel string
	)
	funcTransformBreak := func(cursor *astutil.Cursor) bool {
		level++
		node := cursor.Node()
		switch n := node.(type) {
		case *goast.BranchStmt:
			// only replace continue within the outer for loop
			if n.Tok == token.CONTINUE && len(forLevel) > 0 {
				if beforeConditionLabel == "" {
					panic("Before condition label not set")
				}
				cursor.Replace(&goast.BranchStmt{
					Label: util.NewIdent(beforeConditionLabel),
					Tok:   token.GOTO,
				})
			}
		case *goast.ForStmt:
			if len(forLevel) > 0 {
				// Do not look for continue within an inner for loop
				level--
				return false
			}
			// we have found the outer for loop
			forLevel = append(forLevel, level)
			beforeConditionLabel = setConditionLabel(n, p)
			return true
		case *goast.RangeStmt:
			// Do not look for continue within the children of this AST node type,
			// since a continue would refer to it and not the outer for loop.
			level--
			return false
		}
		return true
	}
	postFunc := func(cursor *astutil.Cursor) bool {
		if len(forLevel) > 0 && level == forLevel[len(forLevel)-1] {
			// we are exiting a for loop
			forLevel = forLevel[:len(forLevel)-1]
		}
		level--
		return true
	}
	// Replace continue's with goto doWhileCondLabel
	astutil.Apply(e, funcTransformBreak, postFunc)
}

func setConditionLabel(f *goast.ForStmt, p *program.Program) (label string) {
	if f.Body == nil || len(f.Body.List) == 0 {
		return ""
	}
	lastStmt := f.Body.List[len(f.Body.List)-1]
	if _, ok := lastStmt.(*goast.IfStmt); ok {
		doWhileCondLabel := p.GetNextIdentifier("DO_WHILE_COND_LABEL_")
		f.Body.List[len(f.Body.List)-1] = &goast.LabeledStmt{
			Label: util.NewIdent(doWhileCondLabel),
			Stmt:  lastStmt,
		}
		return doWhileCondLabel
	}
	return ""
}

// createIfWithNotConditionAndBreak - create operator IF like on next example
// of C code:
// if ( !(condition) ) {
//		break;
// }
// Example of AST tree:
//  `-IfStmt 0x3bb1da8 <line:15:3, line:17:3>
//    |-<<<NULL>>>
//    |-UnaryOperator 0x3bb1d60 <line:15:6, col:11> 'int' prefix '!'
//    | `-ParenExpr 0x3bb1d40 <col:7, col:11> 'int'
//    |   `- CONDITION
//    |-CompoundStmt 0x3bb1d88 <col:13, line:17:3>
//    | `-BreakStmt 0x3bb1d80 <line:16:4>
//    `-<<<NULL>>>
func createIfWithNotConditionAndBreak(condition ast.Node) (ifStmt ast.IfStmt) {
	ifStmt.AddChild(nil)

	var par ast.ParenExpr
	var unitary ast.UnaryOperator
	switch con := condition.(type) {
	case *ast.BinaryOperator:
		par.Type = con.Type
		unitary.Type = con.Type

	case *ast.ImplicitCastExpr:
		par.Type = con.Type
		unitary.Type = con.Type

	case *ast.CStyleCastExpr:
		par.Type = con.Type
		unitary.Type = con.Type

	case *ast.ParenExpr:
		par.Type = con.Type
		unitary.Type = con.Type

	case *ast.UnaryOperator:
		par.Type = con.Type
		unitary.Type = con.Type

	case *ast.IntegerLiteral:
		par.Type = con.Type
		unitary.Type = con.Type

	default:
		panic(
			fmt.Errorf("Type %T is not implemented in createIfWithNotConditionAndBreak", condition))
	}
	par.AddChild(condition)
	unitary.Operator = "!"
	unitary.AddChild(&par)

	ifStmt.AddChild(&unitary)
	var c ast.CompoundStmt
	c.AddChild(&ast.BreakStmt{})
	ifStmt.AddChild(&c)
	ifStmt.AddChild(nil)

	return
}

func transpileContinueStmt(n *ast.ContinueStmt, p *program.Program) (*goast.BranchStmt, error) {
	return &goast.BranchStmt{
		Tok: token.CONTINUE,
	}, nil
}
