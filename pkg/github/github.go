package github

import (
	"context"
	"os"
	"sync"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

var (
	c    *github.Client
	ctx  context.Context
	once sync.Once
)

func Init() {
	once.Do(func() {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{
				AccessToken: os.Getenv("GITHUB_TOKEN"),
			},
		)

		ctx = context.Background()
		tc := oauth2.NewClient(ctx, ts)
		c = github.NewClient(tc)
	})
}
