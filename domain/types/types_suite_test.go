package types_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBuildtype(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Types Suite")
}
