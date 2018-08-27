package main

import (
	"build-tool-detector/app"
	"github.com/goadesign/goa"
	"net/http"
)

// BuildToolDetectorController implements the build-tool-detector resource
type BuildToolDetectorController struct {
	*goa.Controller
}

// NewBuildToolDetectorController creates a build-tool-detector controller
func NewBuildToolDetectorController(service *goa.Service) *BuildToolDetectorController {
	return &BuildToolDetectorController{Controller: service.NewController("BuildToolDetectorController")}
}

// Show runs the show action.
func (c *BuildToolDetectorController) Show(ctx *app.ShowBuildToolDetectorContext) error {
	generatedURL := ctx.URL + "/blob/" + ctx.Branch + "/pom.xml"

	response, err := http.Get(generatedURL)
	if err != nil || response.StatusCode != 200 {
		return ctx.NotFound()
	}

	buildTool := app.GoaBuildToolDetector{
		Tool: "Maven",
	}
	return ctx.OK(&buildTool)
}
