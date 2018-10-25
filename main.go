//go:generate goagen bootstrap -d github.com/tinakurian/build-tool-detector/design

package main

import (
	"errors"
	"strconv"

	"github.com/fabric8-services/fabric8-common/goamiddleware"
	"github.com/fabric8-services/fabric8-common/token"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/goadesign/goa/middleware/security/jwt"
	"github.com/spf13/viper"
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
)

func main() {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	var configuration config.Configuration

	if err := viper.ReadInConfig(); err != nil {
		log.Logger().Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Logger().Fatalf("unable to decode into struct, %v", err)
	}

	if configuration.Github.ClientID == "" || configuration.Github.ClientSecret == "" {
		log.Logger().
			WithField(github.ClientID, configuration.Github.ClientID).
			WithField(github.ClientSecret, configuration.Github.ClientSecret).
			Fatalf(github.ErrFatalMissingGHAttributes.Error())
	}

	// Create service
	service := goa.New(buildToolDetector)

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	tokenManager, err := token.NewManager(&config.AuthConfiguration{URI: configuration.Auth.URI})
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
	c := controllers.NewBuildToolDetectorController(service)
	app.MountBuildToolDetectorController(service, c)

	cs := controllers.NewSwaggerController(service)
	app.MountSwaggerController(service, cs)

	// Start service
	if err := service.ListenAndServe(":" + strconv.Itoa(configuration.Server.Port)); err != nil {
		service.LogError(startup, errorz, err)
	}
}
