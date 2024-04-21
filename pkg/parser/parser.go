package parser

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type parseable interface{
	Parse(any) error
}

type parsedMap map[interface{}]interface{}

type Pipeline struct {
	Stages []string
	Jobs []Job
}

func Parse(file string) (*Pipeline, error) {
	var pipeline Pipeline

	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("file read error: %v", err)
	}

	pipelineMap :=  make(parsedMap)

	err = yaml.Unmarshal(bytes, &pipelineMap)
	if err != nil {
		return nil, fmt.Errorf("yaml error: %v", err)
	}

	myjob := Job{}
	err = myjob.Parse("myjob", pipelineMap["myjob"].(parsedMap))
	if err != nil {
		return nil, fmt.Errorf("parse error: %v", err)
	}


	pipeline.Jobs = []Job{
		myjob,
	}

	return &pipeline, nil
}
