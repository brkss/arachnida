package main

import (
	"flag"
	"log"

	"github.com/brkss/cybersecurity-pool/spider/scrap"
)

func main() {

	recursive := flag.Bool("r", false, "Recursively download images")
	maxDepth := flag.Int("l", 5, "Maximum depth for recursive download")
	savePath := flag.String("p", "./data/", "Path to save downloaded files")

	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("URL is required!")
	}

	URL := args[0]

	visited := make(map[string]bool)
	links := []string{URL} // Starting URL

	currentDepth := 0
	for currentDepth < *maxDepth {
		if !*recursive && currentDepth > 0 {
			break
		}

		newLinks := []string{}

		for _, link := range links {
			if visited[link] {
				continue
			}
			visited[link] = true

			pageContent, err := scrap.ScrapPage(link)
			if err != nil {
				continue
			}

			err = scrap.ScrapImages(link, *savePath, pageContent)
			if err != nil {
				continue
			}

			if currentDepth+1 < *maxDepth {
				extractedLinks, err := scrap.ExtractLinkURLs(pageContent, link)
				if err != nil {
					continue
				}

				for _, newLink := range extractedLinks {
					if !visited[newLink] {
						newLinks = append(newLinks, newLink)
					}
				}
			}
		}

		links = newLinks
		currentDepth++
	}
}
