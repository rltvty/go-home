package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/rltvty/go-home/downloaders/heidi/bbc/filter"
	"github.com/rltvty/go-home/htmlutils"

	"golang.org/x/net/html"
)

//List of all BBC Radio 1's Residency Episodes (including other hosts than Heidi)
const initialURL = "https://www.bbc.co.uk/programmes/b01d76k4/episodes/guide"
const paginationURL = "https://www.bbc.co.uk/programmes/b01d76k4/episodes/guide?page=%v"

func getNumberOfPages() int {
	rootNode, _ := htmlutils.GetHTMLFromURL(initialURL)
	nodes := htmlutils.FilterHTML([]*html.Node{rootNode}, func(n *html.Node) bool {
		return n.Type == html.ElementNode && htmlutils.IncludesAttr(n, "class", "pagination__page--last")
	})
	nodes = htmlutils.FilterHTML(nodes, func(n *html.Node) bool {
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
		rootNode, _ := htmlutils.GetHTMLFromURL(fmt.Sprintf(paginationURL, i))
		nodes := htmlutils.FilterHTML([]*html.Node{rootNode}, func(n *html.Node) bool {
			if htmlutils.IncludesAttr(n, "class", "br-blocklink__link") {
				html, _ := htmlutils.RenderHTMLNodes([]*html.Node{n})
				return strings.Contains(strings.ToLower(html), strings.ToLower(artist))
			}
			return false
		})
		for _, node := range nodes {
			url := htmlutils.GetAttr(node, "href")
			if url != nil {
				urls = append(urls, *url)
			}
		}
	}
	return urls
}



func main() {
	//urls := getUrls("heidi")
	urls := []string{
		"https://www.bbc.co.uk/programmes/b07pd511",
	}

	for _, url := range urls {
		rootNode, _ := htmlutils.GetHTMLFromURL(url)
		

	}
}
