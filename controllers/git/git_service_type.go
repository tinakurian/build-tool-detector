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
	GITHUB = "github"

	// UNKNOWN service.
	UNKNOWN = "unknown"
)

const (
	sDOTCOM = ".com"
	sSLASH  = "/"
)

// ServiceType contains the service
// and path segments
type ServiceType struct {
	Service  string
	Segments []string
}

// GetServiceType performs a simple url parse and split
// in order to retrieve the owner, repository
// and potentially the branch.
//
// Note: This method will likely need to be enhanced
// to handle different github url formats.
func GetServiceType(urlToParse string) (*errs.HTTPTypeError, ServiceType) {
	u, err := url.Parse(urlToParse)
	var service ServiceType

	// Fail on error or empty host or empty scheme
	if err != nil || u.Host == "" || u.Scheme == "" {
		return errs.ErrBadRequest(github.ErrBadRequest), service
	}

	segments := strings.Split(u.Path, sSLASH)
	if len(segments) < 2 {
		return errs.ErrBadRequest(github.ErrBadRequest), service
	}

	// Only support github service today
	service = ServiceType{
		Service:  UNKNOWN,
		Segments: segments,
	}

	if u.Host == GITHUB+sDOTCOM || len(segments) < 2 {
		service.Service = GITHUB
	}

	return nil, service
}
