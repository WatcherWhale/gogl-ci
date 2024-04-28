package graph

import (
	"slices"

	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
)

type JobGraph struct {
	jobs  map[string]gitlab.Job
	edges map[string][]string
}

func (g *JobGraph) New(jobs map[string]gitlab.Job) {
	g.jobs = jobs
	g.edges = make(map[string][]string)

	for job := range jobs {
		g.edges[job] = make([]string, 0)
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

func (g *JobGraph) HasDependency(start string, dependency string) bool {
	for _, edge := range g.edges[start] {
		if edge == dependency {
			return true
		}

		if g.HasDependency(edge, dependency) {
			return true
		}
	}

	return false
}

func CalculateJobGraph(pipeline gitlab.Pipeline, variables map[string]string) *JobGraph {
	var jg JobGraph

	jg.New(pipeline.GetJobs())

	for _, job := range pipeline.GetJobs() {
		if job.IsEnabled(variables) {
			jg.AddJob(pipeline, job)
		}
	}

	return &jg
}
