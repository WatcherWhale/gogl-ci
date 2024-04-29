package interpreter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	INPUT_MAP map[string]string = map[string]string{
		"VAR":  "hello",
		"VARS": "test",
	}

	TEST_CASES []struct {
		string
		bool
	} = []struct {
		string
		bool
	}{
		{`($VAR == "hello" || "hello" == "test")`, true},
		{`$VAR == "hello" || "hello" == "lol" && "hello" == "lol"`, false},
		{`$VAR != "hello" || ("lol" == "lol" && "lol" == "lol")`, true},
		{`$VAR == $VARS`, false},
		{`$VAR != $VARS`, true},
		{`"test" != $VARS`, false},
		{`"test" == $VARS`, true},
		{`$VAR == $VAR`, true},
		{`$VAR != $VAR`, false},
		{`(($VAR == $VAR || ($VAR == $VAR && $VAR == $VAR)) || ($VAR == $VAR || ($VAR == $VAR))) && $VAR == $VAR`, true},
		{`$DOES_NOT_EXIST == null`, true},
		{`$DOES_NOT_EXIST != null`, false},
		{`$VAR != null`, true},
		{`$VAR == null`, false},
		{`$VAR`, true},
		{`$DOES_NOT_EXIST`, false},
		{`$VAR && $VAR == "test"`, false},
		{`$VAR || $VAR == "test"`, true},
		{``, true},
	}
)

func TestEvaluation(t *testing.T) {
	for _, tc := range TEST_CASES {
		output, err := Evaluate(tc.string, INPUT_MAP)

		require.NoError(t, err)

		if tc.bool != output {
			t.Errorf("Rule '%s' evaluated to %v instead of %v", tc.string, output, tc.bool)
		}
	}
}
