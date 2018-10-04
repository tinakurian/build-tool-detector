package log

import (
	"os"

	"github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
)

// Logger something
func Logger() *logrus.Entry {

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.WarnLevel)

	// TODO: have env variable to specify we are running tests which
	// will return a different logger.
	sentryDSN := os.Getenv("SENTRY_DSN")

	if sentryDSN != "" {
		hook, err := logrus_sentry.NewSentryHook(sentryDSN, []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		})

		if err != nil {
			panic(err)
		}

		// Add sentry hook
		logrus.AddHook(hook)
	} else {
		logrus.SetOutput(os.Stdout)
	}
	return logrus.WithField("applicationName", "build-tool-detector")
}
