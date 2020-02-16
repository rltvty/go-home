package testhelpers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/rltvty/go-home/testhelpers"
)

var _ = Describe("Testhelpers", func() {
	Describe("IsJSON", func() {
		It("should verify json correctly", func() {
			Expect(IsJSON(`Australia`)).To(BeFalse())
			Expect(IsJSON(`{"Country":"Australia"}`)).To(BeTrue())
		})
	})
})
