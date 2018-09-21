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
	"github.com/google/go-github/github"
	"github.com/tinakurian/build-tool-detector/app"
	"github.com/tinakurian/build-tool-detector/controllers/buildtype"
	errs "github.com/tinakurian/build-tool-detector/controllers/error"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Attributes used for retrieving
// data using the go-github library.
type serviceAttributes struct {
	Owner      string
	Repository string
	Branch     string
}

const (
	sMASTER = "master"
	sTREE   = "tree"
)

var (
	// ErrInternalServerErrorFailedContentRetrieval to return if unable to get contents
	ErrInternalServerErrorFailedContentRetrieval = errors.New("Unable to retrieve contents")

	// ErrInternalServerErrorUnsupportedGithubURL BadRequest github url is invalid
	ErrInternalServerErrorUnsupportedGithubURL = errors.New("Unsupported github url")

	// ErrBadRequestInvalidPath BadRequest github url is invalid
	ErrBadRequestInvalidPath = errors.New("URL is invalid")

	// ErrInternalServerErrorUnsupportedService git service unsupported
	ErrInternalServerErrorUnsupportedService = errors.New("Unsupported service")

	// ErrNotFoundResource no resource found
	ErrNotFoundResource = errors.New("Resource not found")
)

// IGitService git service interface
type IGitService interface {
	GetContents(ctx *app.ShowBuildToolDetectorContext) (*errs.HTTPTypeError, *app.GoaBuildToolDetector)
}

// GitService struct
type GitService struct{}

// GetContents gets the contents for the service
func (g GitService) GetContents(ctx *app.ShowBuildToolDetectorContext) (*errs.HTTPTypeError, *app.GoaBuildToolDetector) {
	// GetAttributes returns a BadRequest error and
	// will print the error to the user
	u, err := url.Parse(ctx.URL)
	if err != nil {
		return errs.ErrBadRequest(ErrBadRequestInvalidPath), nil
	}

	segments := strings.Split(u.Path, "/")
	httpTypeError, attributes := getServiceAttributes(segments, ctx.Branch)
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

// getServiceAttributes will use the path segments and
// query params to populate the Attributes
// struct. The attributes struct will be used
// to make a request to github to determine
// the build tool type.
func getServiceAttributes(segments []string, ctxBranch *string) (*errs.HTTPTypeError, serviceAttributes) {

	var attributes serviceAttributes

	// Default branch that will be used if a branch
	// is not passed in though the optional 'branch'
	// query parameter and is not part of the url.
	branch := sMASTER

	if len(segments) <= 2 {
		return errs.ErrBadRequest(ErrBadRequestInvalidPath), attributes
	}

	// If the query parameter field 'branch' is not
	// empty then set the branch name to the query
	// parameter value.
	if ctxBranch != nil {
		branch = *ctxBranch
	} else if len(segments) > 4 {
		// If the user has not specified the branch
		// check whether it is passed in through
		// the URL.
		if segments[3] == sTREE {
			branch = segments[4]
		}
	}

	attributes = serviceAttributes{
		Owner:      segments[1],
		Repository: segments[2],
		Branch:     branch,
	}

	return nil, attributes
}

func isMaven(ctx *app.ShowBuildToolDetectorContext, attributes serviceAttributes) *errs.HTTPTypeError {

	t := github.UnauthenticatedRateLimitedTransport{
		ClientID:     "a0e1ce33654a8446356b",
		ClientSecret: "003e451564af39a5e29f768cbb9bcfd749577a31",
	}
	client := github.NewClient(t.Client())

	// Check that the repository + branch exists first
	_, _, err := client.Repositories.GetBranch(ctx, attributes.Owner, attributes.Repository, attributes.Branch)
	if err != nil {
		return errs.ErrNotFoundError(ErrNotFoundResource)
	}

	// If the repository and branch exists, get the contents for the repository
	_, _, resp, err := client.Repositories.GetContents(
		ctx, attributes.Owner,
		attributes.Repository,
		"pom.xml",
		&github.RepositoryContentGetOptions{Ref: attributes.Branch})
	if err != nil && resp.StatusCode != http.StatusOK {
		return errs.ErrInternalServerError(ErrInternalServerErrorFailedContentRetrieval)
	}

	return nil
}
