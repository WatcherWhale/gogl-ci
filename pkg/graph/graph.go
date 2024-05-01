package graph

import (
	"fmt"
	"slices"

	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
)

type JobGraph struct {
	jobs  map[string]gitlab.Job
	edges map[string][]string
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

func (g *JobGraph) AddJob(pipeline gitlab.Pipeline, job gitlab.Job) {
	if job.Needs.NoNeeds {
		stageIndex := slices.Index(pipeline.Stages, job.Stage)

		for i := 0; i < stageIndex; i++ {
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

// Checks if a job has a (indirect) dependency on the given dependency
func (g *JobGraph) HasDependency(dependency string, job string) bool {
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
