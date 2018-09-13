package git

import (
	"log"
	"net/url"
	"strings"
)

// constants to define the different
// supported git services
const (
	// GITHUB service
	GITHUB = "github"

	// UNKNOWN service
	UNKNOWN = "unknown"
)

const (
	sDOTCOM = ".com"
	sSLASH  = "/"
)

// GetType something
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
