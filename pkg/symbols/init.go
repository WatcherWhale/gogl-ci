package symbols

import "reflect"

//go:generate yaegi extract github.com/watcherwhale/gogl-ci/pkg/gitlab
//go:generate yaegi extract github.com/watcherwhale/gogl-ci/pkg/graph
//go:generate yaegi extract github.com/watcherwhale/gogl-ci/pkg/rules/interpreter

var Symbols = map[string]map[string]reflect.Value{}
