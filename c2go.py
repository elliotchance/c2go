#!/usr/bin/env python
# -*- coding: utf-8 -*-

import sys
import pprint
import re
import subprocess
import json
import os
import platform

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

    # darwin/math.h
    '__builtin_fabs': ('double', ('double',)),
    '__builtin_fabsf': ('float', ('float',)),
    '__builtin_fabsl': ('double', ('double',)),
    '__builtin_inf': ('double', ()),
    '__builtin_inff': ('float', ()),
    '__builtin_infl': ('double', ()),

    '__sincospi_stret': ('Double2', ('float',)),
    '__sincospif_stret': ('Float2', ('float',)),
    '__sincos_stret': ('Double2', ('float',)),
    '__sincosf_stret': ('Float2', ('float',)),

    # darwin/assert.h
    '__builtin_expect': ('int', ('int', 'int')),
    '__assert_rtn': ('bool', ('const char*', 'const char*', 'int', 'const char*')),

    # linux/assert.h
    '__assert_fail': ('bool', ('const char*', 'const char*', 'unsigned int', 'const char*')),
}

function_subs = {
    # math.h
    'acos': 'math.Acos',
    'asin': 'math.Asin',
    'atan': 'math.Atan',
    'atan2': 'math.Atan2',
    'ceil': 'math.Ceil',
    'cos': 'math.Cos',
    'cosh': 'math.Cosh',
    'exp': 'math.Exp',
    'fabs': 'math.Abs',
    'floor': 'math.Floor',
    'fmod': 'math.Mod',
    'ldexp': 'math.Ldexp',
    'log': 'math.Log',
    'log10': 'math.Log10',
    'pow': 'math.Pow',
    'sin': 'math.Sin',
    'sinh': 'math.Sinh',
    'sqrt': 'math.Sqrt',
    'tan': 'math.Tan',
    'tanh': 'math.Tanh',

    # darwin/math.h
    '__builtin_fabs': 'github.com/elliotchance/c2go/darwin.Fabs',
    '__builtin_fabsf': 'github.com/elliotchance/c2go/darwin.Fabsf',
    '__builtin_fabsl': 'github.com/elliotchance/c2go/darwin.Fabsl',
    '__builtin_inf': 'github.com/elliotchance/c2go/darwin.Inf',
    '__builtin_inff': 'github.com/elliotchance/c2go/darwin.Inff',
    '__builtin_infl': 'github.com/elliotchance/c2go/darwin.Infl',

    '__sincospi_stret': 'github.com/elliotchance/c2go/darwin.SincospiStret',
    '__sincospif_stret': 'github.com/elliotchance/c2go/darwin.SincospifStret',
    '__sincos_stret': 'github.com/elliotchance/c2go/darwin.SincosStret',
    '__sincosf_stret': 'github.com/elliotchance/c2go/darwin.SincosfStret',

    # stdio
    'printf': 'fmt.Printf',
    'scanf': 'fmt.Scanf',

    # assert
    '__builtin_expect': 'github.com/elliotchance/c2go/darwin.BuiltinExpect',
    '__assert_rtn': 'github.com/elliotchance/c2go/darwin.AssertRtn',

    # linux/assert.h
    '__assert_fail': 'github.com/elliotchance/c2go/linux.AssertFail',
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

    # Darwin specific
    '__darwin_ct_rune_t': '__darwin_ct_rune_t',
    'union __mbstate_t': '__mbstate_t',
    'fpos_t': 'int',
    'struct __float2': 'github.com/elliotchance/c2go/darwin.Float2',
    'struct __double2': 'github.com/elliotchance/c2go/darwin.Double2',

    # Linux specific:
    '_IO_FILE': 'github.com/elliotchance/c2go/linux.File',

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

types_already_defined = set([
    # Linux specific
    '_LIB_VERSION_TYPE',
    '_IO_FILE',

    # Darwin specific
    '__float2',
    '__double2',
])

imports = ["fmt"]

class NoSuchTypeException(Exception):
    pass

def ucfirst(word):
    return word[0].upper() + word[1:]

def get_exported_name(field):
    return ucfirst(field.lstrip('_'))

def add_import(import_name):
    if import_name not in imports:
        imports.append(import_name)

def import_type(type_name):
    if '.' in type_name:
        add_import('.'.join(type_name.split('.')[:-1]))
        return type_name.split('/')[-1]

    return type_name

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

    if s == 'fpos_t':
        return 'int'

    # The simple resolve types are the types that we know there is an exact Go
    # equivalent. For example float, int, etc.
    if s in simple_resolve_types:
        return import_type(simple_resolve_types[s])

    # If the type is already defined we can proceed with the same name.
    if s in types_already_defined:
        return import_type(s)

    # Structures are by name.
    if s[:7] == 'struct ':
        if s[-1] == '*':
            s = s[7:-2]

            if s in simple_resolve_types:
                return '*' + import_type(simple_resolve_types[s])

            return '*' + s
        else:
            s = s[7:]

            if s in simple_resolve_types:
                return import_type(simple_resolve_types[s])

            return s

    # Enums are by name.
    if s[:5] == 'enum ':
        if s[-1] == '*':
            return '*' + s[5:-2]
        else:
            return s[5:]

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
    search = re.search(r"[\w ]+ \(.*\)", s)
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

    add_import('github.com/elliotchance/c2go/noarch')
    return 'noarch.%sTo%s(%s)' % (ucfirst(from_type), ucfirst(to_type), expr)

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

        if (operator == '!=' or operator == '==') and right == '(0)':
            right = 'nil'

        return '%s %s %s' % (left, operator, right), return_type

    if node['node'] == 'UnaryOperator':
        operator = node['operator']
        expr = render_expression(node['children'][0])

        if operator == '!':
            if expr[1] == 'bool':
                return '!(%s)' % expr[0], expr[1]

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
        return '"%s"' % node['value'].replace("\n", "\\n"), 'const char *'

    if node['node'] == 'FloatingLiteral':
        return node['value'], 'double'

    if node['node'] == 'IntegerLiteral':
        literal = node['value']
        if str(literal)[-1] == 'L':
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
            add_import('.'.join(function_subs[func_name].split('.')[:-1]))
            func_name = function_subs[func_name].split('/')[-1]

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

        return '%s(%s)' % (func_name, ', '.join([str(a) for a in args])), func_def[0]

    if node['node'] == 'ArraySubscriptExpr':
        children = node['children']
        return '%s[%s]' % (render_expression(children[0])[0],
            render_expression(children[1])[0]), 'unknown1'

    if node['node'] == 'MemberExpr':
        children = node['children']

        lhs = render_expression(children[0])
        lhs_type = resolve_type(lhs[1])
        rhs = node['name']

        if lhs_type in ('darwin.Float2', 'darwin.Double2'):
            rhs = get_exported_name(rhs)

        return '%s.%s' % (lhs[0], rhs), children[0]['type']

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

            if suffix == ' = (0)':
                suffix = ' = nil'

        return '%s%s %s%s' % (prefix, name, type, suffix), 'unknown3'

    if node['node'] == 'RecordDecl':
        return '/* RecordDecl */', 'unknown5'

    if node['node'] == 'ParenExpr':
        a, b = render_expression(node['children'][0])
        return '(%s)' % a, b

    if node['node'] == 'PredefinedExpr':
        if node['name'] == '__PRETTY_FUNCTION__':
            # FIXME
            return '"void print_number(int *)"', 'const char*'

        if node['name'] == '__func__':
            # FIXME
            return '"%s"' % 'print_number', 'const char*'

        raise Exception('render_expression: unknown PredefinedExpr: %s' % node['name'])

    if node['node'] == 'ConditionalOperator':
        a = render_expression(node['children'][0])[0]
        b = render_expression(node['children'][1])[0]
        c = render_expression(node['children'][2])[0]

        add_import('github.com/elliotchance/c2go/noarch')
        return 'noarch.Ternary(%s, func () interface{} { return %s }, func () interface{} { return %s })' % (a, b, c), node['type']

    raise Exception('render_expression: %s' % node['node'])

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

def render(out, node, function_name, indent=0, return_type=None):
    if node['node'] == 'TranslationUnitDecl':
        for c in node['children']:
            render(out, c, function_name, indent, return_type)
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
                    render(out, c, function_name, indent + 1, node['type'])

            print_line(out, '}\n', indent)

        function_defs[node['name']] = (get_function_return_type(node['type']),
            [a['type'] for a in get_function_params(node)])

        return

    if node['node'] == 'CompoundStmt':
        for c in node['children']:
            render(out, c, function_name, indent, return_type)
        return

    if node['node'] == 'IfStmt':
        children = node['children']

        e = render_expression(children[0])
        print_line(out, 'if %s {' % cast(e[0], e[1], 'bool'), indent)

        render(out, children[1], function_name, indent + 1, return_type)

        if len(children) > 2:
            print_line(out, '} else {', indent)
            render(out, children[2], function_name, indent + 1, return_type)

        print_line(out, '}', indent)

        return

    if node['node'] == 'WhileStmt':
        children = node['children']

        e = render_expression(children[0])
        print_line(out, 'for %s {' % cast(e[0], e[1], 'bool'), indent)

        render(out, children[1], function_name, indent + 1, return_type)

        print_line(out, '}', indent)

        return

    if node['node'] == 'ForStmt':
        children = node['children']

        a, b, c = [render_expression(e)[0] for e in children[:3]]
        print_line(out, 'for %s; %s; %s {' % (a, b, c), indent)

        render(out, children[3], function_name, indent + 1, return_type)

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

        if 'children' in node and function_name != 'main':
            expr, type = render_expression(node['children'][0])
            r = 'return ' + cast(expr, type, 'int')

        print_line(out, r, indent)
        return

    if node['node'] in ('BinaryOperator', 'INTEGER_LITERAL', 'CallExpr'):
        print_line(out, render_expression(node)[0], indent)
        return

    if node['node'] == 'TypedefDecl':
        name = node['name'].strip()
        if name in types_already_defined:
            return

        types_already_defined.add(name)

        # FIXME: All of the logic here is just to avoid errors, it needs to be
        # fixed up.
        # if ('struct' in node['type'] or 'union' in node['type']) and :
        #     return
        node['type'] = node['type'].replace('unsigned', '')

        resolved_type = resolve_type(node['type'])

        if name == '__mbstate_t':
            add_import('github.com/elliotchance/c2go/darwin')
            resolved_type = 'darwin.C__mbstate_t'

        if name == '__darwin_ct_rune_t':
            add_import('github.com/elliotchance/c2go/darwin')
            resolved_type = 'darwin.C__darwin_ct_rune_t'

        if name in ('__builtin_va_list', '__qaddr_t', 'definition',
            '_IO_lock_t', 'va_list', 'fpos_t', '__NSConstantString',
            '__darwin_va_list', '__fsid_t', '_G_fpos_t', '_G_fpos64_t'):
            return

        print_line(out, "type %s %s\n" % (name, resolved_type), indent)

        return

    if node['node'] == 'EnumDecl':
        return

    if node['node'] == 'FieldDecl':
        print_line(out, render_expression(node)[0], indent + 1)
        return

    if node['node'] == 'RecordDecl':
        name = node['name'].strip()
        if name in types_already_defined or name == '':
            return

        types_already_defined.add(name)

        if node['kind'] == 'union':
            return

        print_line(out, "type %s %s {" % (name, node['kind']), indent)
        if 'children' in node:
            for c in node['children']:
                render(out, c, function_name, indent + 1)

        print_line(out, "}\n", indent)
        return

    if node['node'] == 'DeclStmt':
        for child in node['children']:
            print_line(out, render_expression(child)[0], indent)
        return

    if node['node'] == 'VarDecl':
        # FIXME?
        return

    if node['node'] == 'ParenExpr':
        print_line(out, render_expression(node)[0], indent)
        return

    raise Exception(node['node'])

# 1. Compile it first (checking for errors)
c_file_path = sys.argv[1]

# 2. Preprocess
pp = subprocess.Popen(["clang", "-E", c_file_path], stdout=subprocess.PIPE).communicate()[0]

pp_file_path = '/tmp/pp.c'
with open(pp_file_path, 'wb') as pp_out:
    pp_out.write(pp)

# 3. Generate JSON from AST
ast_pp = subprocess.Popen(["clang", "-Xclang", "-ast-dump", "-fsyntax-only", pp_file_path], stdout=subprocess.PIPE)
pp = subprocess.Popen(["./c2go"], stdin=ast_pp.stdout, stdout=subprocess.PIPE).communicate()[0]

json_file_path = 'pp.json'
with open(json_file_path, 'wb') as json_out:
    json_out.write(pp)

with open(json_file_path, 'r') as json_in:
    # 3. Parse C and output Go
    go_file_path = '%s.go' % c_file_path.split('/')[-1][:-2]
    go_out = io.StringIO()
    all_json = json_in.read()

    try:
        l = json.loads(all_json)
    except ValueError as e:
        # This occurs if the JSON cannot be parsed
        print(all_json)
        raise e

    render(go_out, l[0], function_name=None)

    print("package main\n")
    print("import (")
    for import_name in sorted(imports):
        print('\t"%s"' % import_name)
    print(")\n")

    print(go_out.getvalue())
