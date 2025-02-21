package gitlab

import "fmt"

type Variables map[string]string

func NewVariables() Variables {
	return Variables{}
}

func (vars Variables) Copy() Variables {
	newVars := make(Variables)

	for k, v := range vars {
		newVars[k] = v
	}

	return newVars
}

func (vars Variables) Overwrite(overwrites Variables) Variables {
	for k, v := range overwrites {
		vars[k] = v
	}

	return vars
}

// Creates a copy of the current variables and with additional job variables set
func (vars Variables) ForJob(job Job) Variables {
	jobVars := Variables{
		"CI_JOB_IMAGE":      job.Image.Name,
		"CI_JOB_STATUS":     "running",
		"CI_JOB_TIMEOUT":    "1",
		"CI_JOB_ID":         "1",
		"CI_JOB_URL":        "https://gitlab.example.com/project/-/jobs/1",
		"CI_JOB_STARTED_AT": "2020-01-01T00:00:00Z",
		"CI_JOB_NAME":       job.Name,
		"CI_JOB_NAME_SLUG":  fmt.Sprintf("%s-%s", job.Stage, job.Name),
		"CI_JOB_STAGE":      job.Stage,
	}

	return jobVars.Overwrite(vars)
}

func (vars Variables) ForPipeline(pipeline Pipeline) Variables {
	plVars := Variables{
		"CI_PIPELINE_ID":  "1",
		"CI_PIPELINE_URL": "https://gitlab.example.com/project/-/pipelines/1",
	}

	return plVars.Overwrite(vars)
}

func (vars Variables) WithDefaultBranch(defaultBranch string) Variables {
	vars["CI_DEFAULT_BRANCH"] = defaultBranch
	return vars
}

func (vars Variables) WithBranch(branch string) Variables {
	vars["CI_COMMIT_BRANCH"] = branch
	return vars
}

func (vars Variables) WithMergeRequest() Variables {
	mrVars := Variables{
		"CI_MERGE_REQUEST_ID":                      "1",
		"CI_MERGE_REQUEST_IID":                     "1",
		"CI_MERGE_REQUEST_REF_PATH":                "refs/merge-requests/1/head",
		"CI_MERGE_REQUEST_PROJECT_ID":              "1",
		"CI_MERGE_REQUEST_PROJECT_PATH":            "project/gogl",
		"CI_MERGE_REQUEST_PROJECT_URL":             "https://gitlab.example.com/project/gogl",
		"CI_MERGE_REQUEST_TITLE":                   "title",
		"CI_MERGE_REQUEST_TARGET_BRANCH_NAME":      vars["CI_DEFAULT_BRANCH"],
		"CI_MERGE_REQUEST_TARGET_BRANCH_PROTECTED": "true",
		"CI_PIPELINE_SOURCE":                       "merge_request_event",
	}

	return mrVars.Overwrite(vars)
}
