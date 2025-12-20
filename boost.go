// --------------------------------------------------------------------------------
// File:        main.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Global instance declarations for go-boost library
// --------------------------------------------------------------------------------

package boost

var (
	// Debugger is the global debugger instance
	Debugger = NewDebugger()
)

// BOOST is a utility type for converting strings to various types
type BOOST struct {
	_expr interface{}
}

// W converts an interface to a BOOST instance
// _expr: Expression to convert
// Returns: BOOST instance with the given expression
// Usage:
// boost := boost.W("/path/to/file.txt")
func W(expr interface{}) BOOST {
	return BOOST{_expr: expr}
}

// AsDirectory converts the BOOST instance to a DIRECTORY instance
// Returns: DIRECTORY instance with the given expression
// Usage:
// dir := boost.Boost("/path/to/dir").AsDirectory()
func (b BOOST) AsDirectory() DIRECTORY {
	path := ""
	if exprStr, ok := b._expr.(string); ok {
		path = exprStr
	}
	return DIRECTORY{path: path}
}

// AsFile converts the BOOST instance to a FILE instance
// Returns: FILE instance with the given expression
// Usage:
// file := boost.Boost("/path/to/file.txt").AsFile()
func (b BOOST) AsFile() FILE {
	path := ""
	if exprStr, ok := b._expr.(string); ok {
		path = exprStr
	}
	return FILE{path: path}
}

// AsFilePath converts the BOOST instance to a FILEPATH instance
// Returns: FILEPATH instance with the given expression
// Usage:
// filePath := boost.Boost("/path/to/file.txt").AsFilePath()
func (b BOOST) AsFilePath() FILEPATH {
	path := ""
	if exprStr, ok := b._expr.(string); ok {
		path = exprStr
	}
	return FILEPATH{_filepath: path}
}

// AsJson converts the BOOST instance to a JSON instance
// Returns: JSON instance with the given expression
// Usage:
// jsonInstance := boost.W(`{"key":"value"}`).AsJson()
// returns JSON instance with parsed JSON value
func (b BOOST) AsJson() JSON {
	return *NewJSON(b._expr)
}
