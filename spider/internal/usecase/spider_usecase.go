package usecase

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/brkss/arachnida/spider/internal/domain"
)



type SpiderUsecase struct {
	service domain.SpiderService
	fileSaver domain.FileSaver
	visitedURLs map[string]bool // keep track of visited urls to avoid loops 
}


// NewSpiderUsecase creates a new spider usecase
func NewSpiderUsecase(service domain.SpiderService, fileSaver domain.FileSaver) *SpiderUsecase {
	return &SpiderUsecase{
		service: service,
		fileSaver: fileSaver,
		visitedURLs: make(map[string]bool),
	}
}

func (u *SpiderUsecase) DownloadImages(url string, maxDepth int, currentDepth int, savePath string) error {


	if currentDepth > maxDepth {
		return nil;
	}
	if(u.visitedURLs[url]) {
		return nil;
	}

	u.visitedURLs[url] = true;

	html, err := u.service.FetchHTML(url)
	if err != nil {
		return err;
	}

	imageLinks, pageLinks, err := u.service.ExtractLinks(url, html)
	if err != nil {
		return err;
	}

	for _, imageLink := range imageLinks {
		path := strings.Split(imageLink.URL, "/")
		fileName := path[len(path) - 1]
		if fileName == "" {
			fileName = "image.jpg" // set up a default name if the file name is not provided
		}
		dest := filepath.Join(savePath, fileName)
		err := u.fileSaver.SaveImage(imageLink.URL, dest)
		if err != nil {
			fmt.Printf("Failed to save image %s: %v\n", imageLink.URL, err)
		}else {
			fmt.Printf("Saved image %s to %s\n", imageLink.URL, dest)
		}
	}

	maxDepth = currentDepth + 1;

	for _, pageLink := range pageLinks {
		err := u.DownloadImages(pageLink.URL, maxDepth, currentDepth + 1, savePath)
		if err != nil {
			fmt.Printf("Failed to download images from %s: %v\n", pageLink.URL, err)
		}
	}

	return nil;

}


func ValidateDepth(depth int) (int, error){
	if depth < 0 {
		return 0, errors.New("depth must be greater than 0")
	}

	if depth == 0 {
		return 5, nil; // set a default depth of 5 if no depth is provided
	}

	return depth, nil;
}
