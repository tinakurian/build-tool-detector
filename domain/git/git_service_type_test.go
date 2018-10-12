/*

Package git_test is used to test the functionality
within the git package.

*/
package git_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tinakurian/build-tool-detector/domain/git"
	"github.com/tinakurian/build-tool-detector/domain/git/github"
)

var _ = Describe("GitServiceType", func() {
	Context("GetGitServiceType", func() {
		It("Faulty Host - empty", func() {
			serviceType, errz := git.GetGitServiceType("")
			err := *errz
			Expect(serviceType).Should(BeNil(), "service type should be 'nil'")
			Expect(err.Error()).Should(BeEquivalentTo(github.ErrInvalidPath.Error()), "service type should be '400'")
		})

		It("Faulty Host - non-existent", func() {
			serviceType, errz := git.GetGitServiceType("test/test")
			err := *errz
			Expect(serviceType).Should(BeNil(), "service type should be 'nil'")
			Expect(err.Error()).Should(BeEquivalentTo(github.ErrInvalidPath.Error()), "service type should be '400'")
		})

		It("Faulty Host - not github.com", func() {
			serviceType, errz := git.GetGitServiceType("http://test.com/test/test")
			err := *errz
			Expect(serviceType).Should(BeNil(), "service type should be 'nil'")
			Expect(err.Error()).Should(BeEquivalentTo(github.ErrUnsupportedService.Error()), "service type should be '500'")
		})

		It("Faulty url - no repository", func() {
			serviceType, errz := git.GetGitServiceType("http://github.com/test")
			err := *errz
			Expect(serviceType).Should(BeNil(), "service type should be 'nil'")
			Expect(err.Error()).Should(BeEquivalentTo(github.ErrUnsupportedGithubURL.Error()), "service type should be '400'")
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
			service := git.Service{}.GetGitHubService("test", "test")
			Expect(service).Should(BeEquivalentTo(&github.GitService{ClientID: "test", ClientSecret: "test"}), "service type should be a git service")
		})
	})
})
