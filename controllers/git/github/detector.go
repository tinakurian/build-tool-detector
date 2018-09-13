/*

Package github implements a way to extract
and construct a request to github in order
to retrieve a pom file. If the pom file is
not present, we assume the project is not
build using maven.

*/
package github

import (
	"build-tool-detector/app"
	"build-tool-detector/controllers/git/buildtype"
	"net/http"
)

// DetectBuildTool retrieves the attributes
// from the url path parameter and attempts to
// retrieve the pom.xml file from the specified
// repository.
func DetectBuildTool(ctx *app.ShowBuildToolDetectorContext, githubURL []string) (int, app.GoaBuildToolDetector) {
	attributes := GetAttributes(githubURL, ctx.Branch)
	httpStatusCode := getGithubRepositoryPom(ctx, attributes)

	buildTool := buildtype.Unknown()

	if httpStatusCode == http.StatusInternalServerError {
		return http.StatusInternalServerError, buildTool
	}

	// Reset the buildToolType to maven since
	// the pom.xml was retrievable.
	buildTool.BuildToolType = buildtype.MAVEN

	return http.StatusOK, buildTool
}
