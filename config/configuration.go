package config

// Configuration for build tool detector
type Configuration struct {
	Auth   AuthConfiguration
	Github GithubConfiguration
	Sentry SentryConfiguration
	Server ServerConfiguration
}
