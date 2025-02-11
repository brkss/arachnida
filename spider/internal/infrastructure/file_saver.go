package infrastructure

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type FileSaverImpl struct {
	httpClient *http.Client
}

// NewFileSaver creates a new FileSaverImpl
func NewFileSaver(client *http.Client) *FileSaverImpl {
	return &FileSaverImpl{
		httpClient: client,
	}
}

func (fs *FileSaverImpl) SaveImage(imageURL string, destPath string) error {
	// make sure directory exists
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to download image from %s: %w", imageURL, err)
	}

	// download image
	resp, err := fs.httpClient.Get(imageURL)
	if err != nil {
		return fmt.Errorf("failed to download image from %s: %w", imageURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download image from %s: %s", imageURL, resp.Status)
	}

	// create file
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", destPath, err)
	}
	defer out.Close()

	// write image to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write image to file %s: %w", destPath, err)
	}

	return nil
}



