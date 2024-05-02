package gitlab

import (
	"fmt"
	"reflect"
)

type WorkFlow struct {
	Name string
	//AutoCancel AutoCancel `gitlabci:"auto_cancel"`
	//Rules      []WorkflowRule
}

func (workflow *WorkFlow) Parse(template any) error {
	tmplType := reflect.TypeOf(template)

	if tmplType.Kind() == reflect.Map {
		value := reflect.ValueOf(workflow).Elem()
		err := parseMap(&value, template.(map[string]any))
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("cannot parse workflow")
}

type AutoCancel struct {
	Commit     string `gitlabci:"on_new_commit"`
	JobFailure string `gitlabci:"on_job_failure"`
}

func (ac *AutoCancel) Parse(template any) error {
	tmplType := reflect.TypeOf(template)

	if tmplType.Kind() == reflect.Map {
		value := reflect.ValueOf(ac).Elem()
		err := parseMap(&value, template.(map[string]any))
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("cannot parse workflow")
}

type WorkflowRule struct {
	Rule
	AutoCancel AutoCancel `gitlabci:"auto_cancel"`
}

func (rule *WorkflowRule) Parse(template any) error {
	tmplType := reflect.TypeOf(template)

	if tmplType.Kind() == reflect.Map {
		value := reflect.ValueOf(rule).Elem()
		err := parseMap(&value, template.(map[string]any))
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("cannot parse workflow")
}
