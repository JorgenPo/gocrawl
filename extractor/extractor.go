package extractor

import (
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

// forEachNode вызывает функции pre(x) и post(x) для каждого узла х
// в дереве с корнем п. Обе функции необязательны.
// рге вызывается до посещения дочерних узлов, a post - после,
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(n)
	}
}

// ExtractLinks extracts all local and external links from web page
func ExtractLinks(doc *html.Node, request *http.Request) (localUrls []*url.URL, externalUrls []*url.URL) {
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key != "href" {
					continue
				}

				link, err := request.URL.Parse(attr.Val)
				if err != nil {
					continue
				}

				if link.Host == request.Host {
					localUrls = append(localUrls, link)
					continue
				}

				externalUrls = append(externalUrls, link)
			}
		}
	}

	forEachNode(doc, visitNode, nil)

	return
}
