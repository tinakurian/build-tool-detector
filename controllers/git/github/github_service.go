/*

Package github implements a way to extract
and construct a request to github in order
to retrieve a pom file. If the pom file is
not present, we assume the project is not
build using maven.

*/
package github

import (
	"errors"
	"github.com/tinakurian/build-tool-detector/app"
	"github.com/tinakurian/build-tool-detector/controllers/buildtype"
	errs "github.com/tinakurian/build-tool-detector/controllers/error"
	"log"
	"net/http"

	"github.com/google/go-github/github"
)

var (
	// ErrInternalServerError to return if unable to get contents
	ErrInternalServerError = errors.New("Unable to retrieve contents")
)

// GetGithubService something
func GetGithubService(ctx *app.ShowBuildToolDetectorContext, githubURL []string) (*errs.HTTPTypeError, *app.GoaBuildToolDetector) {
	// GetAttributes returns a BadRequest error and
	// will print the error to the user
	httpTypeError, attributes := getAttributes(githubURL, ctx.Branch)
	if httpTypeError != nil {
		log.Printf("Error: %v", httpTypeError)
		return httpTypeError, nil
	}

	// getGithubRepositoryPom returns an
	// InternalServerError and will print
	// the buildTool as unknown
	buildTool := buildtype.Unknown()
	httpTypeError = isMaven(ctx, attributes)
	if httpTypeError != nil {
		log.Printf("Error: %v", httpTypeError)
		return httpTypeError, buildTool
	}

	// Reset the buildToolType to maven since
	// the pom.xml was retrievable.
	buildTool.BuildToolType = buildtype.MAVEN

	return nil, buildTool
}

func isMaven(ctx *app.ShowBuildToolDetectorContext, attributes Attributes) *errs.HTTPTypeError {
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
