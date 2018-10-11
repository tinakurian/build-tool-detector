//go:generate goagen bootstrap -d github.com/tinakurian/build-tool-detector/design

package main

import (
	"errors"
	"flag"
	"os"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/tinakurian/build-tool-detector/app"
	controllers "github.com/tinakurian/build-tool-detector/controllers"
	"github.com/tinakurian/build-tool-detector/domain/git/github"
	"github.com/tinakurian/build-tool-detector/log"
	logorus "github.com/tinakurian/build-tool-detector/log"
)

var (
	// errFatalFailedSettingSentryDSN failed to set sentry env var
	errFatalFailedSettingSentryDSN = errors.New("failed to set environment variable for SENTRY_DSN: ")
)

const (
	startup           = "startup"
	errorz            = "err"
	buildToolDetector = "build-tool-detector"
	port              = "PORT"
	defaultPort       = "8080"
)

var (
	portNumber     = flag.String(port, defaultPort, "port")
	ghClientID     = flag.String(github.ClientID, "", "Github Client ID")
	ghClientSecret = flag.String(github.ClientSecret, "", "Github Client Secret")
	sentryDSN      = flag.String(log.SentryDSN, "", "Sentry DSN")
)

func main() {

	flag.Parse()
	if *ghClientID == "" || *ghClientSecret == "" {
		logorus.Logger().
			WithField(github.ClientID, ghClientID).
			WithField(github.ClientSecret, ghClientSecret).
			Fatalf(github.ErrFatalMissingGHAttributes.Error())
	}

	// Export Sentry DSN for logging
	err := os.Setenv(log.BuildToolDetectorSentryDSN, *sentryDSN)
	if err != nil {
		logorus.Logger().
			WithField(log.SentryDSN, sentryDSN).
			Fatalf(errFatalFailedSettingSentryDSN.Error()+"%v", sentryDSN)
	}

	// Create service
	service := goa.New(buildToolDetector)

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
	if err := service.ListenAndServe(":" + *portNumber); err != nil {
		service.LogError(startup, errorz, err)
	}
}
