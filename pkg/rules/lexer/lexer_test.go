package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComplex(t *testing.T) {
	testRule := `($VAR == "test" || $VAR != "test2") && ( $VAR_WITH_UDERSCORE =~ /regex/ || $VAR_WITH_UDERSCORE !~ /regex\/with\/a\\\/slash/ )`

	tokens, err := Tokenize(testRule)

	require.NoError(t, err)

	kinds := make([]TokenKind, len(tokens))
	for i, token := range tokens {
		kinds[i] = token.Kind
	}

	assert.EqualValues(t, []TokenKind{
		OPEN_PAREN,
		IDENTIFIER,
		EQUAL,
		STRING,
		OR,
		IDENTIFIER,
		NOT_EQUAL,
		STRING,
		CLOSE_PAREN,

		AND,

		OPEN_PAREN,
		IDENTIFIER,
		LIKE,
		REGEX,
		OR,
		IDENTIFIER,
		NOT_LIKE,
		REGEX,
		CLOSE_PAREN,

		EOF,
	}, kinds)

	assert.Equal(t, "VAR", tokens[1].Value)
	assert.Equal(t, "test", tokens[3].Value)
	assert.Equal(t, "VAR", tokens[5].Value)
	assert.Equal(t, "test2", tokens[7].Value)
	assert.Equal(t, "VAR_WITH_UDERSCORE", tokens[11].Value)
	assert.Equal(t, "regex", tokens[13].Value)
	assert.Equal(t, "VAR_WITH_UDERSCORE", tokens[15].Value)
	assert.Equal(t, `regex\/with\/a\\\/slash`, tokens[17].Value)
}

func TestSimpleRegex(t *testing.T) {
	testRule := `$CI_COMMIT_BRANCH != null && $CI_COMMIT_BRANCH =~ /^\d+\.\d+\.x$/`
	tokens, err := Tokenize(testRule)

	require.NoError(t, err)

	kinds := make([]TokenKind, len(tokens))
	for i, token := range tokens {
		kinds[i] = token.Kind
	}

	assert.EqualValues(t, []TokenKind{
		IDENTIFIER,
		NOT_EQUAL,
		NULL,
		AND,
		IDENTIFIER,
		LIKE,
		REGEX,
		EOF,
	}, kinds)
}

func TestSyntaxErrors(t *testing.T) {
	testCase := `$VAR * 2`
	tokens, err := Tokenize(testCase)

	assert.Nil(t, tokens)

	require.EqualError(t, err, "syntax error: cannot find token for '* 2'")
}

func TestEmptyString(t *testing.T) {
	testCase := `$TEST_VAR == ""`
	tokens, err := Tokenize(testCase)

	require.NoError(t, err)

	kinds := make([]TokenKind, len(tokens))
	for i, token := range tokens {
		kinds[i] = token.Kind
	}

	assert.EqualValues(t, []TokenKind{
		IDENTIFIER,
		EQUAL,
		STRING,
		EOF,
	}, kinds)
}

func TestEscapedString(t *testing.T) {
	testCase := `$TEST_VAR == "this is a long string containing escaped \" and a substring \"substring\"!!!"`
	tokens, err := Tokenize(testCase)

	require.NoError(t, err)

	kinds := make([]TokenKind, len(tokens))
	for i, token := range tokens {
		kinds[i] = token.Kind
	}

	assert.EqualValues(t, []TokenKind{
		IDENTIFIER,
		EQUAL,
		STRING,
		EOF,
	}, kinds)
}
