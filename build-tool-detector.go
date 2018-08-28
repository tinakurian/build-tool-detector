package main

import (
	"build-tool-detector/app"
	"github.com/goadesign/goa"
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
	// BuildToolDetectorController_Show: start_implement

	// Put your logic here

	res := &app.GoaBuildToolDetector{}
	return ctx.OK(res)
	// BuildToolDetectorController_Show: end_implement
}
