package graph

import (
	"slices"

	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
)

type JobGraph struct {
	jobs  map[string]gitlab.Job
	edges map[string][]string
}

func (g *JobGraph) New(pipeline gitlab.Pipeline, variables map[string]string) {
	enabledJobs := map[string]gitlab.Job{}

	for k, v := range pipeline.GetJobs() {
		if v.IsEnabled(variables) {
			enabledJobs[k] = v
		}
	}

	g.jobs = enabledJobs
	g.edges = make(map[string][]string)

	for job := range g.jobs {
		g.edges[job] = make([]string, 0)
	}

	for _, job := range g.jobs {
		g.AddJob(pipeline, job)
	}
}

func (g *JobGraph) GetJob(job string) (gitlab.Job, error) {
	return gitlab.Job{}, nil
}

func (g *JobGraph) HasJob(job string) bool {
	_, ok := g.jobs[job]
	return ok
}

func (g *JobGraph) AddJob(pipeline gitlab.Pipeline, job gitlab.Job) {
	for _, need := range job.Needs {
		g.AddEdge(need.Job, job.Name)
	}

	if job.Needs == nil {
		stageIndex := slices.Index(pipeline.Stages, job.Stage)

		for i := 0; i < stageIndex; i++ {
			for _, need := range pipeline.GetJobsByStage(pipeline.Stages[i]) {
				g.AddEdge(need.Name, job.Name)
			}
		}
	}
}

func (g *JobGraph) AddEdge(start string, end string) {
	g.edges[start] = append(g.edges[start], end)
}

func (g *JobGraph) HasDependency(dependency string, dependent string) bool {
	for _, edge := range g.edges[dependency] {
		if edge == dependent {
			return true
		}

		if g.HasDependency(edge, dependent) {
			return true
		}
	}

	return false
}

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
