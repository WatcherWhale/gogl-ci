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
