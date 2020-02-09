package filter_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/html"

	. "github.com/rltvty/go-home/downloaders/heidi/bbc/filter"
	"github.com/rltvty/go-home/htmlutils"
)

var _ = Describe("Filter", func() {
	testFiles := []string{
		"./test_data/1-jun-12.htm",
		"./test_data/2-aug-13.htm",
		"./test_data/15-aug-16.htm",
		"./test_data/26-may-16.htm",
	}
	rootNodes := []*html.Node{}
	BeforeEach(func() {
		for _, testFile := range testFiles {
			rootNodes = append(rootNodes, htmlutils.GetHTMLFromFile(testFile))
		}
	})

	Describe("Remove Junk", func() {
		BeforeEach(func() {
			for _, rootNode := range rootNodes {
				RemoveJunk(rootNode)
			}
		})
		It("should return html with no junk tags", func() {
			Expect(len(rootNodes)).To(Equal(4))
			for _, rootNode := range rootNodes {
				formattedHtml, _ := htmlutils.RenderHTMLNode(rootNode)
				Expect(formattedHtml).ToNot(ContainSubstring("<script"))
				Expect(formattedHtml).ToNot(ContainSubstring("<meta"))
				Expect(formattedHtml).ToNot(ContainSubstring("<link"))
				Expect(formattedHtml).ToNot(ContainSubstring("<style"))
			}
		})
	})

	Describe("Remove NotNeeded", func() {
		BeforeEach(func() {
			for _, rootNode := range rootNodes {
				RemoveJunk(rootNode)
				RemoveBBCNotNeeded(rootNode)
			}
		})
		It("should return html with no bbc header, footer, navigation, images", func() {
			for _, rootNode := range rootNodes {
				formattedHtml, _ := htmlutils.RenderHTMLNode(rootNode)
				Expect(formattedHtml).ToNot(ContainSubstring("<img"))
				Expect(formattedHtml).ToNot(ContainSubstring("<header"))
				Expect(formattedHtml).ToNot(ContainSubstring("<footer"))
				Expect(formattedHtml).ToNot(ContainSubstring("<nav"))
			}
		})
	})
})
