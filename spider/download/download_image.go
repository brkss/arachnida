package download

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

const PATH = "./images"

func SaveImage(path string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
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

func ExtractImageName(imageUrl string) (string, error) {
	imageUrlParsed, err := url.Parse(imageUrl)
	if err != nil {
		return "", err
	}

	imagePath := imageUrlParsed.Path

	imageName := path.Base(imagePath)
	return imageName, nil
}

func DownloadImage(URL, src string) error {
	var err error
	imageUrl := src

	if strings.HasPrefix(imageUrl, "/") {
		imageUrl, err = ConvertRelativePath(URL, src)
		if err != nil {
			return err
		}
	}

	imagePath, err := ExtractImageName(imageUrl)
	if err != nil {
		return err
	}

	imagePath = PATH + "/" + imagePath
	err = SaveImage(imagePath, imageUrl)
	if err != nil {
		return err
	}

	return nil
}
