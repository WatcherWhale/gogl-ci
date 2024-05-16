package file

import (
	"fmt"
	"io"
	"net/http"
)

func GetTemplateWeb(url string) (map[any]any, error) {
	//nolint:gosec
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return parseYaml(body)
}
