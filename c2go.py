#!/usr/bin/env python
# -*- coding: utf-8 -*-

import sys
import pprint
import re
import subprocess
import json

try:
    import StringIO as io
except ImportError:
    import io

function_defs = {
    '__istype': ('uint32', ('__darwin_ct_rune_t', 'uint32')),
    '__isctype': ('__darwin_ct_rune_t', ('__darwin_ct_rune_t', 'uint32')),
    '__tolower': ('__darwin_ct_rune_t', ('__darwin_ct_rune_t',)),
    '__toupper': ('__darwin_ct_rune_t', ('__darwin_ct_rune_t',)),
    '__maskrune': ('uint32', ('__darwin_ct_rune_t', 'uint32')),

    # These are provided by functions-Darwin.go
    '__builtin_fabs': ('double', ('double',)),
    '__builtin_fabsf': ('float', ('float',)),
    '__builtin_fabsl': ('double', ('double',)),
    '__builtin_inf': ('double', ()),
    '__builtin_inff': ('float', ()),
    '__builtin_infl': ('double', ()),
}

function_subs = {
    # math.h
    'cos': 'math.Cos',

    # stdio
    'printf': 'fmt.Printf',
    'scanf': 'fmt.Scanf',
}

# TODO: Some of these are based on assumtions that may not be true for all
# architectures (like the size of an int). At some point in the future we will
# need to find out the sizes of some of there and pick the most compatible type.
# 
# Please keep them sorted by name.
simple_resolve_types = {
    'bool': 'bool',
    'char *': 'string',
    'char': 'byte',
    'char*': 'string',
    'double': 'float64',
    'float': 'float32',
    'int': 'int',
    'long double': 'float64',
    'long int': 'int32',
    'long long': 'int64',
    'long unsigned int': 'uint32',
    'long': 'int32',
    'short': 'int16',
    'signed char': 'int8',
    'unsigned char': 'uint8',
    'unsigned int': 'uint32',
    'unsigned long long': 'uint64',
    'unsigned long': 'uint32',
    'unsigned short': 'uint16',
    'void *': 'interface{}',
    'void': '',

    'const char *': 'string',

    # Mac specific
    '__darwin_ct_rune_t': '__darwin_ct_rune_t',

    # These are special cases that almost certainly don't work. I've put them
    # here becuase for whatever reason there is no suitable type or we don't
    # need these platform specific things to be implemented yet.
    '__builtin_va_list': 'int64',
    '__darwin_pthread_handler_rec': 'int64',
    '__int128': 'int64',
    '__mbstate_t': 'int64',
    '__sbuf': 'int64',
    '__sFILEX': 'interface{}',
    '__va_list_tag': 'interface{}',
    'FILE': 'int64',
}

types_already_defined = set()

imports = ["fmt"]

class NoSuchTypeException(Exception):
    pass

def add_import(import_name):
    if import_name not in imports:
        imports.append(import_name)

def is_keyword(w):
    return w in ('char', 'long', 'struct', 'void')

def is_identifier(w):
    return not is_keyword(w) and re.match('[_a-zA-Z][_a-zA-Z0-9]*', w)

def resolve_type(s):
    # Remove any whitespace or attributes that are not relevant to Go.
    s = s.replace('const ', '')
    s = s.replace('*__restrict', '*')
    s = s.replace('*restrict', '*')
    s = s.strip(' \t\n\r')

    # If the type is already defined we can proceed with the same name.
    if s in types_already_defined:
        return s

    # The simple resolve types are the types that we know there is an exact Go
    # equivalent. For example float, int, etc.
    if s in simple_resolve_types:
        return simple_resolve_types[s]

    # Structures are by name.
    if s[:7] == 'struct ':
        if s[-1] == '*':
            return '*' + s[7:-2]
        else:
            return s[7:]

    # I have no idea how to handle this yet.
    if 'anonymous union' in s:
        return 'interface{}'

    # It may be a pointer of a simple type. For example, float *, int *, etc.
    try:
        if re.match(r"[\w ]+\*", s):
            return '*' + resolve_type(s[:-2].strip())
    except NoSuchTypeException:
        # Keep trying the next one.
        pass

    # Function pointers are not yet supported. In th mean time they will be
    # replaced with a type that certainly wont work until we can fix this
    # properly.
    search = re.search(r"[\w ]+\(\*.*?\)\(.*\)", s)
    if search:
        return 'interface{}'

    try:
        # It could be an array of fixed length.
        search = re.search(r"([\w ]+)\[(\d+)\]", s)
        if search:
            return '[%s]%s' % (search.group(2), resolve_type(search.group(1)))

    except NoSuchTypeException as e:
        # Make the nested exception message more contextual.
        raise NoSuchTypeException(e.message + " (from '%s')" % s)

    raise NoSuchTypeException("'%s'" % s)

