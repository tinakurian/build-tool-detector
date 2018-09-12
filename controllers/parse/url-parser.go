package parse

import (
	"log"
	"net/url"
	"strings"
)

// constants to define the different
// supported git types
const (
	GITHUB  = "github"
	UNKNOWN = ""
)

// URL (urlToParse)
func urlSegments(urlToParse *url.URL) []string {

	paths := strings.Split(urlToParse.Path, "/")

	return paths
}

// GitType something
func GitType(urlToParse string) (string, []string) {
	u, err := url.Parse(urlToParse)
	if err != nil {
		log.Fatal(err)
	}

	segments := urlSegments(u)
	if u.Host == "github.com" {
		return GITHUB, segments
	}

	return UNKNOWN, segments
}
