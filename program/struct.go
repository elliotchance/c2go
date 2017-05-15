package program

// Struct represents the definition for a C struct.
type Struct struct {
	// The name of the struct.
	Name string

	// Each of the fields and their C type.
	Fields map[string]string
}
