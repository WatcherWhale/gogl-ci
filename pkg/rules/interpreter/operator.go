package interpreter

import (
	"fmt"
	"regexp"

	"github.com/watcherwhale/gogl-ci/pkg/rules/lexer"
)

type Operator struct {
	token lexer.Token
	left  evaluable
	right evaluable
}

func (o *Operator) Eval(input map[string]string) (any, error) {
	lVal, err := o.left.Eval(input)
	if err != nil {
		return nil, err
	}

	rVal, err := o.right.Eval(input)
	if err != nil {
		return nil, err
	}

	switch o.token.Kind {
	case lexer.AND:
		return evaluate(lVal) && evaluate(rVal), nil
	case lexer.OR:
		return evaluate(lVal) || evaluate(rVal), nil
	case lexer.EQUAL:
		return lVal == rVal, nil
	case lexer.NOT_EQUAL:
		return lVal != rVal, nil
	case lexer.LIKE:
		if lVal == nil {
			return false, nil
		}

		regex, err := regexp.Compile(rVal.(string))
		if err != nil {
			return nil, err
		}

		return regex.Match([]byte(lVal.(string))), nil
	case lexer.NOT_LIKE:
		if lVal == nil {
			return true, nil
		}

		regex, err := regexp.Compile(rVal.(string))
		if err != nil {
			return nil, err
		}

		return !regex.Match([]byte(lVal.(string))), nil
	}

	return nil, fmt.Errorf("invalid operator %s", o.token.Value)
}

func (o *Operator) build(pos int, evaluables []evaluable) {
	lEvs := evaluables[:pos]
	rEvs := evaluables[pos+1:]

	lev := maxEvaluable(lEvs)
	rev := maxEvaluable(rEvs)

	o.left = lEvs[lev]
	o.left.build(lev, lEvs)

	o.right = rEvs[rev]
	o.right.build(rev, rEvs)
}

func (o *Operator) GetPriority() int {
	switch o.token.Kind {
	case lexer.AND:
		return 4
	case lexer.OR:
		return 3
	case lexer.EQUAL, lexer.NOT_EQUAL, lexer.LIKE, lexer.NOT_LIKE:
		return 2
	}

	return -1
}
