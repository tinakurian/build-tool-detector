//go:generate goagen bootstrap -d github.com/tinakurian/build-tool-detector/design

package main

import (
	"flag"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/tinakurian/build-tool-detector/app"
	controllers "github.com/tinakurian/build-tool-detector/controllers"
	logorus "github.com/tinakurian/build-tool-detector/log"
)

var (
	ghClientID     = flag.String("ghClientID", "", "Github Client ID")
	ghClientSecret = flag.String("ghClientSecret", "", "Github Client Secret")
)

func main() {

	flag.Parse()
	if *ghClientID == "" || *ghClientSecret == "" {
		logorus.Logger().
			WithField("ghClientID", ghClientID).
			WithField("ghClientSecret", ghClientSecret).
			Fatalf("Cannot run application without ghClientID and ghClientSecret")
	}

	// Create service
	service := goa.New("build-tool-detector")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount "build-tool-detector" controller
	c := controllers.NewBuildToolDetectorController(service, *ghClientID, *ghClientSecret)
	app.MountBuildToolDetectorController(service, c)

	cs := controllers.NewSwaggerController(service)
	app.MountSwaggerController(service, cs)

	// Start service
	if err := service.ListenAndServe(":8080"); err != nil {
		service.LogError("startup", "err", err)
	}
}
