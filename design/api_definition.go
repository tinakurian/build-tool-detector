/*

Package design is used to develop
the REST endpoints for the build tool.

*/
package design

import (
	a "github.com/goadesign/goa/design/apidsl"
)

var _ = a.API("build-tool-detector", func() {
	a.Title("Build Tool Detector")
	a.Description("Detects the build tool for a specific repository and branch. Currently, this tool only supports detecting the build tool maven for github repositories.")

	a.Origin("/[.*openshift.io|localhost]/", func() {
		a.Methods("GET")
		a.Headers("X-Request-Id", "Content-Type", "Authorization")
		a.MaxAge(600)
		a.Credentials()
	})

	a.Scheme("http")
	a.Host("localhost:8080")
	a.Version("1.0")
	a.BasePath("/build-tool-detector")

	a.License(func() {
		a.Name("Apache License Version 2.0")
		a.URL("http://www.apache.org/licenses/LICENSE-2.0")
	})

	a.JWTSecurity("jwt", func() {
		a.Description("JWT Token Auth")
		a.Header("Authorization")
	})
})
