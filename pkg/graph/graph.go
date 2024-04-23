package graph

import (
	"github.com/watcherwhale/gitlabci-test/pkg/gitlab"
)

type JobGraph struct {
	jobMap     map[string]gitlab.Job
	stageLinks map[string][]string
	jobNeeds   map[string][]string
	stages     []string
}

func (g *JobGraph) AddJob(job gitlab.Job) {
	g.jobMap[job.Name] = job
	g.stageLinks[job.Stage] = append(g.stageLinks[job.Stage], job.Name)

	for _, need := range job.Needs {
		if g.jobNeeds[need.Job] == nil {
			g.jobNeeds[need.Job] = make([]string, 1)
			g.jobNeeds[need.Job][0] = job.Name
		} else {
			g.jobNeeds[need.Job] = append(g.jobNeeds[need.Job], job.Name)
		}
	}
}

func (g *JobGraph) New(stages []string) {
	g.stages = stages
	g.stageLinks = map[string][]string{}
	g.jobMap = map[string]gitlab.Job{}
	g.jobNeeds = map[string][]string{}

	for _, stage := range stages {
		g.stageLinks[stage] = make([]string, 0)
	}
}
func CalculateJobGraph(pipeline gitlab.Pipeline) *JobGraph {
	var g JobGraph
	g.New(pipeline.Stages)

	for _, job := range pipeline.GetJobs() {
		g.AddJob(job)
	}

	return &g
}
