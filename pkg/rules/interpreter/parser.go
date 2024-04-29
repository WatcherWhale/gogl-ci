package interpreter

import (
	"github.com/watcherwhale/gogl-ci/pkg/rules/lexer"
)

type parser struct {
	tokens []lexer.Token
	pos    int
}

func (p *parser) currentToken() lexer.Token {
	return p.tokens[p.pos]
}

func (p *parser) remainder() []lexer.Token {
	return p.tokens[p.pos:]
}

func (p *parser) advance(delta int) {
	p.pos += delta
}

func (p *parser) end() bool {
	return p.pos == len(p.tokens)-1
}

func BuildConditionalTree(tokens []lexer.Token, endToken lexer.TokenKind) (ConditionalTree, error) {
	// In the case of an empty conditional tree always return true
	if len(tokens) == 1 {
		return ConditionalTree{
			Evaluable: &StringLike{
				token: lexer.Token{
					Kind:  lexer.STRING,
					Value: "true", // This can be any non-empty string
				},
			},
		}, nil
	}

	p := parser{
		tokens: tokens,
		pos:    0,
	}

	ct := ConditionalTree{}

	for !p.end() {
		if p.currentToken().Kind == endToken {
			break
		}

		if p.currentToken().Kind == lexer.OPEN_PAREN {
			p.advance(1)
			subTree, err := BuildConditionalTree(p.remainder(), lexer.CLOSE_PAREN)
			if err != nil {
				return ConditionalTree{}, err
			}
			p.advance(len(subTree.tokens))
			ct.evaluables = append(ct.evaluables, &subTree)
		} else {
			ct.tokens = append(ct.tokens, p.currentToken())
			ct.evaluables = append(ct.evaluables, tokenToEvaluable(p.currentToken()))
			p.advance(1)
		}
	}

	maxIndex := maxEvaluable(ct.evaluables)
	ct.Evaluable = ct.evaluables[maxIndex]
	ct.Evaluable.build(maxIndex, ct.evaluables)

	ct.evaluables = []evaluable{}

	return ct, nil
}

func tokenToEvaluable(token lexer.Token) evaluable {
	switch token.Kind {
	case lexer.REGEX, lexer.STRING:
		return &StringLike{token}
	case lexer.IDENTIFIER:
		return &Variable{token}
	case lexer.NULL:
		return &Null{token}
	case lexer.AND, lexer.OR, lexer.EQUAL, lexer.NOT_EQUAL, lexer.LIKE, lexer.NOT_LIKE:
		return &Operator{token: token}
	}

	return nil
}

func maxEvaluable(evaluables []evaluable) int {
	index := -1
	value := 0
	for i, eval := range evaluables {
		if eval == nil {
			continue
		}
		if eval.GetPriority() > value {
			index = i
			value = eval.GetPriority()
		}
	}

	return index
}
