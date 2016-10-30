package core

import (
	"fmt"
)

const (
	TokenTypeLeftBracket = 1 << iota
	TokenTypeRightBracket
	TokenTypeLeftSquareBracket
	TokenTypeRightSquareBracket
	TokenTypeComma
	TokenTypeColon
	TokenTypeString
	TokenTypeBool
	TokenTypeNull
	TokenTypeNumber
)

var tokenTypeToStr = map[int]string{
	TokenTypeLeftBracket:        "{",
	TokenTypeRightBracket:       "}",
	TokenTypeLeftSquareBracket:  "[",
	TokenTypeRightSquareBracket: "]",
	TokenTypeComma:              ",",
	TokenTypeColon:              ":",
	TokenTypeString:             "string",
	TokenTypeBool:               "bool",
	TokenTypeNull:               "null",
	TokenTypeNumber:             "number",
}

type Token struct {
	Type int
	Buf  []byte
	Bool bool
	Num  int
}

func (t *Token) String() string {
	return fmt.Sprintf("type: %s %s %v %d\n", tokenTypeToStr[t.Type], string(t.Buf), t.Bool, t.Num)
}

func (t *Token) ToString() string {
	return string(t.Buf)
}

func (t *Token) ToBool() bool {
	return t.Bool
}

func (t *Token) ToNumber() int {
	return t.Num
}

func (t *Token) ToInterface() interface{} {
	switch t.Type {
	case TokenTypeBool:
		return t.Bool
	case TokenTypeString:
		return t.ToString()
	case TokenTypeNumber:
		return t.ToNumber()
	case TokenTypeNull:
		return nil
	}
	return nil
}
