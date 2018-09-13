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

// DetectBuildTool (githubURL string)
func DetectBuildTool(ctx *app.ShowBuildToolDetectorContext, githubURL []string) (int, app.GoaBuildToolDetector) {
	attributes := GetAttributes(githubURL, ctx.Branch)
	httpStatusCode := getGithubRepositoryPom(ctx, attributes)

	buildTool := buildtype.Unknown()

	if httpStatusCode == http.StatusInternalServerError {
		return http.StatusInternalServerError, buildTool
	}

	buildTool.BuildToolType = buildtype.MAVEN
	return http.StatusOK, buildTool
}