def cast(expr, from_type, to_type):
    from_type = resolve_type(from_type)
    to_type = resolve_type(to_type)

    if from_type == to_type:
        return expr

    types = ('int', 'int64', 'uint32', '__darwin_ct_rune_t', 'byte', 'float32',
        'float64')
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
    if node['node'] == 'BinaryOperator':
        operator = node['operator']

        left, left_type = render_expression(node['children'][0])
        right, right_type = render_expression(node['children'][1])

        return_type = 'bool'
        if operator in ('|', '&', '+', '-', '*', '/'):
            # TODO: The left and right type might be different
            return_type = left_type

        if operator == '&&':
            left = cast(left, left_type, return_type)
            right = cast(right, right_type, return_type)

        return '%s %s %s' % (left, operator, right), return_type

    if node['node'] == 'UnaryOperator':
        operator = node['operator']
        expr = render_expression(node['children'][0])

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

    if node['node'] == 'StringLiteral':
        return node['value'], 'const char *'

    if node['node'] == 'FloatingLiteral':
        return node['value'], 'double'

    if node['node'] == 'IntegerLiteral':
        literal = node['value']
        if literal[-1] == 'L':
            literal = '%s(%s)' % (resolve_type('long'), literal[:-1])

        return literal, 'int'

    if node['node'] == 'DeclRefExpr':
        name = node['name']

        if name == 'argc':
            name = 'len(os.Args)'
            add_import("os")
        elif name == 'argv':
            name = 'os.Args'
            add_import("os")

        return name, node['type']

    if node['node'] == 'ImplicitCastExpr':
        return render_expression(node['children'][0])

    if node['node'] == 'CallExpr':
        children = node['children']
        func_name = render_expression(children[0])[0]

        func_def = function_defs[func_name]

        if func_name in function_subs:
            func_name = function_subs[func_name]
            add_import(func_name.split('.')[0])

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

    if node['node'] == 'ArraySubscriptExpr':
        children = node['children']
        return '%s[%s]' % (render_expression(children[0])[0],
            render_expression(children[1])[0]), 'unknown1'

    if node['node'] == 'MemberExpr':
        children = node['children']
        return '%s.%s' % (render_expression(children[0])[0], node['name']), children[0]['type']

    if node['node'] == 'CStyleCastExpr':
        children = node['children']
        return render_expression(children[0])

    if node['node'] == 'FieldDecl' or node['node'] == 'VarDecl':
        type = resolve_type(node['type'])
        name = node['name'].replace('used', '')

        # Go does not allow the name of a variable to be called "type". For the
        # moment I will rename this to avoid the error.
        if name == 'type':
            name = 'type_'

        prefix = ''
        if node['node'] == 'VarDecl':
            prefix = 'var '

        suffix = ''
        if 'children' in node:
            children = node['children']
            suffix = ' = %s' % render_expression(children[0])[0]

        return '%s%s %s%s' % (prefix, name, type, suffix), 'unknown3'

    if node['node'] == 'RecordDecl':
        return '/* RecordDecl */', 'unknown5'

    if node['node'] == 'ParenExpr':
        a, b = render_expression(node['children'][0])
        return '(%s)' % a, b

    # return node['node'], 'unknown6'

    raise Exception('render_expression: %s' % node['node'])

def print_children(node):
    print(len(list(node.get_children())), [t.spelling for t in node.get_tokens()])
    for child in node.get_children():
        print(child.kind.name, render_expression(child), [t.spelling for t in child.get_tokens()])

def get_function_params(nodes):
    if 'children' not in nodes:
        return []

    return [n for n in nodes['children'] if n['node'] == 'ParmVarDecl']

def get_function_return_type(f):
    # The type of the function will be the complete prototype, like:
    # 
    #     __inline_isfinitef(float) int
    #     
    # will have a type of:
    #
    #     int (float)
    #
    # The arguments will handle themselves, we only care about the
    # return type ('int' in this case)
    return f.split('(')[0].strip()

