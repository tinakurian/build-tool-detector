/*

Package git_test is used to test the functionality
within the git package.

*/
package git_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tinakurian/build-tool-detector/controllers/git"
	"github.com/tinakurian/build-tool-detector/controllers/git/github"
)

var _ = Describe("GitServiceType", func() {
	Context("GetGitServiceType", func() {
		It("Faulty Host - empty", func() {
			serviceType, err := git.GetGitServiceType("")
			Expect(serviceType).Should(BeNil(), "service type should be 'nil'")
			Expect(err.StatusCode).Should(BeEquivalentTo(http.StatusBadRequest), "service type should be '400'")
		})

		It("Faulty Host - non-existent", func() {
			serviceType, err := git.GetGitServiceType("test/test")
			Expect(serviceType).Should(BeNil(), "service type should be 'nil'")
			Expect(err.StatusCode).Should(BeEquivalentTo(http.StatusBadRequest), "service type should be '400'")
		})

		It("Faulty Host - not github.com", func() {
			serviceType, err := git.GetGitServiceType("http://test.com/test/test")
			Expect(serviceType).Should(BeNil(), "service type should be 'nil'")
			Expect(err.StatusCode).Should(BeEquivalentTo(http.StatusInternalServerError), "service type should be '500'")
		})

		It("Faulty url - no repository", func() {
			serviceType, err := git.GetGitServiceType("http://github.com/test")
			Expect(serviceType).Should(BeNil(), "service type should be 'nil'")
			Expect(err.StatusCode).Should(BeEquivalentTo(http.StatusBadRequest), "service type should be '400'")
		})

		It("Correct url - non-existent", func() {
			serviceType, err := git.GetGitServiceType("http://github.com/test/test")
			Expect(serviceType).ShouldNot(BeNil(), "service type should be not be'nil'")
			Expect(err).Should(BeNil(), "err should be 'nil'")
		})
	})

	Context("Constants", func() {
		It("Github", func() {
			Expect(git.Github).Should(BeEquivalentTo("github"), "git.GITHUB should be 'github'")
		})

		It("Unknown", func() {
			Expect(git.Unknown).Should(BeEquivalentTo("unknown"), "git.UNKNOWN should be 'unknown'")
		})
	})

	Context("Service", func() {
		It("Github", func() {
			service := git.Service{}.GetGitHubService()
			Expect(service).Should(BeEquivalentTo(&github.GitService{}), "service type should be a git service")
		})
	})
})
