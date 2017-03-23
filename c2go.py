import sys
import pprint
import re
import subprocess
import StringIO
import json

function_defs = {
    '__istype': ('uint32', ('__darwin_ct_rune_t', 'uint32')),
    '__isctype': ('__darwin_ct_rune_t', ('__darwin_ct_rune_t', 'uint32')),
    '__tolower': ('__darwin_ct_rune_t', ('__darwin_ct_rune_t',)),
    '__toupper': ('__darwin_ct_rune_t', ('__darwin_ct_rune_t',)),
    '__maskrune': ('uint32', ('__darwin_ct_rune_t', 'uint32')),
}

function_subs = {
    # stdio
    'printf': 'fmt.Printf',
    'scanf': 'fmt.Scanf',
}

imports = ["fmt"]

def add_import(import_name):
    if import_name not in imports:
        imports.append(import_name)

def is_keyword(w):
    return w in ('char', 'long', 'struct', 'void')

def is_identifier(w):
    return not is_keyword(w) and re.match('[_a-zA-Z][_a-zA-Z0-9]*', w)

def resolve_type(s):
    s = s.strip()

    if s == 'const char *' or s == 'const char*' or s == 'char *' or \
        s == 'const char *restrict' or s == 'const char *__restrict':
        return 'string'

    if s == 'float':
        return 'float32'

    if s == 'void *':
        return 'interface{}'

    if s == 'char':
        return 'byte'

    if s == 'int *':
        return '*int'

    if s == 'unsigned long':
        return 'uint32'

    if s == 'int' or s == '__darwin_ct_rune_t':
        return s

    if s == 'long':
        return 'int64'

    if s == 'long long':
        return 'int64'

    if s == 'signed char':
        return 'int8'

    if s == 'unsigned char':
        return 'uint8'

    if s == 'unsigned char *':
        return '*uint8'

    if s == 'unsigned short':
        return 'uint16'

    if s == 'short':
        return 'int16'

    if s == 'unsigned int' or s == 'long unsigned int':
        return 'uint32'

    if s == 'unsigned long long':
        return 'uint64'

    if s == 'long int':
        return 'int32'

    if re.match('unsigned char \\[\\d+\\]', s):
        return s[14:] + 'byte'

    if re.match('char \\[\\d+\\]', s):
        return s[5:] + 'byte'

    if re.match('int \\[\\d+\\]', s):
        return s[4:] + 'int'

    if s[:7] == 'struct ':
        return resolve_type(s[7:])

    if '(*)' in s or s == '__sFILEX *' or s == 'fpos_t':
        return "interface{}"

    # return s

    raise Exception('Cannot resolve type "%s"' % s)

def cast(expr, from_type, to_type):
    from_type = resolve_type(from_type)
    to_type = resolve_type(to_type)

    if from_type == to_type:
        return expr

    types = ('int', 'int64', 'uint32', '__darwin_ct_rune_t', 'byte')
    if from_type in types and to_type == 'bool':
        return '%s != 0' % expr

    if from_type == '*int' and to_type == 'bool':
        return '%s != nil' % expr

    if from_type in types and to_type in types:
        return '%s(%s)' % (to_type, expr)

    return '__%s_to_%s(%s)' % (from_type, to_type, expr)

def print_line(out, line, indent):
    out.write('%s%s\n' % ('\t' * indent, line))

