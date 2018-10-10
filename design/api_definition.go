/*

Package design is used to develop
the REST endpoints for the build tool.

*/
package design

import (
	a "github.com/goadesign/goa/design/apidsl"
)

var _ = a.API("build-tool-detector", func() {
	a.Origin("*", func() {
		a.Methods("GET", "POST", "PUT", "PATCH", "DELETE")
		a.Headers("Accept", "Content-Type")
		a.Expose("Content-Type", "Origin")
		a.Credentials()
	})
	a.Title("Build Tool Detector")
	a.Description("Detects the build tool for a specific repository and branch. Currently, this tool only supports detecting the build tool maven for github repositories.")
	a.Scheme("http")
	a.Host("localhost:8080")
	a.BasePath("/build-tool-detector")
})
