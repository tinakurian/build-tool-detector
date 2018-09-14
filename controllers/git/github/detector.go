/*

Package github implements a way to extract
and construct a request to github in order
to retrieve a pom file. If the pom file is
not present, we assume the project is not
build using maven.

*/
package github

import (
  errs "github.com/tinakurian/build-tool-detector/controllers/error"
	"github.com/tinakurian/build-tool-detector/app"
	"github.com/tinakurian/build-tool-detector/controllers/git/buildtype"
	"net/http"
	"log"
)

// DetectBuildTool retrieves the attributes
// from the url path parameter and attempts to
// retrieve the pom.xml file from the specified
// repository.
func DetectBuildTool(ctx *app.ShowBuildToolDetectorContext, githubURL []string) (*errs.HTTPTypeError, *app.GoaBuildToolDetector) {

	// GetAttributes returns a BadRequest error and
	// will print the error to the user
	httpTypeError, attributes := GetAttributes(githubURL, ctx.Branch)
	if httpTypeError != nil {
		log.Printf("Error: %v", httpTypeError)
		return httpTypeError, nil
	}

	// getGithubRepositoryPom returns an
	// InternalServerError and will print
	// the buildTool as unknown
	buildTool := buildtype.Unknown()
	httpTypeError = getGithubRepositoryPom(ctx, attributes)
	if httpTypeError != nil {
		log.Printf("Error: %v", httpTypeError)
		return httpTypeError, buildTool
	}

	// Reset the buildToolType to maven since
	// the pom.xml was retrievable.
	buildTool.BuildToolType = buildtype.MAVEN

	return nil, buildTool
}
