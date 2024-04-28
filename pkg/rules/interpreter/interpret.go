package interpreter

import (
	"github.com/watcherwhale/gogl-ci/pkg/rules/lexer"
)

func Evaluate(rule string, input map[string]string) (bool, error) {
	tokens, err := lexer.Tokenize(rule)
	if err != nil {
		return false, nil
	}

	ct, err := BuildConditionalTree(tokens, lexer.EOF)
	if err != nil {
		return false, nil
	}

	value, err := ct.Eval(input)
	if err != nil {
		return false, nil
	}

	return value.(bool), nil
}
