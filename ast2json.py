import sys
import re
import json

# This script converts the output of clang AST into a JSON file.
# 
# Usage:
#   clang -Xclang -ast-dump -fsyntax-only myfile.c | python ast2json.py
# 
# Yes, there are many better ways to do this. However I chose this method
# because:
# 
# 1. I need to separate the clang AST from the c2go conversion process so that
#    the c2go program can ingest a reliable JSON file and not depend on clang or
#    its different versions at all.
# 2. The clang API is not stable and trying to match up binaries with different
#    versions and operating systems can be tricky and brittle.
# 3. This tool, in time, will become a better binary of some kind that produces
#    much the same JSON output (so minimal changes to c2go.py).
# 4. I needed something quick and dirty to proof the complete toolchain and get
#    it working on different versions of clang and different operating systems
#    before we enough information to really standardise the process.

regex = {
    'ConstAttr': r"^ (?P<address>[0-9a-fx]+) <(?P<position>.*)>(?P<tags>.*)",
    'Enum': r"^ (?P<address>[0-9a-fx]+) '(?P<name>.*)'",
    'EnumConstantDecl': r"^ (?P<address>[0-9a-fx]+) <(?P<position>.*)>(?P<position2> [^ ]+)? (?P<name>.+) '(?P<type>.+?)'",
    'EnumDecl': r"^ (?P<address>[0-9a-fx]+) <(?P<position>.*)>(?P<position2> [^ ]+)?(?P<name>.*)",
    'EnumType': r"^ (?P<address>[0-9a-fx]+) '(?P<name>.*)'",

    'FormatAttr': r'^ (?P<address>[0-9a-fx]+) <(?P<position>.*)>(?P<tags> \w+)? (?P<function>\w+) (?P<unknown1>\d+) (?P<unknown2>\d+)',
    'ForStmt': r"^ (?P<address>[0-9a-fx]+) <(?P<position>.*)>",
    'FunctionDecl': r"^ (?P<address>[0-9a-fx]+) (?P<prev>prev [0-9a-fx]+)? ?<(?P<position1>.*)>(?P<position2> [^ ]+)?(?P<tags> .*)? (?P<name>\w+) '(?P<type>.*)'(?P<tags3> extern)?",
    'FunctionProtoType': r"^ (?P<address>[0-9a-fx]+) \'(?P<type>.*)\' (?P<kind>.*)",
    'IfStmt': r'^ (?P<address>[0-9a-fx]+) <(?P<position>.*)>',
    'ImplicitCastExpr': r'^ (?P<address>[0-9a-fx]+) <(?P<position>.*)> \'(?P<type>.*)\' <(?P<kind>.*)>',
    'IntegerLiteral': r'^ (?P<address>[0-9a-fx]+) <(?P<position>.*)> \'(?P<type>.*)\' (?P<value>.+)',
    'MallocAttr': r"^ (?P<address>[0-9a-fx]+) <(?P<position>.*)>",
    'MemberExpr': r"^ (?P<address>[0-9a-fx]+) <(?P<position>.*)> '(?P<type>.*?)' (?P<tags>.*?)(?P<name>\w+) (?P<address2>[0-9a-fx]+)",
    'ModeAttr': r"^ (?P<address>[0-9a-fx]+) <(?P<position>.*)> (?P<name>.+)",
    'NonNullAttr': r"^ (?P<address>[0-9a-fx]+) <(?P<position>.*)> 1",
    'NoThrowAttr': r"^ (?P<address>[0-9a-fx]+) <(?P<position>.*)>",
    'ParenExpr': r'^ (?P<address>[0-9a-fx]+) <(?P<position>.*)> \'(?P<type>.*?)\'',
    'ParmVarDecl': r"^ (?P<address>[0-9a-fx]+) <(?P<position>.*)>(?P<position2> [^ ]+:[\d:]+)?(?P<used> used)?(?P<name> \w+)? '(?P<type>.*?)'(?P<type2>:'.*?')?",
    'PointerType': r'^ (?P<address>[0-9a-fx]+) \'(?P<type>.*)\'',
    'PredefinedExpr': r"^ (?P<address>[0-9a-fx]+) <(?P<position>.*)> '(?P<type>.*)' (?P<kind>.*) (?P<name>.*)",
    'QualType': r"^ (?P<address>[0-9a-fx]+) \'(?P<type>.*)\' (?P<kind>.*)",
    'Record': r'^ (?P<address>[0-9a-fx]+) \'(?P<type>.*)\'',
    'RecordDecl': r"^ (?P<address>[0-9a-fx]+) (?P<prev>prev 0x[0-9a-f]+ )?<(?P<position>.*)> (?P<position2>[^ ]+ )?(?P<kind>struct|union) (?P<name>\w*)( definition)?",
    'RecordType': r'^ (?P<address>[0-9a-fx]+) \'(?P<type>.*)\'',
    'RestrictAttr': r"^ (?P<address>[0-9a-fx]+) <(?P<position>.*)> (?P<name>.*)",
    'ReturnStmt': r'^ (?P<address>[0-9a-fx]+) <(?P<position>.*)>',
    'StringLiteral': r'^ (?P<address>[0-9a-fx]+) <(?P<position>.*)> \'(?P<type>.*)\'(?P<tags> lvalue)? (?P<value>.*)',
    'TranslationUnitDecl': r'^ (?P<address>[0-9a-fx]+)',
    'Typedef': r'^ (?P<address>[0-9a-fx]+) \'(?P<type>.*)\'',
    'TypedefDecl': r"(?P<address>[0-9a-fx]+) <(?P<position>.+?)> (?P<position2><invalid sloc>|0x[0-9a-f]+)?(?P<tags>.*?)(?P<name>\w+) '(?P<type>.*?)'(?P<type2>:'.*?')?",
    'TypedefType': r'^ (?P<address>[0-9a-fx]+) \'(?P<type>.*)\' (?P<tags>.+)',
    'UnaryOperator': r"^ (?P<address>[0-9a-fx]+) <(?P<position>.*)> '(?P<type>.*?)'(?P<tags1> lvalue)?(?P<tags2> prefix)?(?P<tags3> postfix)? '(?P<operator>.*?)'",
    'VarDecl': r"^ (?P<address>[0-9a-fx]+) <(?P<position>.*)>(?P<position2> [^ ]+)? (?P<name>.+) '(?P<type>.+?)'(?P<type2>:'.*?')?(?P<tags>.*)",
    'WhileStmt': r"^ (?P<address>[0-9a-fx]+) <(?P<position>.*)>",
}

