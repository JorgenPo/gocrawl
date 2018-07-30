package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/jorgenpo/scrawler/crawler"
)

func printUsage() {
	fmt.Printf("Simple web crawler v0.1\n\n")
	fmt.Printf("Usage: %s -h host -d depth\n", os.Args[0])
	flag.PrintDefaults()
}

// Arguments of the scrawler
type Arguments struct {
	host  string
	depth int
}

func parseArguments() (args Arguments, err error) {
	var host = flag.String("h", "", "host on wich should we start crawl")
	var depth = flag.Int("d", 1, "depth of the crawl")

	flag.Parse()

	if *host == "" {
		return args, fmt.Errorf("Host undefined")
	}

	args.host = *host
	args.depth = *depth

	return
}

func main() {
	args, err := parseArguments()

	if err != nil {
		printUsage()
		os.Exit(1)
	}

	testUrl, err := url.Parse(args.host)
	if err != nil {
		fmt.Printf("%s is not a valid url! Error: %s", os.Args[1], err)
		os.Exit(2)
	}

	depth := args.depth

	fmt.Printf("Web Scrawler v0.1\n\n")
	fmt.Printf("Starting scrawl from %s\n", testUrl)
	fmt.Printf("Depth = %d\n", depth)

	crawler := crawler.NewCrawler(testUrl, depth, &crawler.GEXFSaver{})

	crawler.Crawl()

	before := time.Now()
	for !crawler.Finished {
		time.Sleep(100)
	}
	after := time.Now()

	time := after.Sub(before)

	fmt.Printf("Scrawling ended in %f seconds. Links fetched: %d\n", time.Seconds(), crawler.FetchedUrls)

	err = crawler.SaveResults("./")

	if err != nil {
		fmt.Printf("Failed to save page graph. Error: %s\n", err)
	}

	fmt.Printf("Page graph was successfuly saved\n")
}
