package filter

import (
	"log"
	"strings"

	"github.com/rltvty/go-home/htmlutils"
	"golang.org/x/net/html"
)

//RemoveJunk removes all script, meta, link, style, comment nodes
func RemoveJunk(rootNode *html.Node) {
	junkNodes := htmlutils.FilterHTML([]*html.Node{rootNode}, func(n *html.Node) bool {
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

type keyVal struct {
	key string
	val string
}

//RemoveBBCNotNeeded removes a bunch of auxillary nodes from BBC pages
func RemoveBBCNotNeeded(rootNode *html.Node) {
	notNeededNodes := htmlutils.FilterHTML([]*html.Node{rootNode}, func(n *html.Node) bool {
		for _, tag := range []string{"header", "footer", "img", "svg", "aside"} {
			if n.Data == tag {
				return true
			}
		}
		notNeededAttrs := []keyVal{
			keyVal{"class", "segment__buttons"},
			keyVal{"class", "br-masthead__main"},
			keyVal{"class", "episode-playout"},
			keyVal{"data-map-column", "more"},
			keyVal{"for", "segments-moreless"},
			keyVal{"id", "programmes-footer"},
			keyVal{"id", "broadcasts"},
			keyVal{"id", "br-nav-programme"},
			keyVal{"role", "button"},
		}

		for _, keyVal := range notNeededAttrs {
			if htmlutils.IncludesAttr(n, keyVal.key, keyVal.val) {
				return true
			}
		}
		return false
	})
	RemoveNodes(notNeededNodes)
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

func debugNode(prefix string, n *html.Node) {
	log.Printf("%s %v '%s' %v", prefix, nodeType(n), strings.TrimSpace(n.Data), n.Attr)
}

//DebugTree prints the current node tree to stdout
func DebugTree(rootNode *html.Node) {
	var f func(*html.Node, string)
	f = func(n *html.Node, offset string) {
		debugNode(offset, n)

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

	log.Println()
	log.Println("**** Starting squashing****")

	for i := len(nodesToSquash); i > 0; i-- {
		log.Println()
		nodeToSquash := nodesToSquash[i-1]
		if nodeToSquash.FirstChild != nil {
			log.Println(" start: ")
			debugNode("node to squash: ", nodeToSquash)
			onlyChild := nodeToSquash.FirstChild
			debugNode("removing onlyChild: ", onlyChild)
			nodeToSquash.RemoveChild(onlyChild)

			grandChildren := []*html.Node{}
			for grandChild := onlyChild.FirstChild; grandChild != nil; grandChild = grandChild.NextSibling {
				grandChildren = append(grandChildren, grandChild)
			}

			for _, grandChild := range grandChildren {
				debugNode("moving grandChild: ", grandChild)
				onlyChild.RemoveChild(grandChild)
				nodeToSquash.AppendChild(grandChild)
			}
			log.Println()
			log.Println(" new tree: ")
			DebugTree(nodeToSquash)
		}
	}

	//TODO: merge onlyChilds Attr into parents Attr

	DebugTree(rootNode)
}
