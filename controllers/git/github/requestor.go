/*

Package github implements a way to extract
and construct a request to github in order
to retrieve a pom file. If the pom file is
not present, we assume the project is not
build using maven.

*/
package github

import (
	"github.com/tinakurian/build-tool-detector/app"
	errs "github.com/tinakurian/build-tool-detector/controllers/error"
	"errors"
	"github.com/google/go-github/github"
	"net/http"
)

var (
	// ErrInternalServerError to return if unable to get contents
	ErrInternalServerError = errors.New("Unable to retrieve contents")
)

// GetGithubRepositoryPom requests the pom.xl
// file to determine whether the project is
// built using maven.
func getGithubRepositoryPom(ctx *app.ShowBuildToolDetectorContext, attributes Attributes) *errs.HTTPTypeError {
	client := github.NewClient(nil)
	_, _, resp, err := client.Repositories.GetContents(
		ctx, attributes.Owner,
		attributes.Repository,
		"pom.xml",
		&github.RepositoryContentGetOptions{Ref: attributes.Branch})

	if err != nil || resp.StatusCode != http.StatusOK {
		return errs.ErrInternalServerError(ErrInternalServerError)
	}

	return nil
}
