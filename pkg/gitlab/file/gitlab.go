package file

import "github.com/watcherwhale/gogl-ci/pkg/api"

func GetTemplateProject(file, project, ref string) (map[any]any, error) {
	bytes, err := api.GetProjectFile(project, file, ref)
	if err != nil {
		return nil, err
	}

	return parseYaml(bytes)
}
