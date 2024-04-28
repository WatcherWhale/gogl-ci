package lexer

import (
	"fmt"
	"regexp"
)

type regexpHandler func(lex *lexer, regex *regexp.Regexp)

type regexPatterns struct {
	regex   *regexp.Regexp
	handler regexpHandler
}

type lexer struct {
	patterns []regexPatterns
	Tokens   []Token
	source   string
	pos      int
}

func (lex *lexer) adavance(delta int) {
	lex.pos += delta
}

func (lex *lexer) push(token Token) {
	lex.Tokens = append(lex.Tokens, token)
}

func (lex *lexer) remainer() string {
	return lex.source[lex.pos:]
}

func (lex *lexer) at_eof() bool {
	return lex.pos >= len(lex.source)
}

func Tokenize(source string) ([]Token, error) {
	lex := createLexer(source)

	for !lex.at_eof() {
		matched := false

		for _, pattern := range lex.patterns {
			loc := pattern.regex.FindStringIndex(lex.remainer())

			if loc != nil && loc[0] == 0 {
				pattern.handler(lex, pattern.regex)
				matched = true
				break
			}
		}

		if !matched {
			return nil, fmt.Errorf("syntax error: cannot find token for '%s'", lex.remainer())
		}
	}

	lex.push(Token{EOF, "EOF"})

	return lex.Tokens, nil
}

func createLexer(source string) *lexer {
	return &lexer{
		pos:    0,
		source: source,
		Tokens: make([]Token, 0),
		patterns: []regexPatterns{
			{regexp.MustCompile(`\s+`), skipHandler},
			{regexp.MustCompile(`\(`), defaultHandler(OPEN_PAREN, "(")},
			{regexp.MustCompile(`\)`), defaultHandler(CLOSE_PAREN, ")")},
			{regexp.MustCompile(`==`), defaultHandler(EQUAL, "==")},
			{regexp.MustCompile(`!=`), defaultHandler(NOT_EQUAL, "!=")},
			{regexp.MustCompile(`=~`), defaultHandler(LIKE, "=~")},
			{regexp.MustCompile(`!~`), defaultHandler(NOT_LIKE, "!~")},
			{regexp.MustCompile(`\|\|`), defaultHandler(OR, "||")},
			{regexp.MustCompile(`&&`), defaultHandler(AND, "&&")},
			{regexp.MustCompile(`null`), defaultHandler(NULL, "null")},
			{regexp.MustCompile(`\$([A-Za-z0-9_]+)`), capturedHandler(IDENTIFIER)},
			{regexp.MustCompile(`"([^"]+)"`), capturedHandler(STRING)},
			{regexp.MustCompile(`\/(([^\/\\]+|\\\/|\\\\)+)\/`), capturedHandler(REGEX)},
		},
	}
}

func defaultHandler(kind TokenKind, value string) regexpHandler {
	return func(lex *lexer, regex *regexp.Regexp) {
		lex.adavance(len(value))
		lex.push(Token{
			Kind:  kind,
			Value: value,
		})
	}
}

func skipHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainer())
	lex.adavance(len(match))
}

func capturedHandler(kind TokenKind) regexpHandler {
	return func(lex *lexer, regex *regexp.Regexp) {
		fullmatch := regex.FindString(lex.remainer())
		match := regex.FindStringSubmatch(lex.remainer())
		lex.push(Token{
			Kind:  kind,
			Value: match[1],
		})

		lex.adavance(len(fullmatch))
	}
}
