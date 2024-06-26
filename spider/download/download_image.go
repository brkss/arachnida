package download

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

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

func ExtractImageName(imageUrl string) (string, error) {
	imageUrlParsed, err := url.Parse(imageUrl)
	if err != nil {
		return "", err
	}

	imagePath := imageUrlParsed.Path

	imageName := path.Base(imagePath)
	return imageName, nil
}

func DownloadImage(src, PATH string) error {
	var err error
	imageUrl := src

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
