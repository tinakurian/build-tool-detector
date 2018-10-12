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
	slash  = "/"
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
}

// Create instantiate Github repository
func Create() *Service {
	// TODO figure out from where this needs to be fetched - maybe we construct it differently somewhere else - this is just an idea
	return &Service{ClientID: "need-to-pass-it", ClientSecret: "need-to-pass-it-too"}
}

// GetContents gets the contents for the service.
func (g Service) GetContents(ctx context.Context, rawURL string, branchName *string) (*string, error) {
	// GetAttributes returns a BadRequest error and
	// will print the error to the user.
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, ErrInvalidPath
	}

	urlSegments := strings.Split(u.Path, slash)
	serviceAttribute, errs := getServiceAttributes(urlSegments, branchName)
	if err != nil {
		return nil, errs
	}

	// getGithubRepositoryPom returns an
	// InternalServerError and will print
	// the buildTool as unknown.
	buildTool := types.UnknownBuild
	_, errs = isMaven(ctx, g, serviceAttribute)
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

func isMaven(ctx context.Context, ghService Service, requestAttrs types.Repository) (bool, error) {

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
	_, _, err := client.Repositories.GetBranch(ctx, requestAttrs.Owner, requestAttrs.Repository, requestAttrs.Branch)
	if err != nil {
		return false, ErrResourceNotFound
	}

	// If the repository and branch exists, get the contents for the repository.
	_, _, resp, err := client.Repositories.GetContents(
		ctx, requestAttrs.Owner,
		requestAttrs.Repository,
		pom,
		&github.RepositoryContentGetOptions{Ref: requestAttrs.Branch})
	if err != nil && resp.StatusCode != http.StatusOK {
		return false, ErrFailedContentRetrieval
	}
	return true, nil
}
