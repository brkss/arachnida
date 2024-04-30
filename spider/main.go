package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/gocolly/colly"
)

func DownloadImage(path string, url string) error  {
	resp, err := http.Get(url)
	if err != nil {
		return err;
	}
	defer resp.Body.Close();

	out, err := os.Create(path);
	if err != nil {
		return err;
	}

	defer out.Close();

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err;
	}

	return nil;
}

func ConvertRelativePath(URL, imageSRC string) (string, error) {

	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	
	parsedImageSrc, err := url.Parse(imageSRC)
	if err != nil {
		return "", err
	}
	fullImageURL := u.ResolveReference(parsedImageSrc).String()
	return fullImageURL, nil
}

func ExtractImageName(imageUrl string)(string, error){
	imageUrlParsed, err := url.Parse(imageUrl);
	if err != nil {
		return "", err;
	}

	imagePath := imageUrlParsed.Path

	imageName := path.Base(imagePath);
	return imageName, nil;
}

func main() {

	// testing : "https://en.wikipedia.org/wiki/Main_Page"
	const URL = "https://en.wikipedia.org/wiki/Luis_Walter_Alvarez"
	const PATH = "./images";

	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("img", func(e *colly.HTMLElement) {
		src := e.Attr("src");
		var err error;
		imageUrl := src;
		if strings.HasPrefix(imageUrl, "/"){
			imageUrl, err = ConvertRelativePath(URL, src);		
			if err != nil {
				fmt.Println("Error : ", imageUrl);
				return;
			}
		}
		imagePath, err := ExtractImageName(imageUrl)
		imagePath = PATH + "/" + imagePath
		err = DownloadImage(imagePath, imageUrl);
		if err != nil {
			log.Fatal("Error downloading image : ", err);
			fmt.Println("Error >> ", imageUrl);
		}else {
			fmt.Println("Image >> ", imageUrl);
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		})

	c.Visit(URL)


}
