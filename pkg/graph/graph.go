package graph

import (
	"slices"

	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
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

func (g *JobGraph) HasJob(job string) bool {
	_, ok := g.jobMap[job]
	return ok
}

func (g *JobGraph) GetJob(job string) gitlab.Job {
	j := g.jobMap[job]
	return j
}

func (g *JobGraph) HasDependency(job string, dependency string) bool {
	j := g.jobMap[job]
	d := g.jobMap[dependency]

	if j.Needs != nil && len(j.Needs) == 0 {
		return false
	}

	sJ := slices.Index(g.stages, j.Stage)
	dJ := slices.Index(g.stages, d.Stage)

	// check if job is in not in an earlier stage than dependency
	if sJ < dJ {
		return false
	}

	// If no needs defined and in an earlier stage
	// if j.Needs == nil && sJ > dJ {
	// 	return true
	// }

	// If no needs defined and in the same stage
	// if j.Needs == nil && sJ == dJ {
	// 	return false
	// }

	return g.TraverseUntilFound(job, float64(sJ-dJ), func(job *gitlab.Job) bool {
		return job.Name == dependency
	}, func(job *gitlab.Job) float64 {
		return float64(sJ - slices.Index(g.stages, job.Stage))
	})
}

// Traverses the full graph until haltCondition is true or the remaining graph is fully explored
func (g *JobGraph) TraverseUntilFound(start string, startHeuristic float64, haltCondition func(job *gitlab.Job) bool, heuristic func(job *gitlab.Job) float64) bool {
	for _, dependent := range g.jobNeeds[start] {
		dJ := g.jobMap[dependent]
		dH := heuristic(&dJ)
		if dH > startHeuristic {
			continue
		}

		if haltCondition(&dJ) {
			return true
		}

		if g.TraverseUntilFound(dependent, dH, haltCondition, heuristic) {
			return true
		}
	}

	return false
}

func CalculateJobGraph(pipeline gitlab.Pipeline) *JobGraph {
	var g JobGraph
	g.New(pipeline.Stages)

	for _, job := range pipeline.GetJobs() {
		g.AddJob(job)
	}

	return &g
}
