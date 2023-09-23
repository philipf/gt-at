package at

import (
	"log"
	"net/url"
)

func GetBaseURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		log.Printf("could not parse url: %v\n", err)
		return ""
	}

	baseURL := u.Scheme + "://" + u.Host
	return baseURL
}