# ParmVarDecl 0x4167750 <line:56:23> line:493:15 'struct __va_list_tag *':'struct __va_list_tag *'
# ParmVarDecl 0x2839dd0 </usr/include/_G_config.h:32:20> /usr/include/libio.h:496:58 '__ssize_t':'long'

def build_tree(nodes, depth):
    """Convert an array of nodes, each prefixed with a depth into a tree."""
    if len(nodes) == 0:
        return []

    # Split the list into sections, treat each section as a a tree with its own
    # root.
    sections = []
    for node in nodes:
        if node[0] == depth:
            sections.append([node])
        else:
            sections[-1].append(node)

    results = []
    for section in sections:
        children = build_tree([n for n in section if n[0] > depth], depth + 1)
        result = section[0][1]

        if len(children) > 0:
            result['children'] = children

        results.append(result)

    return results

def read_ast():
    stdin = sys.stdin.read()
    uncolored = re.sub(r'\x1b\[[\d;]+m', '', stdin)
    return uncolored.split("\n")

def convert_lines_to_nodes(lines):
    nodes = []
    for line in lines:
        if line.strip() == '':
            continue

        # This will need to be handled more gracefully...  I'm not even sure
        # what this means?
        if '<<<NULL>>>' in line:
            continue

        indent_and_type = re.search(r'^([|\- `]*)(\w+)', line)
        if indent_and_type is None:
            print("Can not understand line '%s'" % line)
            sys.exit(1)

        node_type = indent_and_type.group(2)
        # if node_type == 'VarDecl':
        #     print(line[offset:])

        offset = len(indent_and_type.group(0))
        try:
            result = re.search(regex[node_type], line[offset:])
        except KeyError:
            print("There is no regex for '%s'." % node_type)
            print("I will print out all the lines so a regex can be created:\n")

            for line in lines:
                s = re.search(r'^([|\- `]*)(\w+)', line)
                if s is not None and node_type == s.group(2):
                    print(line[offset:])

            sys.exit(1)

        if result is None:
            print("Can not understand line '%s'" % line)
            sys.exit(1)

        node = result.groupdict()

        node['node'] = node_type

        indent_level = len(indent_and_type.group(1)) / 2
        nodes.append([indent_level, node])

    return nodes

lines = read_ast()
nodes = convert_lines_to_nodes(lines)
tree = build_tree(nodes, 0)

print(json.dumps(tree, sort_keys=True, indent=2, separators=(',', ': ')))
