package main

import (
	"net/http"
	"strings"
)

func parseShortURL(url string) (vid string) {
	if !strings.HasPrefix(strings.ToLower(url), "https://b23.tv/") {
		return url
	}
	b23, err := http.NewRequest("GET", url, nil)
	b23.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36 Edg/89.0.774.54")
	b23Rsp, err := client.Do(b23)
	if err != nil {
		return url
	}
	if b23Rsp.StatusCode != 302 {
		return url
	}
	return b23Rsp.Header.Get("Location")
}
