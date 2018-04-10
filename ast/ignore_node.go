package ast

func ignoreNode(n []string) bool {
	var ignoreMap = map[string]bool{
		"NotTailCalledAttr": true,
		"FormatArgAttr":     true,
		"...":               true,
	}
	return ignoreMap[n[0]]
}
