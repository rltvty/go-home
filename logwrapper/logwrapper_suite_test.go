package logwrapper_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLogwrapper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Logwrapper Suite")
}
