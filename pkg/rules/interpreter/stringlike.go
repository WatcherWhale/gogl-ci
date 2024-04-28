package interpreter

import "github.com/watcherwhale/gogl-ci/pkg/rules/lexer"

type StringLike struct {
	token lexer.Token
}

func (s *StringLike) Eval(input map[string]string) (any, error) {
	return s.token.Value, nil
}

func (s *StringLike) GetPriority() int {
	return 1
}

func (s *StringLike) build(pos int, evaluables []evaluable) {
	if len(evaluables) != 1 {
		panic("invalid amount of evaluables in a string like")
	}
}
