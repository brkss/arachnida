package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/brkss/cybersecurity-pool/spider/collector"
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

	c := collector.NewSpiderCollector(*recursive, *maxDepth, *savePath)
	c.Visit(URL)

	/*
		// Find and visit all links
		c.OnHTML("img", func(e *colly.HTMLElement) {
			src := e.Attr("src")
			err := download.DownloadImage(URL, src, *savePath)
			if err != nil {
				fmt.Printf("Error : %v \n\n\n\n", err)
			} else {
				fmt.Println("Image downloaded successfuly !")
			}
		})

		c.OnHTML("a", func(e *colly.HTMLElement) {
			link := e.Attr("href")
			fmt.Println("link : ", link)
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL)
		})

	*/

}
