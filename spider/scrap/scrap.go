package scrap

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/brkss/cybersecurity-pool/spider/download"
	"github.com/brkss/cybersecurity-pool/spider/utils"
)

func ScrapPage(url string) (string, error) {

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", fmt.Errorf("something went wrong")
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}

func ExtractImageLinks(htmlBody, url string) ([]string, error) {
	if htmlBody == "" {
		return nil, errors.New("html body is empty")
	}
	// Regular expression to find img tags and extract the src attribute
	var re = regexp.MustCompile(`(?i)<img[^>]+src="([^"]+)"`)
	matches := re.FindAllStringSubmatch(htmlBody, -1)

	var links []string
	for _, match := range matches {
		if len(match) > 1 {
			if strings.HasPrefix(match[1], "/") {
				src, err := utils.ConvertRelativePath(url, match[1])
				if err != nil {
					continue
				}
				links = append(links, src)
			} else {
				links = append(links, match[1])
			}
		}
	}

	return links, nil
}

// ExtractLinkURLs function extracts href URLs from anchor tags in the provided HTML body
func ExtractLinkURLs(htmlBody, url string) ([]string, error) {
	if htmlBody == "" {
		return nil, errors.New("html body is empty")
	}
	// Regular expression to find anchor tags and extract the href attribute
	var re = regexp.MustCompile(`(?i)<a[^>]+href="([^"]+)"`)
	matches := re.FindAllStringSubmatch(htmlBody, -1)

	var links []string
	for _, match := range matches {
		if len(match) > 1 {
			// Check if the URL is relative and convert it to an absolute URL using the provided base URL
			if strings.HasPrefix(match[1], "/") {
				href, err := utils.ConvertRelativePath(url, match[1])
				if err != nil {
					continue
				}
				links = append(links, href)
			} else {
				links = append(links, match[1])
			}
		}
	}

	return links, nil
}

// ScrapImages : scrap all images with a page
// it take page content as parameter

func ScrapImages(url, path, body string) error {

	srcs, err := ExtractImageLinks(body, url)
	if err != nil {
		return err
	}

	for _, src := range srcs {
		err = download.DownloadImage(src, path)
		if err != nil {
			continue
		}

		fmt.Println(">>> Image Dowloaded successfuly!")
	}

	return nil
}
