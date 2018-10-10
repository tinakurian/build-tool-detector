/*

Package git handles detecting build tool types
for git services such as github, bitbucket
and gitlab.

Currently the build-tool-detector only
supports github and can only recognize
maven.

*/
package git

import (
	"net/url"
	"strings"

	errs "github.com/tinakurian/build-tool-detector/controllers/error"
	"github.com/tinakurian/build-tool-detector/controllers/git/github"
)

// constants to define the different
// supported git services.
const (
	// GITHUB service.
	Github = "github"

	// UNKNOWN service.
	Unknown = "unknown"
)

const (
	dotcom = ".com"
	slash  = "/"
)

// Service struct.
type Service struct{}

// IService service interface.
type IService interface {
	GetGitHubService(string, string)
}

// GetGitHubService gets the github service.
func (s Service) GetGitHubService(ghClientID string, ghClientSecret string) *github.GitService {
	return &github.GitService{ClientID: ghClientID, ClientSecret: ghClientSecret}
}

// GetGitServiceType performs a simple url parse and split
// in order to retrieve the owner, repository
// and potentially the branch.
//
// Note: This method will likely need to be enhanced
// to handle different github url formats.
func GetGitServiceType(urlToParse string) (*string, *errs.HTTPTypeError) {
	gitServiceType := Github

	u, err := url.Parse(urlToParse)

	// Fail on error or empty host or empty scheme.
	if err != nil || u.Host == "" || u.Scheme == "" {
		return nil, errs.ErrBadRequest(github.ErrBadRequestInvalidPath)
	}

	// Currently only support Github.
	if u.Host != Github+dotcom {
		return nil, errs.ErrInternalServerError(github.ErrInternalServerErrorUnsupportedService)
	}

	urlSegments := strings.Split(u.Path, slash)
	if len(urlSegments) < 3 {
		return nil, errs.ErrBadRequest(github.ErrInternalServerErrorUnsupportedGithubURL)
	}

	return &gitServiceType, nil
}
