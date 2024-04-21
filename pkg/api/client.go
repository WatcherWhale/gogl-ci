package api

import (
	"github.com/xanzy/go-gitlab"
	"golang.org/x/time/rate"
)

var (
	Client *gitlab.Client
)

func Login(token string) error {
	// Create a client with an API_TOKEN and a rate limiter
	client, err := gitlab.NewClient(token, gitlab.WithCustomLimiter(rate.NewLimiter(10, 5)))

	if err != nil {
		return err
	}

	Client = client
	return nil
}

func GetProjectFile(project string, file string, ref string) ([]byte, error) {
	gitRef := ref
	if ref == "" {
		proj, _, err := Client.Projects.GetProject(project, &gitlab.GetProjectOptions{})
		if err != nil {
			return nil, err
		}

		gitRef = proj.DefaultBranch
	}

	bytes, _, err := Client.RepositoryFiles.GetRawFile(project, file, &gitlab.GetRawFileOptions{
		Ref: gitlab.Ptr(gitRef),
	})

	return bytes, err
}
