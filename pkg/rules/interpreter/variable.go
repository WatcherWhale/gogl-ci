package interpreter

import "github.com/watcherwhale/gogl-ci/pkg/rules/lexer"

type Variable struct {
	token lexer.Token
}

func (v *Variable) Eval(input map[string]string) (any, error) {
	value, ok := input[v.token.Value]

	if !ok {
		return nil, nil
	}

	return value, nil
}

func (v *Variable) GetPriority() int {
	return 1
}

func (v *Variable) build(pos int, evaluables []evaluable) {
	if len(evaluables) != 1 {
		panic("invalid amount of evaluables in a string like")
	}
}
