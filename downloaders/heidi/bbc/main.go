package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/rltvty/go-home/netutils"
	"golang.org/x/net/html"
)

//List of all BBC Radio 1's Residency Episodes (including other hosts than Heidi)
const initialURL = "https://www.bbc.co.uk/programmes/b01d76k4/episodes/guide"
const paginationURL = "https://www.bbc.co.uk/programmes/b01d76k4/episodes/guide?page=%v"

func getNumberOfPages() int {
	rootNode, _ := netutils.GetHTMLFromURL(initialURL)
	nodes := netutils.FilterHTML([]*html.Node{rootNode}, func(n *html.Node) bool {
		return n.Type == html.ElementNode && netutils.IncludesAttr(n, "class", "pagination__page--last")
	})
	nodes = netutils.FilterHTML(nodes, func(n *html.Node) bool {
		return n.Data == "a"
	})
	if len(nodes) != 1 {
		log.Fatal("Expected to find one 'a' node, but found zero or many")
	}
	lastPage, err := strconv.ParseInt(nodes[0].FirstChild.Data, 10, 0)
	if err != nil {
		log.Fatalf("Error parsing last page number: %s", err)
	}
	return int(lastPage)
}

func getUrls(artist string) []string {
	urls := []string{}
	lastPage := getNumberOfPages()
	for i := 1; i <= lastPage; i++ {
		rootNode, _ := netutils.GetHTMLFromURL(fmt.Sprintf(paginationURL, i))
		nodes := netutils.FilterHTML([]*html.Node{rootNode}, func(n *html.Node) bool {
			if netutils.IncludesAttr(n, "class", "br-blocklink__link") {
				html, _ := netutils.RenderHTMLNodes([]*html.Node{n})
				return strings.Contains(strings.ToLower(html), strings.ToLower(artist))
			}
			return false
		})
		for _, node := range nodes {
			url := netutils.GetAttr(node, "href")
			if url != nil {
				urls = append(urls, *url)
			}
		}
	}
	return urls
}

type track struct {
	artist string
	name   string
	album  string
	label  string
	info   string
}

func main() {
	//urls := getUrls("heidi")
	urls := []string{
		"https://www.bbc.co.uk/programmes/b07pd511",
	}

	for _, url := range urls {
		rootNode, _ := netutils.GetHTMLFromURL(url)
		trackNodes := netutils.FilterHTML([]*html.Node{rootNode}, func(n *html.Node) bool {
			return netutils.IncludesAttr(n, "class", "segment__track")
		})
		for _, trackNode := range trackNodes {
			trackHTML, _ := netutils.RenderHTMLNodes([]*html.Node{trackNode})
			fmt.Println("**** Original: ")
			fmt.Println(trackHTML)
			artistNodes := netutils.FilterHTML([]*html.Node{trackNode}, func(n *html.Node) bool {
				return netutils.IncludesAttr(n, "class", "artist")
			})
			nameBlocks := netutils.FilterHTML([]*html.Node{trackNode}, func(n *html.Node) bool {
				return n.Data == "p" && netutils.IncludesAttr(n, "class", "no-margin")
			})
			nameNodes := netutils.FilterHTML(nameBlocks, func(n *html.Node) bool {
				return n.Data == "span"
			})
			t := track{
				artist: artistNodes[0].FirstChild.Data,
				name:   nameNodes[0].FirstChild.Data,
			}

			labelNodes := netutils.FilterHTML([]*html.Node{trackNode}, func(n *html.Node) bool {
				return n.Data == "abbr"
			})

			if len(labelNodes) > 0 {
				t.label = labelNodes[0].FirstChild.Data
			}

			for _, node := range append(artistNodes, append(nameNodes, labelNodes...)...) {
				node.Parent.RemoveChild(node)
			}

			trackHTML, _ = netutils.RenderHTMLNodes([]*html.Node{trackNode})
			fmt.Println("**** Trimmed: ")
			fmt.Println(trackHTML)

			infoNodes := netutils.FilterHTML([]*html.Node{trackNode}, func(n *html.Node) bool {
				return n.Type == html.TextNode
			})

			infoHTML, _ := netutils.RenderHTMLNodes(infoNodes)
			fmt.Println(infoHTML)

			fmt.Println(t)
			fmt.Println()
		}

	}
}
