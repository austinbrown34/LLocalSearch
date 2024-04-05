package utils

import (
	"io"
	"net/http"
)

func IsValidURL(url string) bool {
	// Make an HTTP HEAD request to the URL
	resp, err := http.Head(url)
	if err != nil {
		return false // Return false if the request failed
	}
	defer resp.Body.Close()

	// Check if the status code is in the 2xx range
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

func DownloadWebsiteText(url string) (string, error) {
	// Make an HTTP GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return "", err // Return the error if the request failed
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err // Return the error if reading the body failed
	}

	// Convert the body to a string and return it
	return string(body), nil
}
