package core

import (
	"bytes"
)

var (
	trueBytesPrefix  = []byte{'t', 'r', 'u', 'e'}
	falseBytesPrefix = []byte{'f', 'a', 'l', 's', 'e'}
	nullBytesPrefix  = []byte{'n', 'u', 'l', 'l'}
)

// tokens are
// { } [ ] : , string number bool null
type LeftBracketLexer struct{}

func (lbl *LeftBracketLexer) Run(s []byte) (int, *Token, error) {
	if s[0] == '{' {
		token := &Token{
			Type: TokenTypeLeftBracket,
			Buf:  s[0:1],
		}
		return 1, token, nil
	}
	return 0, nil, nil
}

type RightBracketLexer struct{}

func (rbl *RightBracketLexer) Run(s []byte) (int, *Token, error) {
	if s[0] == '}' {
		token := &Token{
			Type: TokenTypeRightBracket,
			Buf:  s[0:1],
		}
		return 1, token, nil
	}
	return 0, nil, nil
}

type LeftSquareBracketLexer struct{}

func (lsbl *LeftSquareBracketLexer) Run(s []byte) (int, *Token, error) {
	if s[0] == '[' {
		token := &Token{
			Type: TokenTypeLeftSquareBracket,
			Buf:  s[0:1],
		}
		return 1, token, nil
	}
	return 0, nil, nil
}

type RightSquareBracketLexer struct{}

func (rsbl *RightSquareBracketLexer) Run(s []byte) (int, *Token, error) {
	if s[0] == ']' {
		token := &Token{
			Type: TokenTypeRightSquareBracket,
			Buf:  s[0:1],
		}
		return 1, token, nil
	}
	return 0, nil, nil
}

type CommaLexer struct{}

func (cl *CommaLexer) Run(s []byte) (int, *Token, error) {
	if s[0] == ',' {
		token := &Token{
			Type: TokenTypeComma,
			Buf:  s[0:1],
		}
		return 1, token, nil
	}
	return 0, nil, nil
}

type ColonLexer struct{}

func (cl *ColonLexer) Run(s []byte) (int, *Token, error) {
	if s[0] == ':' {
		token := &Token{
			Type: TokenTypeColon,
			Buf:  s[0:1],
		}
		return 1, token, nil
	}
	return 0, nil, nil
}

type StringLexer struct {
}

func (sl *StringLexer) Run(s []byte) (int, *Token, error) {
	if s[0] != '"' {
		return 0, nil, nil
	}
	i := 1
	for i < len(s) && s[i] != '"' {
		i++
	}
	// if i == len(s) raise error
	token := &Token{
		Type: TokenTypeString,
		Buf:  s[1:i],
	}
	return i + 1, token, nil
}

type BoolLexer struct {
}

func (bl *BoolLexer) Run(s []byte) (int, *Token, error) {
	if s[0] != 't' && s[0] != 'f' {
		return 0, nil, nil
	}
	token := &Token{
		Type: TokenTypeBool,
	}
	if bytes.HasPrefix(s, trueBytesPrefix) {
		token.Bool = true
		return 4, token, nil
	} else if bytes.HasPrefix(s, falseBytesPrefix) {
		token.Bool = false
		return 5, token, nil
	}
	return 0, token, nil
}

type NullLexer struct {
}

func (nl *NullLexer) Run(s []byte) (int, *Token, error) {
	if s[0] != 'n' {
		return 0, nil, nil
	}
	if bytes.HasPrefix(s, nullBytesPrefix) == false {
		// raise error
	}
	token := &Token{
		Type: TokenTypeNull,
	}
	return 4, token, nil
}

type NumberLexer struct {
}

func (nl *NumberLexer) Run(s []byte) (int, *Token, error) {
	if s[0] < '0' || s[0] > '9' {
		return 0, nil, nil
	}
	i := 0
	v := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		v = v*10 + int(s[i]-'0')
		i++
	}
	token := &Token{
		Type: TokenTypeNumber,
		Num:  v,
	}
	return i, token, nil
}

type BlankLexer struct {
}

func (bl *BlankLexer) Run(s []byte) (int, *Token, error) {
	if s[0] == ' ' || s[0] == '\n' || s[0] == '\t' {
		return 1, nil, nil
	}
	return 0, nil, nil
}

func NewJSONTokenSpliter() *TokenSpliter {
	return NewTokenSpliter(
		&LeftBracketLexer{},
		&RightBracketLexer{},
		&LeftSquareBracketLexer{},
		&RightSquareBracketLexer{},
		&CommaLexer{},
		&ColonLexer{},
		&StringLexer{},
		&BoolLexer{},
		&NumberLexer{},
		&NullLexer{},
		&BlankLexer{},
	)
}
