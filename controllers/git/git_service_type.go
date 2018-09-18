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
	"github.com/tinakurian/build-tool-detector/app"
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

// serviceType contains the service
// and path segments
type serviceType struct {
	Service  string
	Segments []string
}

// GetGitService something
func GetGitService(ctx *app.ShowBuildToolDetectorContext) (*errs.HTTPTypeError, *app.GoaBuildToolDetector) {
	// Only support Github
	err, service := getServiceType(ctx.URL)
	if err != nil {
		return err, nil
	}

	// If the type is Github pass here
	err, buildTool := github.GetGithubService(ctx, service.Segments)
	if err != nil {
		return err, nil
	}

	return nil, buildTool
}

// GetServiceType performs a simple url parse and split
// in order to retrieve the owner, repository
// and potentially the branch.
//
// Note: This method will likely need to be enhanced
// to handle different github url formats.
func getServiceType(urlToParse string) (*errs.HTTPTypeError, serviceType) {
	u, err := url.Parse(urlToParse)
	var service serviceType

	// Fail on error or empty host or empty scheme
	if err != nil || u.Host == "" || u.Scheme == "" {
		return errs.ErrBadRequest(github.ErrBadRequest), service
	}

	segments := strings.Split(u.Path, sSLASH)
	if len(segments) < 2 {
		return errs.ErrBadRequest(github.ErrBadRequest), service
	}

	// Only support github service today
	service = serviceType{
		Service:  UNKNOWN,
		Segments: segments,
	}

	if u.Host == GITHUB+sDOTCOM || len(segments) < 2 {
		service.Service = GITHUB
		return errs.ErrBadRequest(github.ErrBadRequest), service
	}

	return nil, service
}
