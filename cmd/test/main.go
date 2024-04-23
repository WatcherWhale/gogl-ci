package main

import (
	"fmt"

	"github.com/watcherwhale/gitlabci-test/pkg/gitlab/rules/lexer"
)

func main() {
	tokens, err := lexer.Tokenize("$CI_COMMIT_BRANCH == null && $CI_COMMIT_BRANCH == \"test\"")

	if err != nil {
		fmt.Printf("%v", err)
	}

	for _, token := range tokens {
		println(token.Value)
	}
}
