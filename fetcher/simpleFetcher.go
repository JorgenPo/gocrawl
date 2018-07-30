package fetcher

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jorgenpo/scrawler/extractor"
	"golang.org/x/net/html"
)

// SimpleFetcher - Not concurent Fetcher Interface implementation
type SimpleFetcher struct {
}

// Fetch - gets page body and urls
func (s SimpleFetcher) Fetch(url string) (result Result) {
	resp, err := http.Get(url)
	if err != nil {
		result.Err = err
		return
	}

	result.Url = resp.Request.URL

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		result.Err = err
		return
	}

	result.Body = string(body[:])

	bodyReader := strings.NewReader(result.Body)

	// Parsing document
	doc, err := html.Parse(bodyReader)
	if err != nil {
		result.Err = err
		return
	}

	if err != nil {
		result.Err = err
		result.Body = string(body[:])
	}

	// Extract links
	localUrls, externalUrls := extractor.ExtractLinks(doc, resp.Request)
	result.LocalUrls = localUrls
	result.ExternalUrls = externalUrls

	return
}
