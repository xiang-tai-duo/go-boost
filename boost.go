// Package boost
// File:        boost.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Global instance declarations for go-boost library
package boost

type (
	BOOST_WRAPPER struct {
		_expr interface{}
	}
)

var (
	Debugger = NewDebugger()
)

func Boost(expr interface{}) BOOST_WRAPPER {
	return BOOST_WRAPPER{_expr: expr}
}

func (b BOOST_WRAPPER) AsDirectory() DIRECTORY {
	path := ""
	if exprStr, ok := b._expr.(string); ok {
		path = exprStr
	}
	return DIRECTORY{path: path}
}

func (b BOOST_WRAPPER) AsFile() FILE {
	path := ""
	if exprStr, ok := b._expr.(string); ok {
		path = exprStr
	}
	return FILE{path: path}
}

func (b BOOST_WRAPPER) AsFilePath() FILEPATH {
	path := ""
	if exprStr, ok := b._expr.(string); ok {
		path = exprStr
	}
	return FILEPATH{_filepath: path}
}

func (b BOOST_WRAPPER) AsJson() JSON {
	return *NewJSON(b._expr)
}
