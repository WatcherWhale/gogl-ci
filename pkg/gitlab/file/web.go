package file

import (
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/watcherwhale/gogl-ci/internal/cache"
)

func GetTemplateWeb(url string) (map[any]any, error) {
	var body []byte

	cacheKey := cache.CreateGenericCacheKey(url)
	if cache.LoadCache(cacheKey, &body) == nil {
		return parseYaml(body)
	}

	//nolint:gosec
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = cache.SaveCache(cacheKey, &body)
	if err != nil {
		log.Debug().Err(err).Msgf("failed to save cache with key '%s'", cacheKey)
	}

	return parseYaml(body)
}
