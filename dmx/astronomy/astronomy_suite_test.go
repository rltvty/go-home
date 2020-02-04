package astronomy_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAstronomy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Astronomy Suite")
}
