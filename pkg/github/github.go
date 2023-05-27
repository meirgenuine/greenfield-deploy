package github

import (
	"context"
	"os"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

var (
	c   *github.Client
	ctx context.Context
)

func init() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: os.Getenv("GITHUB_TOKEN"),
		},
	)

	ctx = context.Background()
	tc := oauth2.NewClient(ctx, ts)
	c = github.NewClient(tc)
}
