// Code generated by 'yaegi extract github.com/watcherwhale/gogl-ci/pkg/rules/interpreter'. DO NOT EDIT.

package symbols

import (
	"github.com/watcherwhale/gogl-ci/pkg/rules/interpreter"
	"reflect"
)

func init() {
	Symbols["github.com/watcherwhale/gogl-ci/pkg/rules/interpreter/interpreter"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"BuildConditionalTree": reflect.ValueOf(interpreter.BuildConditionalTree),
		"Evaluate":             reflect.ValueOf(interpreter.Evaluate),

		// type definitions
		"ConditionalTree": reflect.ValueOf((*interpreter.ConditionalTree)(nil)),
		"Null":            reflect.ValueOf((*interpreter.Null)(nil)),
		"Operator":        reflect.ValueOf((*interpreter.Operator)(nil)),
		"StringLike":      reflect.ValueOf((*interpreter.StringLike)(nil)),
		"Variable":        reflect.ValueOf((*interpreter.Variable)(nil)),
	}
}
