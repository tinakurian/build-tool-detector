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
	logorus "github.com/tinakurian/build-tool-detector/log"
	"net/url"
	"strings"
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
	dotcom   = ".com"
	slash    = "/"
	segments = "url segments"
	host     = "host"
)

// Service struct.
type Service struct{}

// IService service interface.
type IService interface {
	GetGitHubService(string)
}

// GetGitHubService gets the github service.
func (s Service) GetGitHubService() *github.GitService {
	return &github.GitService{}
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
		logorus.Logger().
			WithError(github.ErrBadRequestInvalidPath).
			WithField(segments, u).
			Warningf(github.ErrBadRequestInvalidPath.Error())
		return nil, errs.ErrBadRequest(github.ErrBadRequestInvalidPath)
	}

	// Currently only support Github.
	if u.Host != Github+dotcom {
		logorus.Logger().
			WithError(github.ErrBadRequestInvalidPath).
			WithField(host, u.Host).
			Warningf(github.ErrInternalServerErrorUnsupportedService.Error())
		return nil, errs.ErrInternalServerError(github.ErrInternalServerErrorUnsupportedService)
	}

	urlSegments := strings.Split(u.Path, slash)
	if len(urlSegments) < 3 {
		logorus.Logger().
			WithError(github.ErrBadRequestInvalidPath).
			WithField(segments, urlSegments).
			Warningf(github.ErrInternalServerErrorUnsupportedGithubURL.Error())
		return nil, errs.ErrBadRequest(github.ErrInternalServerErrorUnsupportedGithubURL)
	}

	return &gitServiceType, nil
}
