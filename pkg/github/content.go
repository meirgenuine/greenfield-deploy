package github

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/google/go-github/v28/github"
)

func DownloadContent(m *github.RepositoryContent) (io.ReadCloser, error) {
	return c.Repositories.DownloadContents(ctx,
		os.Getenv("GITHUB_REPO"), "greenfield-deploy",
		m.GetPath(), &github.RepositoryContentGetOptions{})
}

func Content(project, environment string) (map[string]*github.RepositoryContent, error) {
	_, ff, _, err := c.Repositories.GetContents(ctx,
		os.Getenv("GITHUB_REPO"), "greenfield-deploy",
		fmt.Sprintf("deployments/%s", project),
		&github.RepositoryContentGetOptions{})
	if err != nil {
		return nil, err
	}

	prefix := fmt.Sprintf("k8s_%s_", environment)
	rr := map[string]*github.RepositoryContent{}
	for _, f := range ff {
		if strings.Index(f.GetName(), prefix) != 0 {
			continue
		}
		rr[f.GetName()] = f
	}
	return rr, nil
}
