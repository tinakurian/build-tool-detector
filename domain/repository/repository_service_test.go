/*

Package git_test is used to test the functionality
within the git package.

*/
package repository_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tinakurian/build-tool-detector/domain/repository/github"
	"github.com/tinakurian/build-tool-detector/domain/repository"
)

var _ = Describe("GitServiceType", func() {
	Context("CreateService", func() {
		It("Faulty Host - empty", func() {
			serviceType, err := repository.CreateService("")
			Expect(serviceType).Should(BeNil(), "service type should be 'nil'")
			Expect(err.Error()).Should(BeEquivalentTo(github.ErrInvalidPath.Error()), "service type should be '400'")
		})

		It("Faulty Host - non-existent", func() {
			serviceType, err := repository.CreateService("test/test")
			Expect(serviceType).Should(BeNil(), "service type should be 'nil'")
			Expect(err.Error()).Should(BeEquivalentTo(github.ErrInvalidPath.Error()), "service type should be '400'")
		})

		It("Faulty Host - not github.com", func() {
			serviceType, err := repository.CreateService("http://test.com/test/test")
			Expect(serviceType).Should(BeNil(), "service type should be 'nil'")
			Expect(err.Error()).Should(BeEquivalentTo(github.ErrUnsupportedService.Error()), "service type should be '500'")
		})

		It("Faulty url - no repository", func() {
			serviceType, err := repository.CreateService("http://github.com/test")
			Expect(serviceType).Should(BeNil(), "service type should be 'nil'")
			Expect(err.Error()).Should(BeEquivalentTo(github.ErrUnsupportedGithubURL.Error()), "service type should be '400'")
		})

		It("Correct url - non-existent", func() {
			serviceType, err := repository.CreateService("http://github.com/test/test")
			Expect(serviceType).ShouldNot(BeNil(), "service type should be not be'nil'")
			Expect(err).Should(BeNil(), "err should be 'nil'")
		})
	})

})
