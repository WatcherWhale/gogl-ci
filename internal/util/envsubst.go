package util

import (
	"fmt"
	"github.com/a8m/envsubst/parse"
)

func SubstituteEnv(input string, envMap map[string]string) (string, error) {
	return parse.New("string", mapToEnv(envMap), &parse.Restrictions{
		NoUnset: false,
		NoEmpty: false,
		NoDigit: false,
	}).Parse(input)
}

func mapToEnv(input map[string]string) []string {
	env := make([]string, 0)

	for k, v := range input {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	return env
}
