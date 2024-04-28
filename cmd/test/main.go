package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/watcherwhale/gogl-ci/pkg/rules/interpreter"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	examples := []string{
		`($VAR == "hello" || "hello" == "lol")`,
		`$VAR == "hello" || "hello" == "lol" && "hello" == "lol"`,
		`$VAR != "hello" || ("lol" == "lol" && "lol" == "lol")`,
		`$VAR == $VARS`,
		`$VAR != $VARS`,
		`"test" != $VARS`,
		`"test" == $VARS`,
		`$VAR == $VAR`,
		`$VAR != $VAR`,
		`(($VAR == $VAR || ($VAR == $VAR && $VAR == $VAR)) || ($VAR == $VAR || ($VAR == $VAR))) && $VAR == $VAR`,
		`$DOES_NOT_EXIST == null`,
		`$DOES_NOT_EXIST != null`,
		`$VAR != null`,
		`$VAR == null`,
		`$VAR`,
		`$DOES_NOT_EXIST`,
		`$VAR && $VAR == "test"`,
		`$VAR || $VAR == "test"`,
	}

	for _, ex := range examples {
		log.Info().Msg(ex)
		out, err := interpreter.Evaluate(ex, map[string]string{
			"VAR":  "hello",
			"VARS": "test",
		})
		if err != nil {
			fmt.Printf("%v\n", err)
		}

		fmt.Printf("%v\n", out)
	}
}
