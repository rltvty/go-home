package htmlutils_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/html"

	"github.com/rltvty/go-home/htmlutils"
	. "github.com/rltvty/go-home/htmlutils"
)

var _ = Describe("Htmlutils", func() {
	Describe("RemoveEmpty", func() {
		var inNode *html.Node
		var inHtml string
		var outNode *html.Node
		BeforeEach(func() {
			inNode = htmlutils.GetHTMLFromFile("./test_data/remove-in.htm")
			outNode = htmlutils.GetHTMLFromFile("./test_data/remove-out.htm")
		})
		It("should return html with all empty nodes removed", func() {
			RemoveEmpty(inNode)
			inHtml, _ = htmlutils.RenderHTMLNode(inNode)
			outHtml, _ := htmlutils.RenderHTMLNode(outNode)
			Expect(inHtml).To(Equal(outHtml))
		})
		/*
			AfterEach(func() {
				buf := new(bytes.Buffer)
				io.WriteString(buf, inHtml)
				ioutil.WriteFile("./test_data/remove-current.htm", buf.Bytes(), os.ModePerm)
			})
		*/
	})

	Describe("Squash", func() {
		var inNode *html.Node
		var inHtml string
		var outNode *html.Node
		BeforeEach(func() {
			inNode = htmlutils.GetHTMLFromFile("./test_data/squash-in.htm")
			outNode = htmlutils.GetHTMLFromFile("./test_data/squash-out.htm")
		})
		It("should return html nodes with a single child have been squashed", func() {
			RemoveEmpty(inNode) //this is required for Squash to function properly
			Squash(inNode)
			inHtml, _ = htmlutils.RenderHTMLNode(inNode)
			outHtml, _ := htmlutils.RenderHTMLNode(outNode)
			Expect(inHtml).To(Equal(outHtml))
		})
		AfterEach(func() {
			buf := new(bytes.Buffer)
			io.WriteString(buf, inHtml)
			ioutil.WriteFile("./test_data/squash-current.htm", buf.Bytes(), os.ModePerm)
		})
	})
})
