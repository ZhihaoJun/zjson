package core

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
		for _, lexer := range ts.lexers {
			n, token, err := lexer.Run(s[i:])
			if n > 0 {
				if err != nil {
					return nil, err
				}

				// move to next token's start
				i += n
				if token != nil {
					r = append(r, token)
				}

				// only one lexer will work
				break
			}
		}
	}
	return r, nil
}
