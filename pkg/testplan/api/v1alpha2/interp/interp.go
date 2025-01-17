package interp

import (
	"os"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/watcherwhale/gogl-ci/pkg/symbols"
)

func NewInterpreter() (*interp.Interpreter, error) {
	i := interp.New(interp.Options{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	err := i.Use(stdlib.Symbols)
	if err != nil {
		return nil, err
	}

	err = i.Use(symbols.Symbols)
	if err != nil {
		return nil, err
	}

	return i, nil
}
