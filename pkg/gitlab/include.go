package gitlab

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/creasty/defaults"
	"github.com/watcherwhale/gogl-ci/pkg/gitlab/file"
)

type Include struct {
	File      []string
	Local     string
	Remote    string
	Project   string
	Component string
	Template  string
	Inputs    map[string]string `default:"{}"`
	Ref       string
	Rules     []Rule `default:"[]"`
	_parent   *Include
}

func (include *Include) hash() string {
	_, project, ref := include.isProject()

	return fmt.Sprintf("%v-%s-%s-%s-%s", include.File, include.Local, include.Remote, project, ref)
}

func (include *Include) Parse(template any) error {
	err := defaults.Set(include)
	if err != nil {
		return err
	}

	rTmpl := reflect.ValueOf(template)

	switch rTmpl.Kind() {
	case reflect.String:
		webMatch, err := regexp.Compile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()!@:%_\+.~#?&\/\/=]*)`)
		if err != nil {
			return err
		}

		if webMatch.Match([]byte(template.(string))) {
			include.Remote = template.(string)
		} else {
			include.Local = template.(string)
		}

	case reflect.Map:
		value := reflect.ValueOf(include).Elem()
		err := parseMap(&value, template.(map[any]any))
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("error parsing include: invalid type %s", rTmpl.Kind().String())
	}

	return nil
}

func (include *Include) isProject() (bool, string, string) {
	if include.Project != "" {
		return true, include.Project, include.Ref
	}

	if include._parent == nil {
		return false, "", ""
	}

	return include._parent.isProject()
}

func (include *Include) GetTemplate() ([]map[any]any, error) {
	if include.Local != "" {
		isProject, proj, ref := include.isProject()
		if !isProject {
			// This is a real local file
			templ, err := file.GetTemplateFile(include.Local[1:])
			if err != nil {
				return nil, err
			}
			return []map[any]any{templ}, nil
		} else {
			// This include was included from another project, loading it from this project instead
			templ, err := file.GetTemplateProject(include.Local[1:], proj, ref)
			if err != nil {
				return nil, err
			}
			return []map[any]any{templ}, nil
		}
	}

	if include.Remote != "" {
		templ, err := file.GetTemplateWeb(include.Remote)
		if err != nil {
			return nil, err
		}
		return []map[any]any{templ}, nil
	}

	if include.Project != "" {
		templArr := make([]map[any]any, len(include.File))
		for i, fileName := range include.File {
			templ, err := file.GetTemplateProject(fileName[1:], include.Project, include.Ref)
			if err != nil {
				return nil, err
			}

			templArr[i] = templ
		}
		return templArr, nil
	}

	return nil, fmt.Errorf("error getting included template: invalid include syntax")
}
