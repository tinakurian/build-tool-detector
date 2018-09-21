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
	errs "github.com/tinakurian/build-tool-detector/controllers/error"
	"github.com/tinakurian/build-tool-detector/controllers/git/github"
	"net/url"
	"strings"
)

// constants to define the different
// supported git services.
const (
	// GITHUB service.
	GITHUB = "github"

	// UNKNOWN service.
	UNKNOWN = "unknown"
)

const (
	sDOTCOM = ".com"
	sSLASH  = "/"
)

// Service something
type Service struct{}

// IService something
type IService interface {
	GetGitHubService(string)
}

// GetGitHubService something
func (g Service) GetGitHubService() *github.GoooService {
	return &github.GoooService{}
}

// GetGitServiceType performs a simple url parse and split
// in order to retrieve the owner, repository
// and potentially the branch.
//
// Note: This method will likely need to be enhanced
// to handle different github url formats.
func GetGitServiceType(urlToParse string) (*string, *errs.HTTPTypeError) {
	gitServiceType := GITHUB

	u, err := url.Parse(urlToParse)

	// Fail on error or empty host or empty scheme
	if err != nil || u.Host == "" || u.Scheme == "" {
		return nil, errs.ErrBadRequest(github.ErrBadRequestInvalidPath)
	}

	// Currently only support Github
	if u.Host != GITHUB+sDOTCOM {
		return nil, errs.ErrInternalServerError(github.ErrInternalServerErrorUnsupportedService)
	}

	segments := strings.Split(u.Path, sSLASH)
	if len(segments) < 3 {
		return nil, errs.ErrBadRequest(github.ErrInternalServerErrorUnsupportedGithubURL)
	}

	return &gitServiceType, nil
}