def render_expression(node):
    if node['node'] == 'BINARY_OPERATOR':
        end_of_left = list(node.get_children())[0].extent.end.column
        operator = None
        for t in node.get_tokens():
            if t.extent.start.column >= end_of_left:
                operator = t.spelling
                break

        left, right = [render_expression(t)[0] for t in list(node.get_children())]

        return_type = 'bool'
        if operator == '|' or operator == '&':
            return_type = 'int64'

        return '%s %s %s' % (left, operator, right), return_type

    if node['node'] == 'CONDITIONAL_OPERATOR':
        a, b, c = [render_expression(t) for t in list(node.get_children())]
        try:
            return '__ternary(%s, %s, %s)' % (cast(a[0], 'bool'), b[0], c[0]), b[1]
        except TypeError:
            return '// CONDITIONAL_OPERATOR: %s' % ''.join([t.spelling for t in node.get_tokens()]), 'unknown'

    if node['node'] == 'UNARY_OPERATOR':
        expr_start = list(node.get_children())[0].extent.start.column
        operator = None
        for t in node.get_tokens():
            if t.extent.start.column >= expr_start:
                break

            operator = t.spelling

        if operator is None:
            operator = '++'

        expr = render_expression(list(node.get_children())[0])

        if operator == '!':
            return '%s(%s)' % ('__not_%s' % expr[1], expr[0]), expr[1]

        if operator == '*':
            if expr[1] == 'const char *':
                return '%s[0]' % expr[0], 'char'

            return '*%s' % expr[0], 'int'

        if operator == '++':
            return '%s += 1' % expr[0], expr[1]

        if operator == '~':
            operator = '^'

        return '%s%s' % (operator, expr[0]), expr[1]

    if node['node'] == 'UNEXPOSED_EXPR':
        children = list(node.get_children())
        if len(children) < 1:
            return '// UNEXPOSED_EXPR: %s' % ''.join([t.spelling for t in node.get_tokens()]), 'unknown'

        # if len(children) > 1:
        #     raise Exception('To many children!')

        e = render_expression(children[0])
        name = e[0]

        if name == 'argc':
            name = 'len(os.Args)'
            add_import("os")
        elif name == 'argv':
            name = 'os.Args'
            add_import("os")

        return name, e[1]

    if node['node'] in ('CHARACTER_LITERAL', 'StringLiteral', 'FLOATING_LITERAL'):
        return node['value'], 'const char*'

    if node['node'] == 'INTEGER_LITERAL':
        literal = list(node.get_tokens())[0].spelling
        if literal[-1] == 'L':
            literal = '%s(%s)' % (resolve_type('long'), literal[:-1])

        return literal, 'int'

    if node['node'] == 'PAREN_EXPR':
        e = render_expression(list(node.get_children())[0])
        return '(%s)' % e[0], e[1]

    if node['node'] == 'DeclRefExpr':
        return node['name'], node['type']

    if node['node'] == 'ImplicitCastExpr':
        return render_expression(node['children'][0])

    if node['node'] == 'CallExpr':
        children = node['children']
        func_name = render_expression(children[0])[0]

        func_def = function_defs[func_name]

        if func_name in function_subs:
            func_name = function_subs[func_name]

        args = []
        i = 0
        for arg in children[1:]:
            e = render_expression(arg)

            if i > len(func_def[1]) - 1:
                # This means the argument is one of the varargs so we don't know
                # what type it needs to be cast to.
                args.append(e[0])
            else:
                args.append(cast(e[0], e[1], func_def[1][i]))

            i += 1

        return '%s(%s)' % (func_name, ', '.join(args)), func_def[0]

    if node['node'] == 'ARRAY_SUBSCRIPT_EXPR':
        children = list(node.get_children())
        return '%s[%s]' % (render_expression(children[0])[0],
            render_expression(children[1])[0]), 'unknown'

    if node['node'] == 'MEMBER_REF_EXPR':
        children = list(node.get_children())
        return '%s.%s' % (render_expression(children[0])[0], list(node.get_tokens())[-2].spelling), 'unknown'

    if node['node'] == 'CSTYLE_CAST_EXPR':
        children = list(node.get_children())
        return render_expression(children[0]), 'unknown'

    if node['node'] == 'FIELD_DECL' or node['node'] == 'VAR_DECL':
        type = resolve_type(node.type.spelling)
        name = node.spelling

        prefix = ''
        if node['node'] == 'VAR_DECL':
            prefix = 'var '

        suffix = ''
        children = list(node.get_children())

        # We must check the position of the child is at the end. Otherwise a
        # child can refer to another expression like the size of the data type.
        if len(children) > 0 and children[0].extent.end.column == node.extent.end.column:
            e = render_expression(children[0])
            suffix = ' = %s' % cast(e[0], e[1], type)

        return '%s%s %s%s' % (prefix, name, type, suffix), 'unknown'

    if node['node'] == 'PARM_DECL':
        return resolve_type(node.type.spelling), 'unknown'

    # return node['node'], 'unknown'

    raise Exception('render_expression: %s' % node['node'])

