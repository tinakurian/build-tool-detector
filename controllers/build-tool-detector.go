package controllers

import (
	"build-tool-detector/app"
	"build-tool-detector/controllers/github"
	"build-tool-detector/controllers/parse"
	"github.com/goadesign/goa"
	"net/http"
)

// BuildToolDetectorController implements the build-tool-detector resource.
type BuildToolDetectorController struct {
	*goa.Controller
}

// NewBuildToolDetectorController creates a build-tool-detector controller.
func NewBuildToolDetectorController(service *goa.Service) *BuildToolDetectorController {
	return &BuildToolDetectorController{Controller: service.NewController("BuildToolDetectorController")}
}

// Show runs the show action.
func (c *BuildToolDetectorController) Show(ctx *app.ShowBuildToolDetectorContext) error {

	gitType, url := parse.GitType(ctx.URL)
	switch gitType {
	case parse.GITHUB:
		return handleGitHub(ctx, url)
	case parse.UNKNOWN:
		return ctx.InternalServerError()
	default:
		return ctx.InternalServerError()
	}
}

func handleGitHub(ctx *app.ShowBuildToolDetectorContext, url []string) error {
	statusCode, buildTool := github.GetGithubBuildTool(ctx, url)

	if statusCode == http.StatusInternalServerError {
		return ctx.InternalServerError()
	}

	return ctx.OK(&buildTool)
}
