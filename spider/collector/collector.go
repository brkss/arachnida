package collector

import (
	"fmt"

	"github.com/brkss/cybersecurity-pool/spider/download"
	"github.com/gocolly/colly/v2"
)

func NewSpiderCollector(recursive bool, depth int, path string) *colly.Collector {
	c := colly.NewCollector(
		colly.MaxDepth(depth),
	)

	if recursive {
		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			link := e.Request.AbsoluteURL(e.Attr("href"))
			e.Request.Visit(link)
		})
	}

	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		src := e.Request.AbsoluteURL(e.Attr("src"))
		err := download.DownloadImage(src, path)
		if err != nil {
			fmt.Println("Error!: cannot Download Image !");
		}else {
			fmt.Println(">> Image Downloaded successfuly !")
		}
	})

	return c

}
