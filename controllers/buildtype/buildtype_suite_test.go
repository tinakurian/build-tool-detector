package buildtype_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBuildtype(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Buildtype Suite")
}
