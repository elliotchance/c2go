package types

import (
    "errors"
    "fmt"
    "strconv"
    "strings"

    "github.com/elliotchance/c2go/program"
    "github.com/elliotchance/c2go/util"
)

func removePrefix(s, prefix string) string {
    if strings.HasPrefix(s, prefix) {
        s = s[len(prefix):]
    }

    return s
}

// SizeOf returns the number of bytes for a type. This the same as using the
// sizeof operator/function in C.
func SizeOf(p *program.Program, cType string) (int, error) {
    // Remove keywords that do not effect the size.
    cType = removePrefix(cType, "signed ")
    cType = removePrefix(cType, "unsigned ")
    cType = removePrefix(cType, "const ")
    cType = removePrefix(cType, "volatile ")

    // FIXME: The pointer size will be different on different platforms. We
    // should find out the correct size at runtime.
    pointerSize := 8

    // A structure will be the sum of its parts.
    if strings.HasPrefix(cType, "struct ") {
        totalBytes := 0

        s := p.Structs[cType[7:]]
        if s == nil {
            return 0, errors.New(fmt.Sprintf("could not sizeof: %s", cType))
        }

        for _, t := range s.Fields {
            var bytes int
            var err error

            switch f := t.(type) {
            case string:
                bytes, err = SizeOf(p, f)

            case *program.Struct:
                bytes, err = SizeOf(p, f.Name)
            }

            if err != nil {
                return 0, err
            }
            totalBytes += bytes
        }

        // The size of a struct is rounded up to fit the size of the pointer of
        // the OS.
        if totalBytes%pointerSize != 0 {
            totalBytes += pointerSize - (totalBytes % pointerSize)
        }

        return totalBytes, nil
    }

    // An union will be the max size of its parts.
    if strings.HasPrefix(cType, "union ") {
        byte_count := 0

        s := p.Unions[cType[6:]]
        if s == nil {
            return 0, errors.New(fmt.Sprintf("could not sizeof: %s", cType))
        }

        for _, t := range s.Fields {
            var bytes int
            var err error

            switch f := t.(type) {
            case string:
                bytes, err = SizeOf(p, f)

            case *program.Struct:
                bytes, err = SizeOf(p, f.Name)
            }

            if err != nil {
                return 0, err
            }

            if byte_count < bytes {
                byte_count = bytes
            }
        }

        // The size of an union is rounded up to fit the size of the pointer of
        // the OS.
        if byte_count%pointerSize != 0 {
            byte_count += pointerSize - (byte_count % pointerSize)
        }

        return byte_count, nil
    }

    // Function pointers are one byte?
    if strings.Index(cType, "(") >= 0 {
        return 1, nil
    }

    if strings.HasSuffix(cType, "*") {
        return pointerSize, nil
    }

    switch cType {
    case "char", "void":
        return 1, nil

    case "short":
        return 2, nil

    case "int", "float":
        return 4, nil

    case "long", "double":
        return 8, nil

    case "long double":
        return 16, nil
    }

    // Get size for array types like: `base_type [count]`
    groups := util.GroupsFromRegex(`^(?P<type>.+) ?[(?P<count>\d+)\]$`, cType)
    fmt.Println("Gr:", groups)

    if groups == nil {
        return pointerSize, errors.New(
            fmt.Sprintf("cannot determine size of: %s", cType))
    }

    base_size, err := SizeOf(p, groups["type"])
    if err != nil {
        return 0, err
    }

    count, err := strconv.Atoi(groups["count"])
    if err != nil {
        return 0, err
    }

    return base_size * count, nil
}
