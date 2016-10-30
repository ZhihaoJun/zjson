package core

import (
	"fmt"
)

const (
	JSONTypeMap = 1 << iota
	JSONTypeArray
	JSONTypeString
	JSONTypeBool
	JSONTypeNumber
	JSONTypeNull
)

type JSON struct {
	Type   int
	Map    map[string]*JSON
	Arr    []*JSON
	Str    string
	Bool   bool
	Number int
}

func (j *JSON) String() string {
	switch j.Type {
	case JSONTypeMap:
		return fmt.Sprintf("%v", j.Map)
	case JSONTypeArray:
		return fmt.Sprintf("%v", j.Arr)
	case JSONTypeString:
		return j.Str
	case JSONTypeBool:
		return fmt.Sprintf("%v", j.Bool)
	case JSONTypeNumber:
		return fmt.Sprintf("%v", j.Number)
	}
	return ""
}

func (j *JSON) Value() interface{} {
	switch j.Type {
	case JSONTypeMap:
		return j.Map
	case JSONTypeArray:
		return j.Arr
	case JSONTypeString:
		return j.Str
	case JSONTypeBool:
		return j.Bool
	case JSONTypeNumber:
		return j.Number
	}
	return nil
}
