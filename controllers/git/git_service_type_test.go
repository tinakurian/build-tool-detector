package git_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/tinakurian/build-tool-detector/controllers/git"
)

var _ = Describe("GitServiceType", func() {

	Context("build-tool-detector/controllers/git", func() {
		It("git_service_type - GetServiceType() service github", func() {
			_, gitService := git.GetServiceType("https://github.com/owner/repo/tree/branch")

			gomega.Expect(gitService.Service).Should(gomega.BeEquivalentTo("github"), "git service should be equivalent to 'github'")
			gomega.Expect(gitService.Segments[1]).Should(gomega.BeEquivalentTo("owner"), "1st segment from url should be 'owner'")
			gomega.Expect(gitService.Segments[2]).Should(gomega.BeEquivalentTo("repo"), "2nd segment from url should be 'repo'")
			gomega.Expect(gitService.Segments[3]).Should(gomega.BeEquivalentTo("tree"), "third segment from url should be 'tree'")
			gomega.Expect(gitService.Segments[4]).Should(gomega.BeEquivalentTo("branch"), "fourth segment from url should be 'branch'")
		})

		It("git_service_type - GetServiceType() service unknown", func() {
			_, gitService := git.GetServiceType("https://test.com/test/test/tree/master")
			gomega.Expect(gitService.Service).Should(gomega.BeEquivalentTo("unknown"), "git service should be equivalent to 'unknown'")
		})

		It("git_service_type - GetServiceType() bad request with no owner or repository", func() {
			err, _ := git.GetServiceType("https://test.com")
			gomega.Expect(err.StatusCode).Should(gomega.BeEquivalentTo(http.StatusBadRequest), "git service should be equivalent to 'http.StatusBadRequest'")
		})

		It("git_service_type - GetServiceType() bad request with no schema or host", func() {
			err, _ := git.GetServiceType("test/test/test")
			gomega.Expect(err.StatusCode).Should(gomega.BeEquivalentTo(http.StatusBadRequest), "git service should be equivalent to 'http.StatusBadRequest'")
		})

		It("git_service_type - GetTyGetServiceTypepe() bad request whitespace url", func() {
			err, _ := git.GetServiceType(" ")
			gomega.Expect(err.StatusCode).Should(gomega.BeEquivalentTo(http.StatusBadRequest), "git service should be equivalent to 'http.StatusBadRequest'")
		})
	})
})
