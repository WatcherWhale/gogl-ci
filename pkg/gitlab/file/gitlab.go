package file

import (
	"github.com/rs/zerolog/log"
	"github.com/watcherwhale/gogl-ci/internal/cache"
	"github.com/watcherwhale/gogl-ci/pkg/api"
)

func GetTemplateProject(file, project, ref string) (map[any]any, error) {
	var bytes []byte

	cacheKey := cache.CreateGenericCacheKey(file, project, ref, api.GitlabUrl)
	if cache.LoadCache(cacheKey, &bytes) == nil {
		return parseYaml(bytes)
	}

	bytes, err := api.GetProjectFile(project, file, ref)
	if err != nil {
		return nil, err
	}

	err = cache.SaveCache(cacheKey, &bytes)
	if err != nil {
		log.Debug().Err(err).Msgf("failed to save cache with key '%s'", cacheKey)
	}

	return parseYaml(bytes)
}
