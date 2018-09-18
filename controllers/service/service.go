package service

import (
	"github.com/tinakurian/build-tool-detector/app"
	errs "github.com/tinakurian/build-tool-detector/controllers/error"
	"github.com/tinakurian/build-tool-detector/controllers/git"
)

// GetService Gets the appropriate service
// for the context provided
func GetService(ctx *app.ShowBuildToolDetectorContext) (*errs.HTTPTypeError, *app.GoaBuildToolDetector) {
	// Currently only support Git
	return git.GetGitService(ctx)
}
