package fetcher

import "net/url"

// Fetcher interface
type Fetcher interface {
	// Fetch returns a body of URL and
	// a slice of URLs found on that page
	Fetch(url string) (body string, urls []string, err error)
}

// Result of fetching web page
type Result struct {
	Body         string
	Url          *url.URL
	LocalUrls    []*url.URL
	ExternalUrls []*url.URL
	Err          error
}
