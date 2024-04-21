package api

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

func GetProjectFile(project string, file string, ref string) ([]byte, error) {
	if Client == nil {
		return nil, fmt.Errorf("not logged into gitlab")
	}

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
