package controllers_test

import (
	"build-tool-detector/app/test"
	controllers "build-tool-detector/controllers"
	"github.com/davecgh/go-spew/spew"
	"github.com/goadesign/goa"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("BuildToolDetector", func() {
	Context("Internal Server Error", func() {
		It("Non-existent owner name -- 500 Internal Server Error", func() {
			service := goa.New("build-tool-detector")
			test.ShowBuildToolDetectorInternalServerError(GinkgoT(), nil, nil, controllers.NewBuildToolDetectorController(service), "https://github.com/fabric8-launcherz/launcher-backend", nil)
		})

		It("Non-existent repository name -- 500 Internal Server Error", func() {
			service := goa.New("build-tool-detector")
			test.ShowBuildToolDetectorInternalServerError(GinkgoT(), nil, nil, controllers.NewBuildToolDetectorController(service), "https://github.com/fabric8-launcher/launcher-backendz", nil)
		})

		It("Non-existent branch name -- 500 Internal Server Error", func() {
			service := goa.New("build-tool-detector")
			branch := "masterz"
			test.ShowBuildToolDetectorInternalServerError(GinkgoT(), nil, nil, controllers.NewBuildToolDetectorController(service), "https://github.com/fabric8-launcher/launcher-backend", &branch)
		})
	})

	Context("Okay", func() {
		It("200 Okay", func() {
			service := goa.New("build-tool-detector")
			branch := "master"
			response, buildTool := test.ShowBuildToolDetectorOK(GinkgoT(), nil, nil, controllers.NewBuildToolDetectorController(service), "https://github.com/fabric8-launcher/launcher-backend", &branch)
			spew.Dump(response, buildTool)
		})
	})
})
