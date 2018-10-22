//go:generate goagen bootstrap -d github.com/tinakurian/build-tool-detector/design

package main

import (
	"errors"
	"flag"
	"os"

	"github.com/fabric8-services/fabric8-common/goamiddleware"
	"github.com/fabric8-services/fabric8-common/token"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/goadesign/goa/middleware/security/jwt"
	"github.com/tinakurian/build-tool-detector/app"
	"github.com/tinakurian/build-tool-detector/config"
	"github.com/tinakurian/build-tool-detector/controllers"
	"github.com/tinakurian/build-tool-detector/domain/repository/github"
	"github.com/tinakurian/build-tool-detector/log"
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
	defaultPort       = "8099"
)

var (
	portNumber     = flag.String(port, defaultPort, "port")
	ghClientID     = flag.String(github.ClientID, "", "Github Client ID")
	ghClientSecret = flag.String(github.ClientSecret, "", "Github Client Secret")
	sentryDSN      = flag.String(log.SentryDSN, "", "Sentry DSN")
	authURL        = flag.String(config.AuthURL, config.AuthURLDefault, "URL to Auth")
)

func main() {

	flag.Parse()

	if *ghClientID == "" || *ghClientSecret == "" {
		log.Logger().
			WithField(github.ClientID, ghClientID).
			WithField(github.ClientSecret, ghClientSecret).
			Fatalf(github.ErrFatalMissingGHAttributes.Error())
	}

	// Export Sentry DSN for logging
	err := os.Setenv(log.BuildToolDetectorSentryDSN, *sentryDSN)
	if err != nil {
		log.Logger().
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

	tokenManager, err := token.NewManager(&config.AuthConfig{AuthServiceURL: *authURL})
	if err != nil {
		log.Logger().Panic(nil, map[string]interface{}{
			"err": err,
		}, "failed to create token manager")
	}
	// Middleware that extracts and stores the token in the context
	jwtMiddlewareTokenContext := goamiddleware.TokenContext(tokenManager, app.NewJWTSecurity())
	service.Use(jwtMiddlewareTokenContext)

	service.Use(token.InjectTokenManager(tokenManager))
	app.UseJWTMiddleware(service, jwt.New(tokenManager.PublicKeys(), nil, app.NewJWTSecurity()))

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
