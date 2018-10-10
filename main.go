//go:generate goagen bootstrap -d github.com/tinakurian/build-tool-detector/design

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/tinakurian/build-tool-detector/app"
	controllers "github.com/tinakurian/build-tool-detector/controllers"
	logorus "github.com/tinakurian/build-tool-detector/log"
)

var (
	ghClientID     = flag.String("ghClientID", "", "Github Client ID")
	ghClientSecret = flag.String("ghClientSecret", "", "Github Client Secret")
	sentryDSN      = flag.String("sentryDSN", "", "Sentry DSN")
)

func main() {

	flag.Parse()
	if *ghClientID == "" || *ghClientSecret == "" {
		logorus.Logger().
			WithField("ghClientID", ghClientID).
			WithField("ghClientSecret", ghClientSecret).
			Fatalf("Cannot run application without ghClientID and ghClientSecret")
	}

	fmt.Printf("SENTRYYYYYYYYYYYYYYYYY: %v \n\n\n\n", *sentryDSN)
	// Export Sentry DSN for logging
	err := os.Setenv("BUILD_TOOL_DETECTOR_SENTRY_DSN", *sentryDSN)
	if err != nil {
		logorus.Logger().
			WithField("sentryDSN", sentryDSN).
			Fatalf("Failed to set environment variable for sentryDSN: %v", sentryDSN)
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
