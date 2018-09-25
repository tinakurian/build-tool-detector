package system_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/tinakurian/build-tool-detector/controllers/git"
	. "github.com/tinakurian/build-tool-detector/controllers/system"
)

var _ = Describe("Service", func() {
	Context("GetGitService", func() {
		It("Get git service", func() {
			gitService := System{}.GetGitService()
			gomega.Expect(gitService).Should(gomega.BeEquivalentTo(git.Service{}), "service type should be git service")
		})
	})
})
