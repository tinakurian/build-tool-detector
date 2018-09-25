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
			Expect(err.StatusCode).Should(BeEquivalentTo(http.StatusBadRequest), "service type should be 'nil'")
		})

		It("Faulty Host - non-existent", func() {
			serviceType, err := git.GetGitServiceType("test/test")
			Expect(serviceType).Should(BeNil(), "service type should be 'nil'")
			Expect(err.StatusCode).Should(BeEquivalentTo(http.StatusBadRequest), "service type should be 'nil'")
		})

		It("Faulty Host - not github.com", func() {
			serviceType, err := git.GetGitServiceType("http://test.com/test/test")
			Expect(serviceType).Should(BeNil(), "service type should be 'nil'")
			Expect(err.StatusCode).Should(BeEquivalentTo(http.StatusInternalServerError), "service type should be 'nil'")
		})

		It("Faulty url - no repository", func() {
			serviceType, err := git.GetGitServiceType("http://github.com/test")
			Expect(serviceType).Should(BeNil(), "service type should be 'nil'")
			Expect(err.StatusCode).Should(BeEquivalentTo(http.StatusBadRequest), "service type should be 'nil'")
		})

		It("Correct url - non-existent", func() {
			serviceType, err := git.GetGitServiceType("http://github.com/test/test")
			Expect(serviceType).ShouldNot(BeNil(), "service type should be 'nil'")
			Expect(err).Should(BeNil(), "service type should be 'nil'")
		})
	})

	Context("Constants", func() {
		It("Github", func() {
			Expect(git.GITHUB).Should(BeEquivalentTo("github"), "service type should be 'nil'")
		})

		It("Unknown", func() {
			Expect(git.UNKNOWN).Should(BeEquivalentTo("unknown"), "service type should be 'nil'")
		})
	})

	Context("Service", func() {
		It("Github", func() {
			service := git.Service{}.GetGitHubService()
			Expect(service).Should(BeEquivalentTo(&github.GitService{}), "service type should be 'nil'")
		})
	})
})