def render(out, node, indent=0, return_type=None):
    if node['node'] == 'TranslationUnitDecl':
        for c in node['children']:
            render(out, c, indent, return_type)
        return

    if node['node'] == 'FunctionDecl':
        function_name = node['name'].strip()

        if function_name in ('__istype', '__isctype', '__wcwidth', '__sputc',
            '__inline_signbitf', '__inline_signbitd', '__inline_signbitl'):
            return

        has_body = False
        if 'children' in node:
            for c in node['children']:
                if c['node'] == 'CompoundStmt':
                    has_body = True

        args = []
        for a in get_function_params(node):
            args.append('%s %s' % (a['name'], resolve_type(a['type'])))

        if has_body:
            return_type = get_function_return_type(node['type'])

            if function_name == 'main':
                print_line(out, 'func main() {', indent)
            else:
                print_line(out, 'func %s(%s) %s {' % (function_name,
                    ', '.join(args), resolve_type(return_type)), indent)
            
            for c in node['children']:
                if c['node'] == 'CompoundStmt':
                    render(out, c, indent + 1, node['type'])

            print_line(out, '}\n', indent)

        function_defs[node['name']] = (get_function_return_type(node['type']),
            [a['type'] for a in get_function_params(node)])

        return

    # if node['node'] == 'PARM_DECL':
    #     print_line(out, node.spelling, indent)
    #     return

    if node['node'] == 'CompoundStmt':
        for c in node['children']:
            render(out, c, indent, return_type)
        return

    if node['node'] == 'IfStmt':
        children = node['children']

        e = render_expression(children[0])
        print_line(out, 'if %s {' % cast(e[0], e[1], 'bool'), indent)

        render(out, children[1], indent + 1, return_type)

        if len(children) > 2:
            print_line(out, '} else {', indent)
            render(out, children[2], indent + 1, return_type)

        print_line(out, '}', indent)

        return

    if node['node'] == 'WhileStmt':
        children = node['children']

        e = render_expression(children[0])
        print_line(out, 'for %s {' % cast(e[0], e[1], 'bool'), indent)

        render(out, children[1], indent + 1, return_type)

        print_line(out, '}', indent)

        return

    if node['node'] == 'ForStmt':
        children = node['children']

        a, b, c = [render_expression(e)[0] for e in children[:3]]
        print_line(out, 'for %s; %s; %s {' % (a, b, c), indent)

        render(out, children[3], indent + 1, return_type)

        print_line(out, '}', indent)

        return

    if node['node'] == 'BreakStmt':
        print_line(out, 'break', indent)
        return

    if node['node'] == 'UnaryOperator':
        print_line(out, render_expression(node)[0], indent)
        return

    if node['node'] == 'ReturnStmt':
        r = 'return'

        # This special return type is for main().
        if 'children' in node and return_type != 'int ()':
            expr, type = render_expression(node['children'][0])
            r = 'return ' + cast(expr, type, 'int')

        print_line(out, r, indent)
        return

    if node['node'] in ('BinaryOperator', 'INTEGER_LITERAL', 'CallExpr'):
        print_line(out, render_expression(node)[0], indent)
        return

    if node['node'] == 'TypedefDecl':
        types_already_defined.add(node['name'].strip())

        # FIXME: All of the logic here is just to avoid errors, it needs to be
        # fixed up.
        if 'struct' in node['type'] or 'union' in node['type']:
            return
        node['type'] = node['type'].replace('unsigned', '')
        if node['name'] in ('__builtin_va_list', '__qaddr_t', 'definition',
            '_IO_lock_t', 'va_list', 'fpos_t'):
            return

        print_line(out, "type %s %s\n" % (node['name'], resolve_type(node['type'])), indent)

        return

    if node['node'] == 'EnumDecl':
        return

    if node['node'] == 'FieldDecl':
        print_line(out, render_expression(node)[0], indent + 1)
        return

    if node['node'] == 'RecordDecl':
        if node['kind'] == 'union':
            return

        # FIXME
        if node['name'] in ('definition', '_IO_FILE'):
            return

        print_line(out, "type %s %s {" % (node['name'], node['kind']), indent)
        if 'children' in node:
            for c in node['children']:
                render(out, c, indent + 1)
        print_line(out, "}\n", indent)
        return

    #if node['node'] == 'UNION_DECL' or node['node'] == 'STRUCT_DECL':
    #     tokens = [t.spelling for t in node.get_tokens()]

    #     struct_name = tokens[-1]
    #     start_at = 2
    #     if struct_name == ';':
    #         struct_name = tokens[1]
    #         start_at = 3

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

    if node['node'] == 'DeclStmt':
        for child in node['children']:
            print_line(out, render_expression(child)[0], indent)
        return

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
with open(pp_file_path, 'wb') as pp_out:
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
    go_out = io.StringIO()
    all_json = json_in.read()

    try:
        l = json.loads(all_json)
    except ValueError as e:
        # This occurs if the JSON cannot be parsed
        print(all_json)
        raise e

    render(go_out, l[0])

    print("package main\n")
    print("import (")
    for import_name in sorted(imports):
        print('\t"%s"' % import_name)
    print(")\n")
    print(go_out.getvalue())

    # 4. Compile the Go
    #subprocess.call(["go", "run", "functions.go", go_file_path])
