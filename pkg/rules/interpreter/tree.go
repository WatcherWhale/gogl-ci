package interpreter

import (
	"reflect"

	"github.com/watcherwhale/gogl-ci/pkg/rules/lexer"
)

type evaluable interface {
	GetPriority() int
	Eval(map[string]string) (any, error)
	build(pos int, evaluables []evaluable)
}

type ConditionalTree struct {
	tokens     []lexer.Token
	evaluables []evaluable
	Evaluable  evaluable
}

func (ct *ConditionalTree) Eval(input map[string]string) (any, error) {
	val, err := ct.Evaluable.Eval(input)
	if err != nil {
		return false, err
	}

	return evaluate(val), nil
}

func (ct *ConditionalTree) GetPriority() int {
	return 1
}

func (ct *ConditionalTree) build(_ int, _ []evaluable) {
}

func evaluate(stringlike any) bool {
	if stringlike == nil {
		return false
	}

	if reflect.TypeOf(stringlike).Kind() == reflect.Bool {
		return stringlike.(bool)
	}

	if stringlike.(string) == "" {
		return false
	}

	return true
}
