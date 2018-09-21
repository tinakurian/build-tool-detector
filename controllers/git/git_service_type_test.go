package git_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/tinakurian/build-tool-detector/controllers/git"
	. "github.com/tinakurian/build-tool-detector/controllers/git"
	"github.com/tinakurian/build-tool-detector/controllers/git/github"
	"net/http"
)

var _ = Describe("GitServiceType", func() {
	Context("GetGitServiceType", func() {
		It("Faulty Host - empty", func() {
			serviceType, err := GetGitServiceType("")
			gomega.Expect(serviceType).Should(gomega.BeNil(), "service type should be 'nil'")
			gomega.Expect(err.StatusCode).Should(gomega.BeEquivalentTo(http.StatusBadRequest), "service type should be 'nil'")
		})

		It("Faulty Host - non-existent", func() {
			serviceType, err := GetGitServiceType("test/test")
			gomega.Expect(serviceType).Should(gomega.BeNil(), "service type should be 'nil'")
			gomega.Expect(err.StatusCode).Should(gomega.BeEquivalentTo(http.StatusBadRequest), "service type should be 'nil'")
		})

		It("Faulty Host - not github.com", func() {
			serviceType, err := GetGitServiceType("http://test.com/test/test")
			gomega.Expect(serviceType).Should(gomega.BeNil(), "service type should be 'nil'")
			gomega.Expect(err.StatusCode).Should(gomega.BeEquivalentTo(http.StatusInternalServerError), "service type should be 'nil'")
		})

		It("Faulty url - no repository", func() {
			serviceType, err := GetGitServiceType("http://github.com/test")
			gomega.Expect(serviceType).Should(gomega.BeNil(), "service type should be 'nil'")
			gomega.Expect(err.StatusCode).Should(gomega.BeEquivalentTo(http.StatusBadRequest), "service type should be 'nil'")
		})

		It("Correct url", func() {
			serviceType, err := GetGitServiceType("http://github.com/test/test")
			gomega.Expect(serviceType).ShouldNot(gomega.BeNil(), "service type should be 'nil'")
			gomega.Expect(err).Should(gomega.BeNil(), "service type should be 'nil'")
		})
	})

	Context("Constants", func() {
		It("Github", func() {
			gomega.Expect(git.GITHUB).Should(gomega.BeEquivalentTo("github"), "service type should be 'nil'")
		})

		It("Unknown", func() {
			gomega.Expect(git.UNKNOWN).Should(gomega.BeEquivalentTo("unknown"), "service type should be 'nil'")
		})
	})

	Context("Service", func() {
		It("Github", func() {
			service := Service{}.GetGitHubService()
			gomega.Expect(service).Should(gomega.BeEquivalentTo(&github.GoooService{}), "service type should be 'nil'")
		})
	})
})
