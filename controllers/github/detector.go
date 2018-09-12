package github

import (
	"build-tool-detector/app"
	"net/http"
)

// GetGithubBuildTool (githubURL string)
func GetGithubBuildTool(ctx *app.ShowBuildToolDetectorContext, githubURL []string) (int, app.GoaBuildToolDetector) {

	attributes := GetAttributes(githubURL, ctx.Branch)
	httpStatusCode := getGithubRepositoryPom(ctx, attributes)

	buildTool := app.GoaBuildToolDetector{
		BuildToolType: "Unknown",
	}

	if httpStatusCode == http.StatusInternalServerError {
		return http.StatusInternalServerError, buildTool
	}

	buildTool.BuildToolType = "Maven"
	return http.StatusOK, buildTool
}
