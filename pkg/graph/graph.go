package graph

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
)

type JobGraph struct {
	jobs  map[string]gitlab.Job
	edges map[string][]string
}

func NewGraph(pipeline gitlab.Pipeline, variables map[string]string) (*JobGraph, error) {
	var g JobGraph
	err := g.New(pipeline, variables)
	return &g, err
}

func (g *JobGraph) New(pipeline gitlab.Pipeline, variables map[string]string) error {
	enabledJobs := map[string]gitlab.Job{}

	jobs, err := pipeline.GetActiveJobs(variables)
	if err != nil {
		return err
	}

	for k, v := range jobs {
		if v.When == gitlab.WHEN_NEVER {
			continue
		}

		enabledJobs[k] = v
	}

	g.jobs = enabledJobs
	g.edges = make(map[string][]string)

	for job := range g.jobs {
		g.edges[job] = make([]string, 0)
	}

	for _, job := range g.jobs {
		g.AddJob(pipeline, job)
	}

	return nil
}

func (g *JobGraph) AddJob(pipeline gitlab.Pipeline, job gitlab.Job) {
	if job.Needs.NoNeeds {
		stageIndex := slices.Index(pipeline.Stages, job.Stage)

		for i := 0; i < stageIndex; i++ {
			// TODO: validate this behaviour, with rules evaluated pipelines
			for _, need := range pipeline.GetJobsByStage(pipeline.Stages[i]) {
				g.AddEdge(need.Name, job.Name)
			}
		}
	} else {
		for _, need := range job.Needs.Needs {
			g.AddEdge(need.Job, job.Name)
		}
	}

	for _, dep := range job.Dependencies {
		g.AddEdge(dep, job.Name)
	}
}

// Add a dependency for the end job with the start job
func (g *JobGraph) AddEdge(start string, end string) {
	if _, ok := g.edges[start]; ok {
		g.edges[start] = append(g.edges[start], end)
	}
}

func (g *JobGraph) GetJob(job string) (gitlab.Job, error) {
	j, ok := g.jobs[job]

	if !ok {
		return gitlab.Job{}, fmt.Errorf("job '%s' is not found in graph", job)
	}

	return j, nil
}

func (g *JobGraph) HasJob(job string) bool {
	_, ok := g.jobs[job]
	return ok
}

// Checks if a job has a (indirect) dependency on the given dependency
func (g *JobGraph) HasDependency(dependency string, job string) bool {
	if !g.HasJob(job) || !g.HasJob(dependency) {
		return false
	}

	for _, edge := range g.edges[dependency] {
		if edge == job {
			return true
		}

		if g.HasDependency(edge, job) {
			return true
		}
	}

	return false
}

// Get all direct and indirect dependencies for a given job
func (g *JobGraph) GetDependencies(job string) []string {
	jobs := make([]string, 0)

	for j := range g.jobs {
		if job == j {
			continue
		}

		if g.HasDependency(j, job) {
			jobs = append(jobs, j)
		}
	}

	return jobs
}

func (g JobGraph) Validate() error {
	for _, job := range g.jobs {
		for _, need := range job.Needs.Needs {
			if !g.HasJob(need.Job) {
				return fmt.Errorf("job '%s' needs job '%s' but it was not present in the pipeline", job.Name, need.Job)
			}
		}
	}

	return nil
}

func (g JobGraph) Map() map[string]any {
	return map[string]any{
		"dependencies": g.edges,
		"jobs":         g.jobs,
	}
}

func (g JobGraph) String() string {
	jsonbBytes, err := json.Marshal(g.Map())

	// We know all data is serializable as it originally came from a yaml file
	//  If magically we do still have an error, just panic
	if err != nil {
		panic(err)
	}

	return string(jsonbBytes)
}
