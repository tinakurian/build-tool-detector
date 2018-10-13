/*

Package types implements a simple way to
create a GoaBuildToolDetector struct with
the recognized build type.

The GoaBuildToolDetector struct will result
in the following JSON response:

	[[ maven ]]

		{
		"build-tool-type": "maven"
		}


	[[ unknown ]]

		{
		"build-tool-type": "unknown"
		}

*/
package types

import (
	"github.com/tinakurian/build-tool-detector/app"
)

const (
	// MavenBuild build type detected Maven.
	MavenBuild = "maven"

	// UnknownBuild build type detected Unknown.
	UnknownBuild = "unknown"
)

// Maven will create a buildToolDetector
// struct with the BuildToolType set
// to maven.
func Maven() *app.GoaBuildToolDetector {
	return &app.GoaBuildToolDetector{
		BuildToolType: MavenBuild,
	}
}

// Unknown will create a buildToolDetector
// struct with the BuildToolType set
// to unknown.
func Unknown() *app.GoaBuildToolDetector {
	return &app.GoaBuildToolDetector{
		BuildToolType: UnknownBuild,
	}
}
