package lexer

type TokenKind int

const (
	EOF TokenKind = iota
	STRING
	REGEX
	IDENTIFIER
	NULL

	OPEN_PAREN
	CLOSE_PAREN

	EQUAL
	NOT_EQUAL
	LIKE
	NOT_LIKE
	OR
	AND
)

type Token struct {
	Kind  TokenKind
	Value string
}
