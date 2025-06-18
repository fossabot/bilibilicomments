package common

import (
	"net/http"
	"strings"
)

func ParseShortURL(url string) string {
	if !strings.HasPrefix(strings.ToLower(url), "https://b23.tv/") && strings.HasPrefix(url, "https://b23.tv/ep") {
		return url
	}
	resp, err := http.Get(url)
	if err != nil {
		return url
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusFound {
		return url
	}
	return resp.Header.Get("Location")
}
