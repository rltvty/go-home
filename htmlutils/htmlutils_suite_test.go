package htmlutils_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHtmlutils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Htmlutils Suite")
}
