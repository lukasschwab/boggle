package game

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(url string) (path string, err error) {
	file, err := os.CreateTemp("", "game*.boggle")
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %w", err)
	}
	defer file.Close() //nolint:errcheck

	resp, err := http.Get(url)
	if err != nil {
		return file.Name(), fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if _, err := io.Copy(file, resp.Body); err != nil {
		return file.Name(), fmt.Errorf("error downloading file: %w", err)
	}
	return file.Name(), nil
}
