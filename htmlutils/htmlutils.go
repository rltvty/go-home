package htmlutils

import (
	"bytes"
	"fmt"
	"net/http"
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
