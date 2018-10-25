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
	"github.com/tinakurian/build-tool-detector/config"
	"github.com/tinakurian/build-tool-detector/domain/types"
)

const (
	// ClientID github client id.
	ClientID = "GH_CLIENT_ID"

	// ClientSecret github client secret.
	ClientSecret = "GH_CLIENT_SECRET"
)
const (
	master = "master"
	tree   = "tree"
)

var (
	// ErrFailedContentRetrieval to return if unable to get contents.
	ErrFailedContentRetrieval = errors.New("unable to retrieve contents")

	// ErrUnsupportedGithubURL BadRequest github url is invalid.
	ErrUnsupportedGithubURL = errors.New("unsupported github url")

	// ErrInvalidPath github url is invalid.
	ErrInvalidPath = errors.New("url is invalid")

	// ErrResourceNotFound no resource found.
	ErrResourceNotFound = errors.New("resource not found")

	// ErrFatalMissingGHAttributes github client id and github client secret are unavailable.
	ErrFatalMissingGHAttributes = errors.New("github client id and github client secret are unavailable")
)

// RepositoryService contains
// values pertaining to a github
// repository.
type githubRepository struct {
	owner        string
	repository   string
	branch       string
	clientID     string
	clientSecret string
}

// result used to send results to
// the result channel.
type result struct {
	typeInfo *types.BuildType
	res      *github.Response
	err      error
}

// Create instantiate Github repository
func Create(segment []string, branch *string) (types.RepositoryService, error) {
	return newRepository(segment, branch)
}

// DetectBuildTool gets the contents for the service.
func (g githubRepository) DetectBuildTool(ctx context.Context) (*string, error) {
	// getGithubRepositoryPom returns an
	// InternalServerError and will print
	// the buildTool as unknown.
	buildTool := types.Unknown
	results := getContents(ctx, g)
	if results.err != nil {
		return &buildTool, results.err
	}

	return &results.typeInfo.BuildType, nil
}

// Owner returns the owner of a repository.
func (g githubRepository) Owner() string {
	return g.owner
}

// Repository returns the repository of a repository.
func (g githubRepository) Repository() string {
	return g.repository
}

// Branch returns the repository of a repository.
func (g githubRepository) Branch() string {
	return g.branch
}

// newRepository will use the path segments and
// query params to populate the Attributes
// struct. The attributes struct will be used
// to make a request to github to determine
// the build tool type.
func newRepository(segments []string, ctxBranch *string) (types.RepositoryService, error) {
	var configuration config.Configuration
	var repositoryService types.RepositoryService

	// Default branch that will be used if a branch
	// is not passed in though the optional 'branch'
	// query parameter and is not part of the url.
	branch := master

	if len(segments) <= 2 {
		return repositoryService, ErrInvalidPath
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

	repositoryService = githubRepository{
		owner:        segments[1],
		repository:   segments[2],
		branch:       branch,
		clientID:     configuration.Github.ClientID,
		clientSecret: configuration.Github.ClientSecret,
	}

	return repositoryService, nil
}

// getContents creates a client and
// initiates making requests to github.
func getContents(ctx context.Context, repository githubRepository) result {

	// Get the github client id and github client
	// secret if set to get better rate limits.
	t := github.UnauthenticatedRateLimitedTransport{
		ClientID:     repository.clientID,
		ClientSecret: repository.clientSecret,
	}

	// If the github client id or github client
	// secret are empty, we will log and fail.
	client := github.NewClient(t.Client())

	_, err := getBranchRequest(ctx, client, repository)
	if err != nil {
		return result{nil, nil, err}
	}

	// Parallel get requests.
	results := getContentsRequest(ctx, types.GetTypes(), client, repository)
	for _, result := range results {
		if result.res != nil {
			if result.res.StatusCode == http.StatusOK {
				return result
			}
		}
	}

	return result{nil, nil, ErrFailedContentRetrieval}
}

// getBranchRequest makes a request
// to ensure the repository and
// branch are valid.
func getBranchRequest(ctx context.Context, client *github.Client, repository githubRepository) (bool, error) {
	_, _, err := client.Repositories.GetBranch(ctx, repository.owner, repository.repository, repository.branch)
	if err != nil {
		return false, ErrResourceNotFound
	}

	return true, nil
}

// getContentsRequest makes parellel requests
// and sends the results through a results channel.
func getContentsRequest(ctx context.Context, buildTypes []types.BuildType, client *github.Client, repository githubRepository) []result {
	resultsChannel := make(chan result)
	defer func() {
		close(resultsChannel)
	}()

	for _, buildType := range buildTypes {
		go func(buildType types.BuildType) {
			_, _, resp, err := client.Repositories.GetContents(
				ctx, repository.owner,
				repository.repository,
				buildType.File,
				&github.RepositoryContentGetOptions{Ref: repository.branch})
			resultsChannel <- result{&buildType, resp, err}
		}(buildType)
	}

	var results []result
	for {
		result := <-resultsChannel
		results = append(results, result)

		// if we've reached the expected amount of urls then stop
		if len(results) == len(buildTypes) {
			break
		}
	}
	return results
}
