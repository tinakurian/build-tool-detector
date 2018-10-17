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

	"github.com/google/go-github/github"
	"github.com/tinakurian/build-tool-detector/domain/types"
	"github.com/tinakurian/build-tool-detector/log"
)

const (
	// ClientID github client id
	ClientID = "GH_CLIENT_ID"

	// ClientSecret github client secret
	ClientSecret = "GH_CLIENT_SECRET"
)
const (
	master = "master"
	tree   = "tree"
	pom    = "pom.xml"
)

var (
	// ErrFailedContentRetrieval to return if unable to get contents.
	ErrFailedContentRetrieval = errors.New("unable to retrieve contents")

	// ErrUnsupportedGithubURL BadRequest github url is invalid.
	ErrUnsupportedGithubURL = errors.New("unsupported github url")

	// ErrInvalidPath github url is invalid.
	ErrInvalidPath = errors.New("url is invalid")

	// ErrUnsupportedService git service unsupported.
	ErrUnsupportedService = errors.New("unsupported service")

	// ErrResourceNotFound no resource found.
	ErrResourceNotFound = errors.New("resource not found")

	// ErrFatalMissingGHAttributes github client id and github client secret are unavailable
	ErrFatalMissingGHAttributes = errors.New("github client id and github client secret are unavailable")
)

// Service struct.
type Service struct {
	ClientID     string
	ClientSecret string
	repository   types.Repository
}

// Create instantiate Github repository
func Create(segment []string, branch *string, ghClientID string, ghClientSecret string) (*Service, error) {
	repo, err := getServiceAttributes(segment, branch)
	if err != nil {
		return nil, err
	}
	return &Service{ClientID: ghClientID, ClientSecret: ghClientSecret, repository: repo}, nil
}

// DetectBuildTool gets the contents for the service.
func (g Service) DetectBuildTool(ctx context.Context) (*string, error) {
	// getGithubRepositoryPom returns an
	// InternalServerError and will print
	// the buildTool as unknown.
	buildTool := types.UnknownBuild
	_, errs := getContents(ctx, g)
	if errs != nil {
		return &buildTool, errs
	}

	// Reset the buildToolType to maven since
	// the pom.xml was retrievable.
	buildTool = types.MavenBuild

	return &buildTool, nil
}

// getServiceAttributes will use the path segments and
// query params to populate the Attributes
// struct. The attributes struct will be used
// to make a request to github to determine
// the build tool type.
func getServiceAttributes(segments []string, ctxBranch *string) (types.Repository, error) {

	var requestAttrs types.Repository

	// Default branch that will be used if a branch
	// is not passed in though the optional 'branch'
	// query parameter and is not part of the url.
	branch := master

	if len(segments) <= 2 {
		return requestAttrs, ErrInvalidPath
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

	requestAttrs = types.Repository{
		Owner:      segments[1],
		Repository: segments[2],
		Branch:     branch,
	}

	return requestAttrs, nil
}

func getContents(ctx context.Context, ghService Service) (bool, error) {

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
		log.Logger().
			WithField(ClientID, t.ClientID).
			WithField(ClientSecret, t.ClientSecret).
			Fatalf(ErrFatalMissingGHAttributes.Error())
	}

	// Check that the repository + branch exists first.
	_, _, err := client.Repositories.GetBranch(ctx, ghService.repository.Owner, ghService.repository.Repository, ghService.repository.Branch)
	if err != nil {
		return false, ErrResourceNotFound
	}

	// If the repository and branch exists, get the contents for the repository.
	_, _, resp, err := client.Repositories.GetContents(
		ctx, ghService.repository.Owner,
		ghService.repository.Repository,
		pom,
		&github.RepositoryContentGetOptions{Ref: ghService.repository.Branch})
	if err != nil && resp.StatusCode != http.StatusOK {
		return false, ErrFailedContentRetrieval
	}
	return true, nil
}
