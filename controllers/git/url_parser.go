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
	"log"
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

// GetType performs a simple url parse and split
// in order to retrieve the owner, repository
// and potentially the branch.
//
// Note: This method will likely need to be enhanced
// to handle different github url formats.
func GetType(urlToParse string) (string, []string) {
	u, err := url.Parse(urlToParse)
	if err != nil {
		log.Fatal(err)
	}

	segments := strings.Split(u.Path, sSLASH)
	if u.Host == GITHUB+sDOTCOM {
		return GITHUB, segments
	}

	return UNKNOWN, segments
}
