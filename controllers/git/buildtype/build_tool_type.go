/*

Package buildtype implements a simple way to
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
package buildtype

import (
	"build-tool-detector/app"
)

const (
	// MAVEN build type detected Maven.
	MAVEN = "maven"

	// UNKNOWN build type detected Unknown.
	UNKNOWN = "unknown"
)

// Maven will create a buildToolDetector
// struct with the BuildToolType set
// to maven.
func Maven() *app.GoaBuildToolDetector {
	return &app.GoaBuildToolDetector{
		BuildToolType: MAVEN,
	}
}

// Unknown will create a buildToolDetector
// struct with the BuildToolType set
// to unknown.
func Unknown() *app.GoaBuildToolDetector {
	return &app.GoaBuildToolDetector{
		BuildToolType: UNKNOWN,
	}
}
