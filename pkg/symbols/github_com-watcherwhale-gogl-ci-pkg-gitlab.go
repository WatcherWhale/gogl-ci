// Code generated by 'yaegi extract github.com/watcherwhale/gogl-ci/pkg/gitlab'. DO NOT EDIT.

package symbols

import (
	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
	"go/constant"
	"go/token"
	"reflect"
)

func init() {
	Symbols["github.com/watcherwhale/gogl-ci/pkg/gitlab/gitlab"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"NewVariables":    reflect.ValueOf(gitlab.NewVariables),
		"Parse":           reflect.ValueOf(gitlab.Parse),
		"WHEN_ALWAYS":     reflect.ValueOf(constant.MakeFromLiteral("\"always\"", token.STRING, 0)),
		"WHEN_DELAYED":    reflect.ValueOf(constant.MakeFromLiteral("\"delayed\"", token.STRING, 0)),
		"WHEN_MANUAL":     reflect.ValueOf(constant.MakeFromLiteral("\"manual\"", token.STRING, 0)),
		"WHEN_NEVER":      reflect.ValueOf(constant.MakeFromLiteral("\"never\"", token.STRING, 0)),
		"WHEN_ON_FAILURE": reflect.ValueOf(constant.MakeFromLiteral("\"on_failure\"", token.STRING, 0)),
		"WHEN_ON_SUCCESS": reflect.ValueOf(constant.MakeFromLiteral("\"on_success\"", token.STRING, 0)),

		// type definitions
		"AllowFailure": reflect.ValueOf((*gitlab.AllowFailure)(nil)),
		"Artifacts":    reflect.ValueOf((*gitlab.Artifacts)(nil)),
		"AutoCancel":   reflect.ValueOf((*gitlab.AutoCancel)(nil)),
		"Cache":        reflect.ValueOf((*gitlab.Cache)(nil)),
		"Image":        reflect.ValueOf((*gitlab.Image)(nil)),
		"Include":      reflect.ValueOf((*gitlab.Include)(nil)),
		"Job":          reflect.ValueOf((*gitlab.Job)(nil)),
		"Need":         reflect.ValueOf((*gitlab.Need)(nil)),
		"Needs":        reflect.ValueOf((*gitlab.Needs)(nil)),
		"Pipeline":     reflect.ValueOf((*gitlab.Pipeline)(nil)),
		"Rule":         reflect.ValueOf((*gitlab.Rule)(nil)),
		"Variables":    reflect.ValueOf((*gitlab.Variables)(nil)),
		"WorkFlow":     reflect.ValueOf((*gitlab.WorkFlow)(nil)),
		"WorkflowRule": reflect.ValueOf((*gitlab.WorkflowRule)(nil)),
	}
}
