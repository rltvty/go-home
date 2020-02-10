package playlist

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/rltvty/go-home/htmlutils"
	"golang.org/x/net/html"
)

//RemoveBBCNotNeeded removes a bunch of auxillary nodes from BBC pages
func RemoveBBCNotNeeded(rootNode *html.Node) {
	notNeededNodes := htmlutils.FilterHTML([]*html.Node{rootNode}, func(n *html.Node) bool {
		for _, tag := range []string{"header", "footer", "img", "svg", "aside", "head"} {
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
		"longest-synopsis",
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

//Track contains info about a music track
type Track struct {
	Artist string
	Name   string
	Album  string
	Label  string
	Info   string
}

func (t Track) String() string {
	return fmt.Sprintf("Artist: %s\nTrack: %s\nAlbum: %s\nLabel: %s\nInfo: %s\n", t.Artist, t.Name, t.Album, t.Label, t.Info)
}

type playlist struct {
	title string
}

//ParseTrack pulls track info from a track node
func ParseTrack(trackNode *html.Node) Track {
	fmt.Println("**** Original: ")
	htmlutils.DebugTree(trackNode)

	var t Track

	artistTexts := getText(trackNode, func(n *html.Node) bool {
		return n.Data == "h3" && htmlutils.IncludesAttr(n, "class", "artist")
	})
	t.Artist = strings.Join(artistTexts, " ")

	nameTexts := getText(trackNode, func(n *html.Node) bool {
		return n.Data == "p" && htmlutils.IncludesAttr(n, "class", "no-margin")
	})
	t.Name = strings.Join(nameTexts, " ")

	t.Label = getFirstText(trackNode, func(n *html.Node) bool {
		return htmlutils.IncludesAttr(n, "title", "label")
	})
	t.Label = strings.TrimSuffix(t.Label, ".")

	album := getFirstText(trackNode, func(n *html.Node) bool {
		return htmlutils.IncludesAttr(n, "class", "inline")
	})
	album = strings.TrimSuffix(t.Label, ".")
	if album != t.Label {
		t.Album = album
	}
	info := getLastText(trackNode, func(n *html.Node) bool {
		return htmlutils.IncludesAttr(n, "class", "segment__track")
	})
	info = strings.TrimSuffix(info, ".")

	if info != t.Artist && info != t.Name && info != t.Label && info != t.Album {
		t.Info = info
	}

	fmt.Print(t)
	fmt.Println()

	return t
}

func getText(rootNode *html.Node, filter func(n *html.Node) bool) []string {
	out := []string{}
	nodes := htmlutils.FilterHTML([]*html.Node{rootNode}, filter)

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			out = append(out, strings.TrimSpace(n.Data))
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	for _, node := range nodes {
		f(node)
	}
	return out
}

func getFirstText(rootNode *html.Node, filter func(n *html.Node) bool) string {
	text := getText(rootNode, filter)
	if len(text) > 0 {
		return text[0]
	}
	return ""
}

func getLastText(rootNode *html.Node, filter func(n *html.Node) bool) string {
	text := getText(rootNode, filter)
	if len(text) > 0 {
		return text[len(text)-1]
	}
	return ""
}

func getDescription(rootNode *html.Node) string {
	detailNodes := htmlutils.FilterHTML([]*html.Node{rootNode}, func(n *html.Node) bool {
		return htmlutils.IncludesAttr(n, "data-map-column", "playout-details")
	})

	for _, detailNode := range detailNodes {
		detailNode.Parent.RemoveChild(detailNode)
		//log.Println()
		//log.Println("Details tree: ")
		//htmlutils.DebugTree(detailNode)
		//log.Println()
		description := getText(detailNode, func(n *html.Node) bool {
			return htmlutils.IncludesAttr(n, "class", "longest-synopsis")
		})
		if len(description) == 0 {
			description = getText(detailNode, func(n *html.Node) bool {
				return htmlutils.IncludesAttr(n, "class", "synopsis-toggle__long")
			})
		}
		if len(description) == 0 {
			description = getText(detailNode, func(n *html.Node) bool {
				return htmlutils.IncludesAttr(n, "class", "synopsis-toggle__short")
			})
		}
		return strings.Join(description, " ")
	}
	return ""
}

func getDate(rootNode *html.Node) string {
	dateNodes := htmlutils.FilterHTML([]*html.Node{rootNode}, func(n *html.Node) bool {
		return htmlutils.IncludesAttr(n, "class", "broadcast-event__time")
	})

	for _, dateNode := range dateNodes {
		dateNode.Parent.RemoveChild(dateNode)
		//log.Println()
		//log.Println("Date tree: ")
		//htmlutils.DebugTree(dateNode)
		//log.Println()
		dateText := htmlutils.GetAttr(dateNode, "content")
		if dateText != nil {
			dateTime, _ := time.Parse(time.RFC3339, *dateText)
			//log.Println(dateTime)
			return dateTime.Format("2006-01-02")
		}
	}
	return ""
}

//GetPlaylist extracts a playlist from html
func GetPlaylist(rootNode *html.Node) string {
	SquashAndClean(rootNode)

	log.Println()
	log.Println("*** Starting Playlist")
	log.Printf("Description: %s", getDescription(rootNode))
	log.Printf("Date: %s", getDate(rootNode))

	trackNodes := htmlutils.FilterHTML([]*html.Node{rootNode}, func(n *html.Node) bool {
		return htmlutils.IncludesAttr(n, "class", "segment__track")
	})

	for _, trackNode := range trackNodes {
		trackNode.Parent.RemoveChild(trackNode)
		ParseTrack(trackNode)
	}

	println()
	println("Remaining tree: ")
	htmlutils.DebugTree(rootNode)
	return ""
}
