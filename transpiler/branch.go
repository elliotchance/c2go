// This file contains functions for transpiling common branching and control
// flow, such as "if", "while", "do" and "for". The more complicated control
// flows like "switch" will be put into their own file of the same or sensible
// name.

package transpiler

import (
	"fmt"
	"go/token"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"

	goast "go/ast"
)

func transpileIfStmt(n *ast.IfStmt, p *program.Program) (
	*goast.IfStmt, []goast.Stmt, []goast.Stmt, error) {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}
	children := n.Children

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

	conditional, conditionalType, newPre, newPost, err := transpileToExpr(children[1], p)
	if err != nil {
		return nil, nil, nil, err
	}

	// The condition in Go must always be a bool.
	boolCondition, err := types.CastExpr(p, conditional, conditionalType, "bool")
	p.AddMessage(ast.GenerateWarningOrErrorMessage(err, n, boolCondition == nil))

	if boolCondition == nil {
		boolCondition = util.NewNil()
	}

	body, newPre, newPost, err := transpileToBlockStmt(children[2], p)
	if err != nil {
		return nil, nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

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

		r.Else = elseBody
	}

	return r, newPre, newPost, nil
}

func transpileForStmt(n *ast.ForStmt, p *program.Program) (
	*goast.ForStmt, []goast.Stmt, []goast.Stmt, error) {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	children := n.Children

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

	// If we have 2 and more initializations like
	// in operator for
	// for( a = 0, b = 0, c = 0; a < 5; a ++)
	switch c := children[0].(type) {
	case *ast.BinaryOperator:
		if c.Operator == "," {
			// recursive action to code like that:
			// a = 0;
			// b = 0;
			// for(c = 0 ; a < 5 ; a++)
			before, newPre, newPost, err := transpileToStmt(c.Children[0], p)
			if err != nil {
				return nil, nil, nil, err
			}
			preStmts = append(preStmts, newPre...)
			preStmts = append(preStmts, before)
			preStmts = append(preStmts, newPost...)
			children[0] = c.Children[1]
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
				compound = children[4].(*ast.CompoundStmt)
			} else {
				// if body is not exist
				compound = new(ast.CompoundStmt)
			}
			compound.Children = append(compound.Children, c.Children[0:len(c.Children)]...)
			children[4] = compound
			children[3] = nil
		}
	}

	post, newPre, newPost, err := transpileToStmt(children[3], p)
	if err != nil {
		return nil, nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

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
			tempSlice := c.Children[0 : len(c.Children)-1]

			var condition ast.IfStmt
			condition.AddChild(nil)
			var par ast.ParenExpr
			par.AddChild(c.Children[len(c.Children)-1])
			var unitary ast.UnaryOperator
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
			compound.Children = append(tempSlice, compound.Children...)
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
		condition, conditionType, newPre, newPost, err = transpileToExpr(children[2], p)
		if err != nil {
			return nil, nil, nil, err
		}

		preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

		condition, err = types.CastExpr(p, condition, conditionType, "bool")
		p.AddMessage(ast.GenerateWarningOrErrorMessage(err, n, condition == nil))

		if condition == nil {
			condition = util.NewNil()
		}
	}

	body, newPre, newPost, err := transpileToBlockStmt(children[4], p)
	if err != nil {
		return nil, nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	return &goast.ForStmt{
		Init: init,
		Cond: condition,
		Post: post,
		Body: body,
	}, preStmts, postStmts, nil
}

// transpileWhileStmt - transpiler for operator While.
// We have only operator FOR in Go, but in C we also have
// operator WHILE. So, we have to convert to operator FOR.
// We choose directly convertion  from AST C code to AST C code, for
// - avoid dublicate of code in realization WHILE and FOR.
// - create only one operator FOR powerfull.
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
	*goast.ForStmt, []goast.Stmt, []goast.Stmt, error) {
	var forOperator ast.ForStmt
	forOperator.AddChild(nil)
	forOperator.AddChild(nil)
	forOperator.AddChild(n.Children[1])
	forOperator.AddChild(nil)
	forOperator.AddChild(n.Children[2])
	return transpileForStmt(&forOperator, p)
}

func transpileDoStmt(n *ast.DoStmt, p *program.Program) (
	*goast.ForStmt, []goast.Stmt, []goast.Stmt, error) {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}
	children := n.Children

	body, newPre, newPost, err := transpileToBlockStmt(children[0], p)
	if err != nil {
		return nil, nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	condition, conditionType, newPre, newPost, err := transpileToExpr(children[1], p)
	if err != nil {
		return nil, nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	// Add IfStmt to the end of the loop to check the condition.
	x, err := types.CastExpr(p, condition, conditionType, "bool")
	p.AddMessage(ast.GenerateWarningOrErrorMessage(err, n, x == nil))

	if x == nil {
		x = util.NewNil()
	}

	body.List = append(body.List, &goast.IfStmt{
		Cond: &goast.UnaryExpr{
			Op: token.NOT,
			X:  x,
		},
		Body: &goast.BlockStmt{
			List: []goast.Stmt{&goast.BranchStmt{Tok: token.BREAK}},
		},
	})

	return &goast.ForStmt{
		Body: body,
	}, preStmts, postStmts, nil
}

func transpileContinueStmt(n *ast.ContinueStmt, p *program.Program) (*goast.BranchStmt, error) {
	return &goast.BranchStmt{
		Tok: token.CONTINUE,
	}, nil
}
