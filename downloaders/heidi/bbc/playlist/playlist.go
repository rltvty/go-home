package playlist

import (
	"github.com/rltvty/go-home/htmlutils"
	"golang.org/x/net/html"
)

//RemoveBBCNotNeeded removes a bunch of auxillary nodes from BBC pages
func RemoveBBCNotNeeded(rootNode *html.Node) {
	notNeededNodes := htmlutils.FilterHTML([]*html.Node{rootNode}, func(n *html.Node) bool {
		for _, tag := range []string{"header", "footer", "img", "svg", "aside"} {
			if n.Data == tag {
				return true
			}
		}
		type keyVal struct {
			key string
			val string
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
	htmlutils.RemoveNodes(notNeededNodes)
}

//SquashAndClean removes all the junk we don't need for the playlists
func SquashAndClean(rootNode *html.Node) {
	htmlutils.RemoveJunk(rootNode)
	RemoveBBCNotNeeded(rootNode)
	htmlutils.RemoveEmpty(rootNode)
	htmlutils.Squash(rootNode)
	valuesToKeep := []string{
		"artist",
		"broadcast-event__date",
		"broadcast-event__time",
		"context__item",
		"episode-panel__meta",
		"inline",
		"island",
		"micro",
		"no-margin",
		"programme__title",
		"programme__service",
		"segment__track",
		"segments-list",
		"synopsis-toggle__long",
		"synopsis-toggle__short",
	}
	htmlutils.CleanClassAttr(rootNode, valuesToKeep)
}

type track struct {
	artist string
	name   string
	album  string
	label  string
	info   string
}

//GetPlaylist extracts a playlist from html
func GetPlaylist(rootNode *html.Node) string {
	SquashAndClean(rootNode)
	/*
		trackNodes := htmlutils.FilterHTML([]*html.Node{rootNode}, func(n *html.Node) bool {
			return htmlutils.IncludesAttr(n, "class", "segment__track")
		})

		for _, trackNode := range trackNodes {
			trackHTML, _ := htmlutils.RenderHTMLNodes([]*html.Node{trackNode})
			fmt.Println("**** Original: ")
			fmt.Println(trackHTML)
			artistNodes := htmlutils.FilterHTML([]*html.Node{trackNode}, func(n *html.Node) bool {
				return htmlutils.IncludesAttr(n, "class", "artist")
			})
			nameBlocks := htmlutils.FilterHTML([]*html.Node{trackNode}, func(n *html.Node) bool {
				return n.Data == "p" && htmlutils.IncludesAttr(n, "class", "no-margin")
			})
			nameNodes := htmlutils.FilterHTML(nameBlocks, func(n *html.Node) bool {
				return n.Data == "span"
			})
			t := track{
				artist: artistNodes[0].FirstChild.Data,
				name:   nameNodes[0].FirstChild.Data,
			}

			labelNodes := htmlutils.FilterHTML([]*html.Node{trackNode}, func(n *html.Node) bool {
				return n.Data == "abbr"
			})

			if len(labelNodes) > 0 {
				t.label = labelNodes[0].FirstChild.Data
			}

			for _, node := range append(artistNodes, append(nameNodes, labelNodes...)...) {
				node.Parent.RemoveChild(node)
			}

			trackHTML, _ = htmlutils.RenderHTMLNodes([]*html.Node{trackNode})
			fmt.Println("**** Trimmed: ")
			fmt.Println(trackHTML)

			infoNodes := htmlutils.FilterHTML([]*html.Node{trackNode}, func(n *html.Node) bool {
				return n.Type == html.TextNode
			})

			infoHTML, _ := htmlutils.RenderHTMLNodes(infoNodes)
			fmt.Println(infoHTML)

			fmt.Println(t)
			fmt.Println()
		}
	*/
	return ""
}
