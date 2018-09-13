package buildtype

import (
	"build-tool-detector/app"
)

const (
	// MAVEN build type detected Maven
	MAVEN = "maven"

	// UNKNOWN build type detected Unknown
	UNKNOWN = "unknown"
)

// Maven something
func Maven() app.GoaBuildToolDetector {
	return app.GoaBuildToolDetector{
		BuildToolType: MAVEN,
	}
}

// Unknown something
func Unknown() app.GoaBuildToolDetector {
	return app.GoaBuildToolDetector{
		BuildToolType: UNKNOWN,
	}
}
