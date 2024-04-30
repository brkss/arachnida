package main

import (
	"fmt"

	"github.com/brkss/cybersecurity-pool/spider/download"
	"github.com/gocolly/colly"
)

const URL = "https://en.wikipedia.org/wiki/Luis_Walter_Alvarez"

func main() {

	// testing : "https://en.wikipedia.org/wiki/Main_Page"

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
