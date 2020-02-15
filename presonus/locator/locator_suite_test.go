package locator_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLocator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Locator Suite")
}
