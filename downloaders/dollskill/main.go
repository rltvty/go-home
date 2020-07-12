package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/yosssi/gohtml"
	"golang.org/x/net/html"
)

func main() {
	http.HandleFunc("/", clearance)
	http.ListenAndServe(":8080", nil)
}

func getPage(url string) (string, int) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error getting url: %v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Printf("Error parsing html: %v", err)
		os.Exit(1)
	}

	var f func(*html.Node)
	itemHTML := ""
	itemCount := 0
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "div":
				for _, a := range n.Attr {
					if a.Key == "id" && a.Val == "category_grid_container" {
						parseItems(*n, &itemCount)
						buf := new(bytes.Buffer)
						html.Render(buf, n)
						itemHTML = buf.String()
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}

	}
	f(doc)
	return itemHTML, itemCount
}

func clearance(res http.ResponseWriter, req *http.Request) {
	i := 0
	itemHTML := ""
	for {
		i++
		url := fmt.Sprintf("https://www.dollskill.com/clearance.html?c_=3-446&i_=price_desc&p=%v", i)
		fmt.Println(url)
		items, itemCount := getPage(url)
		fmt.Println(itemCount)
		itemHTML += items
		if itemCount == 0 {
			break
		}
	}

	stylesheet := "<link rel=\"stylesheet\" type=\"text/css\" href=\"https://cdn-cf.dollskill.com/skin/frontend/dollskill/dollskill/css/styles.css?v=3301\" media=\"all\" />"
	renderedHTML := fmt.Sprintf("<!DOCTYPE html><html><head>%v</head><body>%v</body></html>", stylesheet, itemHTML)
	formattedHTML := gohtml.Format(renderedHTML)

	fmt.Fprint(res, formattedHTML)
}

//TODO make this a function on the *html.Node type instead
func includesAttr(n *html.Node, key string, val string) bool {
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

func parseItems(n html.Node, itemCount *int) {
	switch n.Type {
	case html.ElementNode:
		switch n.Data {
		case "script":
			n.Parent.RemoveChild(&n)
			return
		case "div":
			if includesAttr(&n, "class", "swatches-list") {
				n.Parent.RemoveChild(&n)
				return
			}
			if includesAttr(&n, "class", "product-image__size-selector") {
				n.Parent.RemoveChild(&n)
				return
			}
		case "li":
			if includesAttr(&n, "class", "product-item") {
				*itemCount++
			}
		case "a":
			if includesAttr(&n, "class", "product-wishlist-icon") {
				n.Parent.RemoveChild(&n)
				return
			}
		case "span":
			if includesAttr(&n, "style", "display:none") {
				n.Parent.RemoveChild(&n)
				return
			}
		}

	case html.TextNode:

	default:
		n.Parent.RemoveChild(&n)
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "select" {
			c.Parent.RemoveChild(c)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "img" {
			if includesAttr(c, "src", "https://cdn-cf.dollskill.com/skin/frontend/dollskill/dollskill/images/category/wph.png") {
				c.Parent.RemoveChild(c)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "div" {
			if includesAttr(c, "class", "category_page_quick_mobile_cart") {
				c.Parent.RemoveChild(c)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "select" {
			if includesAttr(c, "name", "custom_quick_cart") {
				c.Parent.RemoveChild(c)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseItems(*c, itemCount)
	}
}
