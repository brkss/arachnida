package main

import (
	"flag"
	"fmt"
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

	fmt.Println("-r : ", *recursive)
	fmt.Println("-l : ", *maxDepth)
	fmt.Println("-p : ", *savePath)
	fmt.Println("URL : ", URL)

	fmt.Println("--------------------")

	visited := make(map[string]bool)
	links := []string{URL} // Starting URL

	currentDepth := 0
	for currentDepth < *maxDepth {
		if !*recursive && currentDepth > 0 {
			break
		}

		// New batch of links to process in the next depth level
		newLinks := []string{}

		for _, link := range links {
			if visited[link] {
				continue
			}
			visited[link] = true

			// Fetch page content
			pageContent, err := scrap.ScrapPage(link)
			if err != nil {
				continue
			}

			// Scrap images from the page
			err = scrap.ScrapImages(link, *savePath, pageContent)
			if err != nil {
				continue
			}

			// Extract new links only if within depth limit
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
