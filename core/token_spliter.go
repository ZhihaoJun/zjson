package core

import (
	"log"
)

type Lexer interface {
	Run([]byte) (int, *Token, error)
}

type TokenSpliter struct {
	lexers []Lexer
}

func NewTokenSpliter(lexers ...Lexer) *TokenSpliter {
	return &TokenSpliter{
		lexers: lexers,
	}
}

func (ts *TokenSpliter) Run(s []byte) ([]*Token, error) {
	r := []*Token{}
	for i := 0; i < len(s); {
		log.Println("lex start at ", i)
		log.Println(string(s[i:]))
		for _, lexer := range ts.lexers {
			n, token, err := lexer.Run(s[i:])
			if n > 0 {
				if err != nil {
					return nil, err
				}
				i += n
				if token != nil {
					r = append(r, token)
				}
				break
			}
		}
	}
	return r, nil
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
