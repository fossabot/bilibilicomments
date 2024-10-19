package main

import (
	"log"
	"regexp"
	"strings"
)

func matchVIDURL(url string) (vid string) {
	// Define the regular expression pattern to match the Bilibili video ID
	// Pattern explanation:
	// - (?i): case insensitive flag
	// - ^(av|bv): starts with "av" or "bv", case insensitive
	// - [0-9A-Za-z]+: followed by one or more alphanumeric characters
	re := regexp.MustCompile(`(?i)^(av|bv)[0-9A-Za-z]+`)

	// Extract the path from the URL
	splitURL := strings.Split(url, "/")
	for _, part := range splitURL {
		vidMatch := re.FindString(part)
		if vidMatch != "" {
			return vidMatch
		}
	}

	log.Println("匹配不到输入的视频")
	return
}
