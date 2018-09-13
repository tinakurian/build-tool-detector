/*

Package github implements a way to extract
and construct a request to github in order
to retrieve a pom file. If the pom file is
not present, we assume the project is not
build using maven.

*/
package github

import (
	"build-tool-detector/app"
	"github.com/google/go-github/github"
	"net/http"
)

// GetGithubRepositoryPom requests the pom.xl
// file to determine whether the project is
// built using maven.
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