def print_children(node):
    print(len(list(node.get_children())), [t.spelling for t in node.get_tokens()])
    for child in node.get_children():
        print(child.kind.name, render_expression(child), [t.spelling for t in child.get_tokens()])

def get_function_params(nodes):
    if 'children' not in nodes:
        return []

    return [n for n in nodes['children'] if n['node'] == 'ParmVarDecl']

def render(out, node, indent=0, return_type=None):
    if node['node'] == 'TranslationUnitDecl':
        for c in node['children']:
            render(out, c, indent, return_type)
        return

    if node['node'] == 'FunctionDecl':
        function_name = node['name']

        if function_name in ('__istype', '__isctype', '__wcwidth', '__sputc'):
            return

        has_body = False
        if 'children' in node:
            for c in node['children']:
                if c['node'] == 'CompoundStmt':
                    has_body = True

        args = []
        # for a in get_function_params(node):
        #     args.append('%s %s' % (a['name'], resolve_type(a['type'])))

        if has_body:
            return_type = ' ' + node['type']
            if return_type == ' void':
                return_type = ''

            if function_name == 'main':
                print_line(out, 'func main() {', indent)
            else:
                print_line(out, 'func %s(%s)%s {' % (function_name,
                    ', '.join(args), return_type), indent)
            
            for c in node['children']:
                if c['node'] == 'CompoundStmt':
                    render(out, c, indent + 1, node['type'])

            print_line(out, '}\n', indent)

        function_defs[node['name']] = (node['type'], [a['type'] for a in get_function_params(node)])

        return

    # if node['node'] == 'PARM_DECL':
    #     print_line(out, node.spelling, indent)
    #     return

    if node['node'] == 'CompoundStmt':
        for c in node['children']:
            render(out, c, indent, return_type)
        return

    # if node['node'] == 'IF_STMT':
    #     children = list(node.get_children())

    #     e = render_expression(children[0])
    #     print_line(out, 'if %s {' % cast(e[0], e[1], 'bool'), indent)

    #     render(out, children[1], indent + 1, return_type)

    #     if len(children) > 2:
    #         print_line(out, '} else {', indent)
    #         render(out, children[2], indent + 1, return_type)

    #     print_line(out, '}', indent)

        # return

    # if node['node'] == 'WHILE_STMT':
    #     children = list(node.get_children())

    #     e = render_expression(children[0])
    #     print_line(out, 'for %s {' % cast(e[0], e[1], 'bool'), indent)

    #     render(out, children[1], indent + 1, return_type)

    #     print_line(out, '}', indent)

    #     return

    # if node['node'] == 'FOR_STMT':
    #     children = list(node.get_children())

    #     a, b, c = [render_expression(e)[0] for e in children[:3]]
    #     print_line(out, 'for %s; %s; %s {' % (a, b, c), indent)

    #     render(out, children[3], indent + 1, return_type)

    #     print_line(out, '}', indent)

    #     return

    # if node['node'] == 'BREAK_STMT':
    #     print_line(out, 'break', indent)
    #     return

    # if node['node'] == 'UNARY_OPERATOR':
    #     variable, operator = [t.spelling for t in list(node.get_tokens())[0:2]]
    #     if operator == '++':
    #         print_line(out, '%s += 1' % variable, indent)
    #         #print_line(out, '%s = string(%s[1:])' % (variable, variable), indent)
    #         return

    #     print_line(out, '%s%s' % (operator, variable), indent)
    #     return

    #     #raise Exception('UNARY_OPERATOR: %s' % operator)

    if node['node'] == 'ReturnStmt':
        # try:
        #     e = render_expression(list(node.get_children())[0])
        #     print_line(out, 'return %s' % cast(e[0], e[1], return_type), indent)
        # except IndexError:
        print_line(out, 'return', indent)
        
        return

    if node['node'] in ('BINARY_OPERATOR', 'INTEGER_LITERAL', 'CallExpr'):
        print_line(out, render_expression(node)[0], indent)
        return

    if node['node'] == 'TypedefDecl':
        print_line(out, "type %s %s\n" % (node['type'], node['name']), indent)
        # print(node)
        return

        tokens = [t.spelling for t in node.get_tokens()]
        if len(list(node.get_children())) == 0:
            print_line(out, "type %s %s\n" % (tokens[-2], resolve_type(' '.join(tokens[1:-2]))), indent)
        #else:
        #    print_line(out, "type %s %s\n" % (tokens[-2], render(out, list(node.get_children())[0], indent, return_type)), indent)

        return

    if node['node'] == 'RecordDecl':
        return

    #if node['node'] == 'UNION_DECL' or node['node'] == 'STRUCT_DECL':
    #     tokens = [t.spelling for t in node.get_tokens()]

    #     struct_name = tokens[-1]
    #     start_at = 2
    #     if struct_name == ';':
    #         struct_name = tokens[1]
    #         start_at = 3

    #     if struct_name in ('__darwin_pthread_handler_rec', '_opaque_pthread_t',
    #         '_RuneEntry', '_RuneRange', '_RuneCharClass', '_RuneLocale'):
    #         return

    #     print_line(out, "type %s struct {" % struct_name, indent)

    #     for attribute in node.get_children():
    #         print_line(out, render_expression(attribute)[0], indent + 1)
    #         # print(struct_name, render_expression(attribute))

    #     # name = ''
    #     # type = ''
    #     # for token in tokens[start_at:-2]:
    #     #     if token == ';':
    #     #         print_line(out, '%s %s' % (name, resolve_type(type)), indent + 1)
    #     #         type = ''
    #     #     elif is_identifier(token):
    #     #         name = token
    #     #     else:
    #     #         type += ' ' + token

    #     print_line(out, "}\n", indent)
    #     return

    # if node['node'] == 'UNEXPOSED_DECL':
    #     tokens = [t.spelling for t in node.get_tokens()]
    #     print_line(out, '// ' + ' '.join(tokens[1:-2]), indent)
    #     return

    # if node['node'] == 'DECL_STMT':
    #     for child in node.get_children():
    #         print_line(out, render_expression(child)[0], indent)
    #     return

    if node['node'] == 'VarDecl':
    #     tokens = [t.spelling for t in node.get_tokens()]
    #     if tokens[0] == 'extern':
    #         return

    #     children = list(node.get_children())
    #     if len(children) > 0:
    #         print_line(out, 'var %s %s = %s\n' % (tokens[2], tokens[1], render_expression(children[0])[0]), indent)
    #     else:
    #         print_line(out, 'var %s %s\n' % (tokens[2], tokens[1]), indent)
        
        return

    # if node['node'] == 'ENUM_DECL':
    #     print_line(out, '// enum', indent)
    #     return

    raise Exception(node['node'])

