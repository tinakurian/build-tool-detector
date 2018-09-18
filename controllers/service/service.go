package service

import (
	"github.com/tinakurian/build-tool-detector/app"
	errs "github.com/tinakurian/build-tool-detector/controllers/error"
	"github.com/tinakurian/build-tool-detector/controllers/git"
)

// GetService something
func GetService(ctx *app.ShowBuildToolDetectorContext) (*errs.HTTPTypeError, *app.GoaBuildToolDetector) {
	// Only support Git
	return git.GetGitService(ctx)
}
