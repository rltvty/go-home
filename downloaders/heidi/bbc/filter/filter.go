package filter

import (
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
