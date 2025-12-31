// Package boost
// File:        boost.go
// Author:      TRAE AI
// Created:     2025/12/30 11:03:46
// Description: Global instance declarations for go-boost library
package boost

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
type (
	BOOST struct {
		value interface{}
	}
)

var (
	Debugger = NewDebugger()
)

func M(value interface{}) BOOST {
	return BOOST{value: value}
}

func (b BOOST) AsDirectory() *DIRECTORY {
	path := ""
	if s, ok := b.value.(string); ok {
		path = s
	}
	return &DIRECTORY{path: path}
}

func (b BOOST) AsFile() *FILE {
	path := ""
	if s, ok := b.value.(string); ok {
		path = s
	}
	return &FILE{path: path}
}

func (b BOOST) AsFilePath() *FILEPATH {
	path := ""
	if s, ok := b.value.(string); ok {
		path = s
	}
	return &FILEPATH{_filepath: path}
}

func (b BOOST) AsJson() *JSON {
	return NewJSON(b.value)
}

func (b BOOST) AsProcess() *PROCESS {
	pid := 0
	if pidInt, ok := b.value.(int); ok {
		pid = pidInt
	}
	return NewProcess(pid)
}
