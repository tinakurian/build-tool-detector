package github_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	git "github.com/tinakurian/build-tool-detector/controllers/git"
	"github.com/tinakurian/build-tool-detector/controllers/git/github"
	"net/http"
)

var _ = Describe("Attributes", func() {

	Context("build-tool-detector/controllers/git/github", func() {
		It("attributes - GetAttributes() use branch from url", func() {
			_, gitService := git.GetServiceType("https://github.com/owner/repo/tree/branch")
			_, attributes := github.GetAttributes(gitService.Segments, nil)
			gomega.Expect(attributes.Owner).Should(gomega.BeEquivalentTo("owner"), "Owner field should be equivalent to 'owner'")
			gomega.Expect(attributes.Repository).Should(gomega.BeEquivalentTo("repo"), "Repository field should be equivalent to 'repo'")
			gomega.Expect(attributes.Branch).Should(gomega.BeEquivalentTo("branch"), "Branch field should be equivalent to 'branch'")

		})

		It("attributes - GetAttributes() use default branch master", func() {
			_, gitService := git.GetServiceType("https://github.com/owner/repo")
			_, attributes := github.GetAttributes(gitService.Segments, nil)
			gomega.Expect(attributes.Owner).Should(gomega.BeEquivalentTo("owner"), "Owner field should be equivalent to 'owner'")
			gomega.Expect(attributes.Repository).Should(gomega.BeEquivalentTo("repo"), "Repository field should be equivalent to 'repo'")
			gomega.Expect(attributes.Branch).Should(gomega.BeEquivalentTo("master"), "Branch field should be equivalent to 'master'")
		})

		It("attributes - GetAttributes() not enough segments extracted from path", func() {
			segments := []string{"test1", "test2"}
			err, attributes := github.GetAttributes(segments, nil)
			gomega.Expect(err.StatusCode).Should(gomega.BeEquivalentTo(http.StatusBadRequest), "invalid url segments should be equivalent to a Bad Request")
			gomega.Expect(attributes).Should(gomega.Equal(github.Attributes{}), "invalid url segments cause empty attributes")
		})
	})
})
