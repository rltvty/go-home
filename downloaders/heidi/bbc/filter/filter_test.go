package filter_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/html"

	. "github.com/rltvty/go-home/downloaders/heidi/bbc/filter"
	"github.com/rltvty/go-home/htmlutils"
)

func getHtmlFromFile(fileName string) *html.Node {
	reader, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("cannot open test file: %s", err)
	}
	rootNode, err := html.Parse(reader)
	if err != nil {
		log.Fatalf("cannot parse html in test file: %s", err)
	}
	return rootNode
}

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
			rootNodes = append(rootNodes, getHtmlFromFile(testFile))
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

				/*
					if i == 0 {
						buf := new(bytes.Buffer)
						io.WriteString(buf, formattedHtml)
						ioutil.WriteFile("./test_data/sofar.htm", buf.Bytes(), os.ModePerm)
					}
				*/
			}
		})
	})

	Describe("RemoveEmpty", func() {
		var inNode *html.Node
		var inHtml string
		var outNode *html.Node
		BeforeEach(func() {
			inNode = getHtmlFromFile("./test_data/remove-in.htm")
			outNode = getHtmlFromFile("./test_data/remove-out.htm")
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
			inNode = getHtmlFromFile("./test_data/squash-in.htm")
			outNode = getHtmlFromFile("./test_data/squash-out.htm")
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
