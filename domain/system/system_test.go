/*

Package git_test is used to test the
functionality within the github package.

*/
package system_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/tinakurian/build-tool-detector/domain/git"
	. "github.com/tinakurian/build-tool-detector/domain/system"
)

var _ = Describe("Service", func() {
	Context("GetGitService", func() {
		It("Get git service", func() {
			gitService := System{}.GetGitService()
			gomega.Expect(gitService).Should(gomega.BeEquivalentTo(git.Service{}), "service type should be git service")
		})
	})
})
