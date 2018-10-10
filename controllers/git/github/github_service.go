/*

Package github implements a way to extract
and construct a request to github in order
to retrieve a pom file. If the pom file is
not present, we assume the project is not
build using maven.

*/
package github

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-github/github"
	"github.com/tinakurian/build-tool-detector/controllers/buildtype"
	errs "github.com/tinakurian/build-tool-detector/controllers/error"
	logorus "github.com/tinakurian/build-tool-detector/log"
)

// serviceAttributes used for retrieving
// data using the go-github library.
type requestAttributes struct {
	Owner      string
	Repository string
	Branch     string
}

const (
	master         = "master"
	tree           = "tree"
	slash          = "/"
	pom            = "pom.xml"
	ghClientID     = "ghClientID"
	ghClientSecret = "ghClientSecret"
)

var (
	// ErrInternalServerErrorFailedContentRetrieval to return if unable to get contents.
	ErrInternalServerErrorFailedContentRetrieval = errors.New("unable to retrieve contents")

	// ErrInternalServerErrorUnsupportedGithubURL BadRequest github url is invalid.
	ErrInternalServerErrorUnsupportedGithubURL = errors.New("unsupported github url")

	// ErrBadRequestInvalidPath BadRequest github url is invalid.
	ErrBadRequestInvalidPath = errors.New("url is invalid")

	// ErrInternalServerErrorUnsupportedService git service unsupported.
	ErrInternalServerErrorUnsupportedService = errors.New("unsupported service")

	// ErrNotFoundResource no resource found.
	ErrNotFoundResource = errors.New("resource not found")

	// ErrFatalLimitedRateLimits github client id and github client secret are unavailable
	ErrFatalLimitedRateLimits = errors.New("github client id and github client secret are unavailable")
)

// IGitService git service interface.
type IGitService interface {
	GetContents(ctx context.Context) (*errs.HTTPTypeError, *string)
}

// GitService struct.
type GitService struct {
	ClientID     string
	ClientSecret string
}

// GetContents gets the contents for the service.
func (g GitService) GetContents(ctx context.Context, rawURL string, branchName *string) (*string, *errs.HTTPTypeError) {
	// GetAttributes returns a BadRequest error and
	// will print the error to the user.
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, errs.ErrBadRequest(ErrBadRequestInvalidPath)
	}

	urlSegments := strings.Split(u.Path, slash)
	serviceAttribute, httpTypeError := getServiceAttributes(urlSegments, branchName)
	if httpTypeError != nil {
		return nil, httpTypeError
	}

	// getGithubRepositoryPom returns an
	// InternalServerError and will print
	// the buildTool as unknown.
	buildTool := buildtype.UNKNOWN
	_, httpTypeError = isMaven(ctx, g, serviceAttribute)
	if httpTypeError != nil {
		return &buildTool, httpTypeError
	}

	// Reset the buildToolType to maven since
	// the pom.xml was retrievable.
	buildTool = buildtype.MAVEN

	return &buildTool, nil
}

// getServiceAttributes will use the path segments and
// query params to populate the Attributes
// struct. The attributes struct will be used
// to make a request to github to determine
// the build tool type.
func getServiceAttributes(segments []string, ctxBranch *string) (requestAttributes, *errs.HTTPTypeError) {

	var requestAttrs requestAttributes

	// Default branch that will be used if a branch
	// is not passed in though the optional 'branch'
	// query parameter and is not part of the url.
	branch := master

	if len(segments) <= 2 {
		return requestAttrs, errs.ErrBadRequest(ErrBadRequestInvalidPath)
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
		if segments[3] == tree {
			branch = segments[4]
		}
	}

	requestAttrs = requestAttributes{
		Owner:      segments[1],
		Repository: segments[2],
		Branch:     branch,
	}

	return requestAttrs, nil
}

func isMaven(ctx context.Context, ghService GitService, requestAttrs requestAttributes) (bool, *errs.HTTPTypeError) {

	// Get the github client id and github client
	// secret if set to get better rate limits.
	t := github.UnauthenticatedRateLimitedTransport{
		ClientID:     ghService.ClientID,
		ClientSecret: ghService.ClientSecret,
	}

	// If the github client id or github client
	// secret are empty, we will log and fail.
	client := github.NewClient(t.Client())
	if t.ClientID == "" || t.ClientSecret == "" {
		logorus.Logger().
			WithField(ghClientID, t.ClientID).
			WithField(ghClientSecret, t.ClientSecret).
			Fatalf(ErrFatalLimitedRateLimits.Error())
	}

	// Check that the repository + branch exists first.
	_, _, err := client.Repositories.GetBranch(ctx, requestAttrs.Owner, requestAttrs.Repository, requestAttrs.Branch)
	if err != nil {
		return false, errs.ErrNotFoundError(ErrNotFoundResource)
	}

	// If the repository and branch exists, get the contents for the repository.
	_, _, resp, err := client.Repositories.GetContents(
		ctx, requestAttrs.Owner,
		requestAttrs.Repository,
		pom,
		&github.RepositoryContentGetOptions{Ref: requestAttrs.Branch})
	if err != nil && resp.StatusCode != http.StatusOK {
		return false, errs.ErrInternalServerError(ErrInternalServerErrorFailedContentRetrieval)
	}
	return true, nil
}
