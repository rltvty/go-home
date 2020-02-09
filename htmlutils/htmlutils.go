package htmlutils

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/yosssi/gohtml"
	"golang.org/x/net/html"
)

//RemoveJunk removes all script, meta, link, style, comment nodes
func RemoveJunk(rootNode *html.Node) {
	junkNodes := FilterHTML([]*html.Node{rootNode}, func(n *html.Node) bool {
		if n.Type == html.CommentNode {
			return true
		}

		for _, tag := range []string{"script", "noscript", "meta", "link", "style"} {
			if n.Data == tag {
				return true
			}
		}
		return false
	})
	RemoveNodes(junkNodes)
}

//RemoveNodes removes the passed nodes from their parent nodes
func RemoveNodes(nodes []*html.Node) {
	for _, node := range nodes {
		node.Parent.RemoveChild(node)
	}
}

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

//WriteHTMLToFile writes formatted html to a file
func WriteHTMLToFile(rootNode *html.Node, fileName string) error {
	buf := new(bytes.Buffer)
	renderedHTML, err := RenderHTMLNode(rootNode)
	if err != nil {
		return err
	}
	_, err = io.WriteString(buf, renderedHTML)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fileName, buf.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
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

func mergeAttr(keep *html.Node, merge *html.Node) {
	for _, mergeAttr := range merge.Attr {
		if strings.TrimSpace(mergeAttr.Val) != "" {
			foundKeepAttr := false
			for i, keepAttr := range keep.Attr {
				if keepAttr.Key == mergeAttr.Key {
					foundKeepAttr = true
					keep.Attr[i].Val = strings.TrimSpace(fmt.Sprintf("%s %s", keep.Attr[i].Val, mergeAttr.Val))
				}
			}
			if !foundKeepAttr {
				keep.Attr = append(keep.Attr, mergeAttr)
			}
		}
	}
}

//CleanClassAttr removes all class attribute values that arent in the valuesToKeep list
func CleanClassAttr(rootNode *html.Node, valuesToKeep []string) {
	keep := map[string]string{}
	for _, valueToKeep := range valuesToKeep {
		keep[valueToKeep] = ""
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		newAttrs := []html.Attribute{}
		for _, oldAttr := range n.Attr {
			if oldAttr.Key == "class" {
				oldValues := strings.Split(oldAttr.Val, " ")
				newValues := []string{}
				for _, oldValue := range oldValues {
					if _, found := keep[oldValue]; found {
						newValues = append(newValues, oldValue)
					}
				}
				oldAttr.Val = strings.Join(newValues, " ")
				if oldAttr.Val != "" {
					newAttrs = append(newAttrs, oldAttr)
				}
			} else {
				newAttrs = append(newAttrs, oldAttr)
			}
		}
		n.Attr = newAttrs

		//iterate on children
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(rootNode)
}

//Squash combines nodes that only have a single child, grouping attributes together
func Squash(rootNode *html.Node) {
	isSquashableNode := func(parent *html.Node) bool {
		if parent.Type != html.ElementNode || parent.Data == "head" {
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
			mergeAttr(nodeToSquash, onlyChild)

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
}
