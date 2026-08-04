package file

import (
	"fmt"
	"regexp"

	"gopkg.in/yaml.v3"
)

func parseYaml(b []byte) (map[any]any, error) {
	pipelineMap := make(map[any]any)

	replacer, err := regexp.Compile(`!reference\s+\[([^\]]+)\]`)
	if err != nil {
		return nil, fmt.Errorf("yaml error: %v", err)
	}

	// The !reference tag is not supported by the builtin yaml parser, so we need to replace it with a string that can be parsed later.
	b = replacer.ReplaceAll(b, []byte("\"!reference [$1]\""))

	// The builtin yaml parser does not handle \/ like the gitlab pipeline parser does, so we need to fix that before unmarshalling.
	b = fixSlashEscape(b)

	err = yaml.Unmarshal(b, &pipelineMap)
	if err != nil {
		return nil, fmt.Errorf("yaml error: %v", err)
	}

	return pipelineMap, nil
}

func fixSlashEscape(b []byte) []byte {
	out := make([]byte, 0, len(b))
	for i := 0; i < len(b); i++ {
		if b[i] == '\\' && i+1 < len(b) {
			switch b[i+1] {
			case '/':
				out = append(out, '\\', '\\', '/')
				i++
			case '\\':
				out = append(out, '\\', '\\')
				i++
			default:
				out = append(out, b[i])
			}
		} else {
			out = append(out, b[i])
		}
	}
	return out
}
