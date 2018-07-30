package crawler

import (
	"log"
	"net/url"

	"github.com/jorgenpo/scrawler/fetcher"
)

// PageNode - page graph node
type PageNode struct {
	url *url.URL
}

// PageEdge - page graph edge
type PageEdge struct {
	source      *PageNode
	destination *PageNode
}

// Crawler - crawler result object
type Crawler struct {
	StartUrl    *url.URL
	Depth       int
	Nodes       []*PageNode
	Edges       []*PageEdge
	Root        *PageNode
	Fetcher     fetcher.SimpleFetcher
	queuedUrls  int
	FetchedUrls int
	Finished    bool
	Saver       Saver
}

// NewCrawler - creates new crawler object
func NewCrawler(url *url.URL, depth int, saver Saver) (crawler Crawler) {
	return Crawler{StartUrl: url, Depth: depth, Saver: saver}
}

// SaveResults - saves page graph with choosed format
func (c *Crawler) SaveResults(path string) error {
	var filename = path + c.StartUrl.Hostname() + c.Saver.GetExtension()

	log.Printf("[Crawler] Saving graph to %s...\n", filename)
	return c.Saver.Save(filename, c.Nodes, c.Edges)
}

func (c *Crawler) fetchHost(node *PageNode, depth int) {
	log.Printf("[Crawler]: crawling host %s\n", node.url.String())

	result := c.Fetcher.Fetch(node.url.String())

	log.Printf("[Crawler]: found %d external links on %s\n", len(result.ExternalUrls), node.url)

	if depth > 0 {
		c.queuedUrls += len(result.ExternalUrls)

		for _, link := range result.ExternalUrls {
			destNode := &PageNode{link}
			edge := &PageEdge{source: node, destination: destNode}
			c.Nodes = append(c.Nodes, destNode)
			c.Edges = append(c.Edges, edge)

			go c.fetchHost(destNode, depth-1)
		}
	}

	c.FetchedUrls++

	if c.queuedUrls == c.FetchedUrls {
		c.Finished = true
	}
}

// Crawl - start scrawl web from startUrl and to defined depth
func (c *Crawler) Crawl() {
	log.Printf("[Crawler]: start scrawl from %s\n", c.StartUrl)

	c.Root = &PageNode{c.StartUrl}
	c.Nodes = append(c.Nodes, c.Root)

	c.queuedUrls = 1
	c.FetchedUrls = 0

	go c.fetchHost(c.Root, c.Depth)
}
