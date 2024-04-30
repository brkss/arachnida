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

	pageContent, err := scrap.ScrapPage(URL)
	if err != nil {
		log.Fatal("something went wrng : ", err)
	}

	err = scrap.ScrapImages(URL, *savePath, pageContent)
	if err != nil {
		log.Fatal("Cannot scrap images !")
	}

	/*
		srcs, err := scrap.ExtractImageLinks(pageContent, URL)

	*/

	// links, err := scrap.ExtractLinkURLs(pageContent, URL)
	// if err != nil {
	// 	log.Fatal("something went wrong extracting image links")
	// }

	// for _, link := range links {
	// 	fmt.Println("link : ", link)
	// }

	/*
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
