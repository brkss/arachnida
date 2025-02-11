package infrastructure

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/brkss/arachnida/spider/internal/domain"
)


type SpiderServiceImpl struct {
	httpClient *http.Client
}

func NewSpiderService(client *http.Client) *SpiderServiceImpl {
	return &SpiderServiceImpl{
		httpClient: client,
	}
}


// fetchHTML imaplements domain.SpdierService
func (s *SpiderServiceImpl) FetchHTML(url string) (string, error) {
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return "", err;
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err;
	}

	return string(bodyBytes), nil;
	
}

func (s *SpiderServiceImpl) ExtractLinks(baseURL string, content string) ([]domain.ImageLink, []domain.PageLink, error) {
	var images []domain.ImageLink;
	var pages []domain.PageLink;

	// regex to find all image links
	imageRegex := regexp.MustCompile(`(?i)<img[^>]+src=["']([^"']+)["']`)
	matches := imageRegex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		src := match[1]
		if absolute, ok := isAbsoluteURL(baseURL, src); ok {
			images = append(images, domain.ImageLink{URL: absolute})
		}
	}

	// regex to find all page links
	linkRegex := regexp.MustCompile(`(?i)<a[^>]+href=["']([^"']+)["']`)
	matches = linkRegex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		href := match[1]
		if absolute, ok := isAbsoluteURL(baseURL, href); ok {
			pages = append(pages, domain.PageLink{URL: absolute})
		}
	}

	return images, pages, nil
}


func isAbsoluteURL(baseURL string, link string) (string, bool) {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return link, false;
	}

	parsedLink, err := url.Parse(strings.TrimSpace(link))
	if err != nil {
		return link, false
	}
	
	resolved := parsedBaseURL.ResolveReference(parsedLink)
	return resolved.String(), true;
}