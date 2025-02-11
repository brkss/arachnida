package domain

// SpiderService define the behavior of our use case needs from html parser and http client
type SpiderService interface {
	// fetch the html of a page
	FetchHTML(url string) (string, error)

	// extract the links from the html and return the image links and page links
	ExtractLinks(baseURL string, content string) ([]ImageLink, []PageLink, error)
}


// FileSaver define the behavior for saving files locally 
type FileSaver interface {
	SaveImage(imageURL, destPath string) error
}





