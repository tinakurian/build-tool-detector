/*

Package design is used to develop
the REST endpoints for the build tool.

*/
package design

import (
	d "github.com/goadesign/goa/design"
	a "github.com/goadesign/goa/design/apidsl"
)

// Endpoint to detect the build tool type based off the
// the repository url and branch
var _ = a.Resource("build-tool-detector", func() {
	a.BasePath("/")
	a.DefaultMedia(BuildToolDetectorMedia)
	a.Action("show", func() {
		a.Description("Detects the build tool for a given repository and branch.")
		a.Routing(
			a.GET("/:url"),
		)
		a.Params(func() {
			a.Param("url", d.String, "repository url")
			a.Param("branch", d.String, "repository branch")
		})
		a.Response(d.OK)
		a.Response(d.InternalServerError)
		a.Response(d.BadRequest)
		a.Response(d.NotFound)
	})
})

// BuildToolDetectorMedia defines the media type used to render the build tool
var BuildToolDetectorMedia = a.MediaType("application/vnd.goa.build.tool.detector+json", func() {
	a.Description("Detected build tool type.")
	a.Attributes(func() {
		a.Attribute("build-tool-type", d.String, "Name of build tool")
		a.Required("build-tool-type")
	})
	a.View("default", func() {
		a.Attribute("build-tool-type")
	})
})
