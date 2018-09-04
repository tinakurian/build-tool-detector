package controllers

import (
	"build-tool-detector/app"
	"encoding/json"
	"fmt"
	"github.com/goadesign/goa"
	"github.com/google/go-github/github"
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

	client := github.NewClient(nil)
	_, _, resp, err := client.Repositories.GetContents(ctx, ctx.Owner, ctx.Repository, "pom.xml", &github.RepositoryContentGetOptions{Ref: ctx.Branch})

	// If there was an error or the status code returned
	// was not 200, return interal server error
	if err != nil || resp.StatusCode != http.StatusOK {

		ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		ctx.WriteHeader(http.StatusInternalServerError)
		buildTool := app.GoaBuildToolDetector{
			BuildToolType: "Unknown",
		}
		out, err := json.Marshal(buildTool)
		if err != nil {
			panic(err)
		}
		fmt.Fprint(ctx.ResponseWriter, string(out))
		return ctx.InternalServerError()
	}

	// If the response code was 200, return
	// the detected build tool type as Maven
	buildTool := app.GoaBuildToolDetector{
		BuildToolType: "Maven",
	}
	return ctx.OK(&buildTool)
}
