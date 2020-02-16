package testhelpers_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTesthelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testhelpers Suite")
}
