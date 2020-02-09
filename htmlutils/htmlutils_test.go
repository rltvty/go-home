package htmlutils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/html"

	"github.com/rltvty/go-home/htmlutils"
	. "github.com/rltvty/go-home/htmlutils"
)

var _ = Describe("Htmlutils", func() {
	Describe("RemoveJunk", func() {
		var rootNode *html.Node
		BeforeEach(func() {
			rootNode = GetHTMLFromFile("./test_data/1-jun-12.htm")
			RemoveJunk(rootNode)
		})
		It("should return html with no junk tags", func() {
			formattedHtml, _ := RenderHTMLNode(rootNode)
			Expect(formattedHtml).ToNot(ContainSubstring("<script"))
			Expect(formattedHtml).ToNot(ContainSubstring("<meta"))
			Expect(formattedHtml).ToNot(ContainSubstring("<link"))
			Expect(formattedHtml).ToNot(ContainSubstring("<style"))
		})
	})

	Describe("RemoveEmpty", func() {
		var inNode *html.Node
		var outNode *html.Node
		BeforeEach(func() {
			inNode = htmlutils.GetHTMLFromFile("./test_data/remove-in.htm")
			outNode = htmlutils.GetHTMLFromFile("./test_data/remove-out.htm")
		})
		It("should return html with all empty nodes removed", func() {
			RemoveEmpty(inNode)
			//htmlutils.WriteHTMLToFile(inNode, "./test_data/remove-current.htm")
			inHtml, _ := htmlutils.RenderHTMLNode(inNode)
			outHtml, _ := htmlutils.RenderHTMLNode(outNode)
			Expect(inHtml).To(Equal(outHtml))
		})
	})

	Describe("Squash", func() {
		var inNode *html.Node
		var outNode *html.Node
		BeforeEach(func() {
			inNode = htmlutils.GetHTMLFromFile("./test_data/squash-in.htm")
			outNode = htmlutils.GetHTMLFromFile("./test_data/squash-out.htm")
		})
		It("should return html nodes where single children have been squashed", func() {
			RemoveEmpty(inNode) //this is required for Squash to function properly
			Squash(inNode)
			//htmlutils.WriteHTMLToFile(inNode, "./test_data/squash-current.htm")
			inHtml, _ := htmlutils.RenderHTMLNode(inNode)
			outHtml, _ := htmlutils.RenderHTMLNode(outNode)
			Expect(inHtml).To(Equal(outHtml))
		})
	})

	Describe("CleanClassAttr", func() {
		var inNode *html.Node
		var outNode *html.Node
		BeforeEach(func() {
			inNode = htmlutils.GetHTMLFromFile("./test_data/clean-in.htm")
			outNode = htmlutils.GetHTMLFromFile("./test_data/clean-out.htm")
		})
		It("should return html nodes with a single child have been squashed", func() {
			RemoveEmpty(inNode) //this is required for CleanClassAttr to function properly
			valuesToKeep := []string{
				"island",
				"context__item",
				"segments-list",
				"segment__track",
				"artist",
				"no-margin",
			}
			CleanClassAttr(inNode, valuesToKeep)
			//htmlutils.WriteHTMLToFile(inNode, "./test_data/clean-current.htm")
			inHtml, _ := htmlutils.RenderHTMLNode(inNode)
			outHtml, _ := htmlutils.RenderHTMLNode(outNode)
			Expect(inHtml).To(Equal(outHtml))
		})
	})
})
