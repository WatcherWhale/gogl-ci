package gitlab

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"slices"

	"github.com/watcherwhale/gogl-ci/pkg/rules/interpreter"

	"github.com/creasty/defaults"
	"github.com/rs/zerolog/log"
)

type Job struct {
	Name string `json:"name"`

	Image Image `default:"{}" json:"image"`

	Stage string `default:"test" json:"stage"`

	Script       []string `json:"script"`
	BeforeScript []string `default:"[]" gitlabci:"before_script" json:"before_script"`
	AfterScript  []string `default:"[]" gitlabci:"after_script" json:"after_script"`

	When         string   `default:"on_success" json:"on_success"`
	Rules        []Rule   `json:"rules"`
	Needs        Needs    `gitlabci:"needs" json:"needs"`
	Dependencies []string `json:"dependencies"`

	Variables map[string]string `json:"variables"`

	Interruptible bool `json:"interruptible"`

	Extends []string `json:"extends"`

	AllowFailure AllowFailure `gitlabci:"allow_failure" json:"allow_failure"`

	Artifacts Artifacts `json:"artifacts"`
	Cache     Cache     `json:"cache"`

	Coverage string `json:"coverage"`

	_keysWithValue []string `default:"[]" parser:"ignore"`
	_filled        bool     `default:"false" parser:"ignore"`
}

func (job *Job) Parse(name string, template map[string]any) error {
	err := defaults.Set(job)

	if err != nil {
		return fmt.Errorf("setting defaults error: %v", err)
	}

	keyMap := getFieldKeys(reflect.TypeOf(*job))

	job.Name = name

	structPtr := reflect.ValueOf(job).Elem()
	for yamlKey, value := range template {
		key, ok := keyMap[yamlKey]
		if !ok {
			log.Logger.Debug().Msgf("found unknown keyword %s", yamlKey)
			continue
		}

		field := structPtr.FieldByName(key)
		err := parseField(&field, key, value)
		if err != nil {
			return fmt.Errorf("error parsing key %s: %v", key, err)
		}

		job._keysWithValue = append(job._keysWithValue, key)
	}

	return nil
}

func (job *Job) Fill(pipeline *Pipeline) error {
	if job._filled {
		return nil
	}

	job.fill(pipeline.Default)

	for _, extendKey := range job.Extends {
		extendJob, ok := pipeline.Jobs[extendKey]

		if !ok {
			return fmt.Errorf("job %s not found, cannot parse job %s", extendKey, job.Name)
		}

		err := extendJob.Fill(pipeline)
		if err != nil {
			return err
		}
		job.fill(extendJob)
	}

	refRegex := regexp.MustCompile(`!reference \[(.*), rules\]`)

	rules := make([]Rule, 0)
	for _, rule := range job.Rules {
		if refRegex.MatchString(rule._reference) {
			ruleJob := pipeline.Jobs[string(refRegex.FindSubmatch([]byte(rule._reference))[1])]
			err := ruleJob.Fill(pipeline)
			if err != nil {
				return err
			}
			rules = append(rules, ruleJob.Rules...)
		} else {
			rules = append(rules, rule)
		}
	}

	job.Rules = rules

	job._filled = true

	return nil
}

func (job *Job) fill(template Job) {
	jobVal := reflect.ValueOf(job).Elem()
	templateVal := reflect.ValueOf(template)

	for _, fieldName := range template._keysWithValue {
		if slices.Contains(job._keysWithValue, fieldName) {
			continue
		}

		jobVal.FieldByName(fieldName).Set(templateVal.FieldByName(fieldName))
	}
}

func (job *Job) String() string {
	bytes, err := json.Marshal(job)
	if err != nil {
		return "{}"
	}

	return string(bytes)
}

// Get a copy of the current job, but overwritten with the first active rule
func (job Job) GetActiveJob(variables map[string]string) (Job, error) {
	activeJob := job

	for _, rule := range job.Rules {
		ok, err := interpreter.Evaluate(rule.If, variables)

		if err != nil {
			log.Warn().Err(err).Msgf("error occurred while evaluating rule '%s'", rule.If)
			return Job{}, err
		}

		if !ok {
			continue
		}

		activeJob.When = rule.When

		// Overwrites needs
		if !rule.Needs.NoNeeds {
			activeJob.Needs = rule.Needs
		}

		// Merge or overwrite variables
		for k, v := range rule.Variables {
			activeJob.Variables[k] = v
		}

		return activeJob, nil
	}

	activeJob.When = WHEN_NEVER

	return activeJob, nil
}
