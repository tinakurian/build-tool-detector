package github

import (
	"build-tool-detector/app"
	"github.com/google/go-github/github"
	"net/http"
)

// GetGithubRepositoryPom something
func getGithubRepositoryPom(ctx *app.ShowBuildToolDetectorContext, attributes Attributes) int {
	client := github.NewClient(nil)

	_, _, resp, err := client.Repositories.GetContents(
		ctx, attributes.Owner,
		attributes.Repository,
		"pom.xml",
		&github.RepositoryContentGetOptions{Ref: attributes.Branch})

	if err != nil || resp.StatusCode != http.StatusOK {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}
