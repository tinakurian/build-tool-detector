/*

Package design is used to develop
the REST endpoints for the build tool.

*/
package design

import (
	a "github.com/goadesign/goa/design/apidsl"
)

// Endpoint to serve the swagger specification to use
// with a swagger editor: https://editor.swagger.io/
var _ = a.Resource("swagger", func() {
	a.Origin("*", func() {
		// Allow all origins to retrieve the Swagger JSON (CORS)
		a.Methods("GET")
	})
	a.Files("/swagger.json", "swagger/swagger.json")
})
