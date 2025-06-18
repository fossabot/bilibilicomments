package common

import (
	"log"
	"regexp"
	"strings"
)

func MatchVIDURL(url string) (vid string) {
	// (?i): case-insensitive
	// ^(av|bv|ep): starts with av, bv, or ep
	// [0-9A-Za-z]+: followed by numbers or letters
	re := regexp.MustCompile(`(?i)^(av|bv|ep)[0-9A-Za-z]+`)

	for _, part := range strings.Split(url, "/") {
		if vidMatch := re.FindString(part); vidMatch != "" {
			return vidMatch
		}
	}

	log.Println("匹配不到输入的视频")
	return
}
