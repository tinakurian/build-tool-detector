package design

import (
	d "github.com/goadesign/goa/design"
	a "github.com/goadesign/goa/design/apidsl"
)

var _ = a.API("build-tool-detector", func() {
	a.Title("Detects the build tool for a specific repository and branch")
	a.Description("A simple goa service")
	a.Scheme("http")
	a.Host("localhost:8080")

})

// Endpoint to detect the build tool type based off the
// the repository url and branch
var _ = a.Resource("build-tool-detector", func() {
	a.BasePath("/build-tool")
	a.DefaultMedia(BuildToolDetectorMedia)

	a.Action("show", func() {
		a.Description("Get build tool")
		a.Routing(
			a.GET("/:url/:branch"),
		)
		a.Params(func() {
			a.Param("url", d.String, "github url")
			a.Param("branch", d.String, "branch name")
		})
		a.Response(d.OK)
		a.Response(d.NotFound)
	})
})

// BuildToolDetectorMedia defines the media type used to render the build tool
var BuildToolDetectorMedia = a.MediaType("application/vnd.goa.build.tool.detector+json", func() {
	a.Description("Detect the build tool for the specified repository and branch")
	a.Attributes(func() {
		a.Attribute("href", d.String, "API href for making requests on the build tool detection")
		a.Attribute("tool", d.String, "Name of build tool")
		a.Required("tool", "href")
	})
	a.View("default", func() {
		a.Attribute("tool")
		a.Attribute("href")
	})
})

// Endpoint to serve the swagger specification to use
// with a swagger editor: https://editor.swagger.io/
var _ = a.Resource("swagger", func() {
	a.Origin("*", func() {
		// Allow all origins to retrieve the Swagger JSON (CORS)
		a.Methods("GET")
	})
	a.Files("/swagger.json", "swagger/swagger.json")
})
