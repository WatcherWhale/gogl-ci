package interpreter

import "github.com/watcherwhale/gogl-ci/pkg/rules/lexer"

type Null struct {
	token lexer.Token
}

func (n *Null) Eval(input map[string]string) (any, error) {
	return nil, nil
}

func (n *Null) GetPriority() int {
	return 1
}

func (n *Null) build(pos int, evaluables []evaluable) {
	if len(evaluables) != 1 {
		panic("invalid amount of evaluables in a null")
	}
}
