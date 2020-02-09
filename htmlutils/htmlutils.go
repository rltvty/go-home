package htmlutils

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/yosssi/gohtml"
	"golang.org/x/net/html"
)

//GetHTMLFromURL gets parsed html from an URL
func GetHTMLFromURL(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error getting url: %v", err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error parsing html: %v", err)
	}
	return doc, nil
}

//GetHTMLFromFile gets parsed html from a file
func GetHTMLFromFile(fileName string) *html.Node {
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

//FilterHTML returns a slice of nodes that pass a filter
func FilterHTML(inNodes []*html.Node, filter func(*html.Node) bool) []*html.Node {
	var outNodes []*html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if filter(n) {
			outNodes = append(outNodes, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	for _, inNode := range inNodes {
		f(inNode)
	}
	return outNodes
}

//RenderHTMLNode renders a HTML node into a pretty-printed string
func RenderHTMLNode(node *html.Node) (string, error) {
	buf := new(bytes.Buffer)
	err := html.Render(buf, node)
	return gohtml.Format(buf.String()), err
}

//RenderHTMLNodes renders a set of HTML nodes into a pretty-printed string
func RenderHTMLNodes(nodes []*html.Node) (string, []error) {
	var renderedHTML string
	var errors []error
	for _, node := range nodes {
		ren, err := RenderHTMLNode(node)
		if err != nil {
			errors = append(errors, err)
		}
		renderedHTML += ren
	}
	return renderedHTML, errors
}

//IncludesAttr returns true if the HTML node includes the 'value' in the attrubute 'key'
// example: for <div id="page2" class="class="pagination__page pagination__page--offset2">
//          IncludesAttr(n, "class", "pagination__page") would return true
func IncludesAttr(n *html.Node, key string, val string) bool {
	for _, a := range n.Attr {
		if a.Key == key {
			for _, v := range strings.Split(a.Val, " ") {
				if v == val {
					return true
				}
			}
		}
	}
	return false
}

//GetAttr returns the attribute value for the key if it exists
func GetAttr(n *html.Node, key string) *string {
	for _, a := range n.Attr {
		if a.Key == key {
			return &a.Val
		}
	}
	return nil
}

func nodeType(n *html.Node) string {
	switch n.Type {
	case html.ElementNode:
		return "ElementNode"
	case html.ErrorNode:
		return "ErrorNode"
	case html.CommentNode:
		return "CommentNode"
	case html.DoctypeNode:
		return "DoctypeNode"
	case html.DocumentNode:
		return "DocumentNode"
	case html.TextNode:
		return "TextNode"
	default:
		return "UnknownNodeType"
	}
}

//DebugNode prints info about a node
func DebugNode(prefix string, n *html.Node) {
	log.Printf("%s %v '%s' %v", prefix, nodeType(n), strings.TrimSpace(n.Data), n.Attr)
}

//DebugTree prints the current node tree to stdout
func DebugTree(rootNode *html.Node) {
	var f func(*html.Node, string)
	f = func(n *html.Node, offset string) {
		DebugNode(offset, n)

		//iterate on children
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, offset+"  ")
		}
	}
	f(rootNode, "")
}

//RemoveEmpty removes nodes that have no text content
func RemoveEmpty(rootNode *html.Node) {
	for {
		nodesToRemove := []*html.Node{}
		var f func(*html.Node)
		f = func(n *html.Node) {

			if n.Type == html.TextNode && strings.TrimSpace(n.Data) == "" {
				//remove empty text nodes
				nodesToRemove = append(nodesToRemove, n)
			} else if n.Type == html.ElementNode && n.FirstChild == nil {
				//remove element nodes with no children
				nodesToRemove = append(nodesToRemove, n)
			}

			//iterate on children
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(rootNode)

		for _, nodeToRemove := range nodesToRemove {
			nodeToRemove.Parent.RemoveChild(nodeToRemove)
		}
		if len(nodesToRemove) == 0 {
			break
		}
	}
}

//Squash combines nodes that only have a single child, grouping attributes together
func Squash(rootNode *html.Node) {
	isSquashableNode := func(parent *html.Node) bool {
		if parent.Type != html.ElementNode {
			return false
		}
		firstChild := parent.FirstChild
		if firstChild == nil || firstChild.Type != html.ElementNode {
			return false
		}
		secondChild := firstChild.NextSibling
		return secondChild == nil
	}

	var nodesToSquash []*html.Node
	var findNodesToSquash func(*html.Node)
	findNodesToSquash = func(n *html.Node) {
		if isSquashableNode(n) {
			nodesToSquash = append(nodesToSquash, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findNodesToSquash(c)
		}
	}
	findNodesToSquash(rootNode)

	for i := len(nodesToSquash); i > 0; i-- {
		nodeToSquash := nodesToSquash[i-1]
		if nodeToSquash.FirstChild != nil {
			onlyChild := nodeToSquash.FirstChild
			nodeToSquash.RemoveChild(onlyChild)

			grandChildren := []*html.Node{}
			for grandChild := onlyChild.FirstChild; grandChild != nil; grandChild = grandChild.NextSibling {
				grandChildren = append(grandChildren, grandChild)
			}

			for _, grandChild := range grandChildren {
				onlyChild.RemoveChild(grandChild)
				nodeToSquash.AppendChild(grandChild)
			}
		}
	}

	//TODO: merge onlyChilds Attr into parents Attr
}
