/*

Package repository handles detecting build tool types
for git services such as github, bitbucket
and gitlab.

Currently the build-tool-detector only
supports github and can only recognize
maven.

*/
package repository

import (
	"net/url"
	"strings"

	"context"

	"github.com/tinakurian/build-tool-detector/domain/repository/github"
)

const (
	slash      = "/"
	githubHost = "github.com"
)

// Service service interface.
type Service interface {
	DetectBuildTool(ctx context.Context) (*string, error)
}

// CreateService performs a simple url parse and split
// in order to retrieve the owner, repository
// and potentially the branch.
//
// Note: This method will likely need to be enhanced
// to handle different github url formats.
func CreateService(urlToParse string, branch *string, ghClientID string, ghClientSecret string) (Service, error) {

	u, err := url.Parse(urlToParse)

	// Fail on error or empty host or empty scheme.
	if err != nil || u.Host == "" || u.Scheme == "" {
		return nil, github.ErrInvalidPath
	}

	// Currently only support Github.
	if u.Host != githubHost {
		return nil, github.ErrUnsupportedService
	}

	urlSegments := strings.Split(u.Path, slash)
	if len(urlSegments) < 3 {
		return nil, github.ErrUnsupportedGithubURL
	}

	return github.Create(urlSegments, branch, ghClientID, ghClientSecret)
}
