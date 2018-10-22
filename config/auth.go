package config

const (
	// AuthURL holds the key name based on which flag or environment variable is looked up
	AuthURL = "F8_BUILD_DETECTOR_AUTH_URL"
	// AuthURLDefault is a default value in case it's not provided
	AuthURLDefault = "https://auth.openshift.io"
)

// AuthConfig holds required values to use with authorization service (while setting JWT middleware for example)
type AuthConfig struct {
	AuthServiceURL string
}

// GetAuthServiceURL provides URL used to call Authorization service
func (c *AuthConfig) GetAuthServiceURL() string {
	return c.AuthServiceURL
}

// GetAuthKeysPath provides a URL path to be called for retrieving the keys
func (c *AuthConfig) GetAuthKeysPath() string {
	// Fixed with https://github.com/fabric8-services/fabric8-common/pull/25
	return "/api/token/keys"
}

// GetDevModePrivateKey not used right now
func (c *AuthConfig) GetDevModePrivateKey() []byte {
	// No need for now
	return nil
}
