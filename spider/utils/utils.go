package utils

import "net/url"

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
