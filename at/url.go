package at

import (
	"log"
	"net/url"
)

// GetBaseURL extracts the base URL (scheme and host) from the given raw URL string.
// If parsing fails, it logs the error and returns an empty string.
func GetBaseURL(rawURL string) string {
	// Parse the provided raw URL.
	u, err := url.Parse(rawURL)
	if err != nil {
		// Log the error if the URL is invalid.
		log.Printf("could not parse url: %v\n", err)
		return ""
	}

	// Construct the base URL using the scheme and host.
	baseURL := u.Scheme + "://" + u.Host
	return baseURL
}
