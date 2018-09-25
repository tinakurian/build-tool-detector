/*

Package buildtype_test tests
the buildtype package.

*/
package buildtype_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/tinakurian/build-tool-detector/controllers/buildtype"
)

var _ = Describe("BuildToolType", func() {
	Context("Maven", func() {
		It("Get Maven", func() {
			maven := Maven()
			Expect(maven.BuildToolType).Should(BeEquivalentTo("maven"), "build tool type type should be 'maven'")
		})
	})

	Context("Unknown", func() {
		It("Get Unknown", func() {
			unknown := Unknown()
			Expect(unknown.BuildToolType).Should(BeEquivalentTo("unknown"), "build tool type should be 'unknown'")
		})
	})
})
