package design

import (
	d "github.com/goadesign/goa/design"
	a "github.com/goadesign/goa/design/apidsl"
)

// Endpoint to detect the build tool type based off the
// the repository url and branch
var _ = a.Resource("build-tool-detector", func() {
	a.BasePath("/build-tool")
	a.DefaultMedia(BuildToolDetectorMedia)
	a.Action("show", func() {
		a.Description("Get build tool")
		a.Routing(
			a.GET("/:owner/:repository/:branch"),
		)
		a.Params(func() {
			a.Param("owner", d.String, "github owner")
			a.Param("repository", d.String, "github repository name")
			a.Param("branch", d.String, "github branch")
		})
		a.Response(d.OK)
		a.Response(d.InternalServerError)
	})
})

// BuildToolDetectorMedia defines the media type used to render the build tool
var BuildToolDetectorMedia = a.MediaType("application/vnd.goa.build.tool.detector+json", func() {
	a.Description("Detect the build tool for the specified repository and branch")
	a.Attributes(func() {
		a.Attribute("build-tool-type", d.String, "Name of build tool")
		a.Required("build-tool-type")
	})
	a.View("default", func() {
		a.Attribute("build-tool-type")
	})
})