# 1. Compile it first (checking for errors)
c_file_path = sys.argv[1]
#subprocess.call(["clang", c_file_path])

# 2. Preprocess
pp = subprocess.Popen(["clang", "-E", c_file_path], stdout=subprocess.PIPE).communicate()[0]

pp_file_path = 'pp.c'
with open(pp_file_path, 'w') as pp_out:
    pp_out.write(pp)

# 3. Generate JSON from AST
ast_pp = subprocess.Popen(["clang", "-Xclang", "-ast-dump", "-fsyntax-only", pp_file_path], stdout=subprocess.PIPE)
pp = subprocess.Popen(["python", "ast2json.py"], stdin=ast_pp.stdout, stdout=subprocess.PIPE).communicate()[0]

json_file_path = 'pp.json'
with open(json_file_path, 'w') as json_out:
    json_out.write(pp)

with open(json_file_path, 'r') as json_in:
    # 3. Parse C and output Go
    # index = clang.cindex.Index.create()
    # tu = index.parse(pp_file_path)

    go_file_path = '%s.go' % c_file_path.split('/')[-1][:-2]
    # go_out = sys.stdout
    go_out = StringIO.StringIO()
    #with open(go_file_path, 'w') as go_out:
    # print_line(go_out, "package main\n", 0)
    #print_line(go_out, 'import ("fmt"; "os")\n', 0)
    render(go_out, json.loads(json_in.read())[0])

    print("package main\n")
    print("import (")
    for import_name in sorted(imports):
        print('\t"%s"' % import_name)
    print(")\n")
    print(go_out.getvalue())

    # 4. Compile the Go
    #subprocess.call(["go", "run", "functions.go", go_file_path])
