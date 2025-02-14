package v1alpha2

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"regexp"

	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
	"github.com/watcherwhale/gogl-ci/pkg/testplan/api/v1alpha2/interp"
)

var packageRegexp = regexp.MustCompile(`(?m)^package +([^\s]+) *$`)
var testFuncRegexp = regexp.MustCompile(`^Test.*$`)

type TestFunc func(gitlab.Pipeline) (bool, string)

func loadTestFile(path string) (map[string]TestFunc, error) {
	ctx := context.Background()

	pluginSrc, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	int, err := interp.NewInterpreter()
	if err != nil {
		return nil, err
	}

	_, err = int.EvalWithContext(ctx, string(pluginSrc))
	if err != nil {
		return nil, err
	}

	packageMatch := packageRegexp.FindStringSubmatch(string(pluginSrc))
	if len(packageMatch) != 2 {
		return nil, fmt.Errorf("cannot find package")
	}

	packageName := packageMatch[1]

	testFuncs := map[string]TestFunc{}

	for name, val := range int.Symbols(packageName)[packageName] {
		if !testFuncRegexp.MatchString(name) {
			continue
		}

		if val.Kind() == reflect.Func {
			fn, ok := val.Interface().(func(gitlab.Pipeline) (bool, string))
			if ok {
				testFuncs[name] = fn
			}
		}
	}

	return testFuncs, nil
}
