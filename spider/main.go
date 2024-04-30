package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/brkss/cybersecurity-pool/spider/download"
	"github.com/gocolly/colly"
)

//const URL = "https://en.wikipedia.org/wiki/Luis_Walter_Alvarez"

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

	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("img", func(e *colly.HTMLElement) {
		src := e.Attr("src")
		err := download.DownloadImage(URL, src)
		if err != nil {
			fmt.Printf("Error : %v \n\n\n\n", err)
		} else {
			fmt.Println("Image downloaded successfuly !")
		}

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(URL)

}
