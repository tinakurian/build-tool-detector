package config

// AuthConfiguration holds required values to use with authorization service (while setting JWT middleware for example).
type AuthConfiguration struct {
	URI string
}

// GetAuthServiceURL provides URL used to call Authorization service.
func (c *AuthConfiguration) GetAuthServiceURL() string {
	return c.URI
}

// GetAuthKeysPath provides a URL path to be called for retrieving the keys.
func (c *AuthConfiguration) GetAuthKeysPath() string {
	// Fixed with https://github.com/fabric8-services/fabric8-common/pull/25.
	return "/api/token/keys"
}

// GetDevModePrivateKey not used right now.
func (c *AuthConfiguration) GetDevModePrivateKey() []byte {
	// No need for now
	return nil
}
