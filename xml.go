package main

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type Commit struct {
	XMLName    xml.Name  `xml:"i"`
	ChatServer string    `xml:"chatserver"`
	ChatID     string    `xml:"chatid"`
	Mission    int       `xml:"mission"`
	MaxLimit   int       `xml:"maxlimit"`
	State      int       `xml:"state"`
	RealName   int       `xml:"real_name"`
	Source     string    `xml:"source"`
	Comments   []Comment `xml:"d"`
}

type Comment struct {
	P       string `xml:"p,attr"`
	Content string `xml:",chardata"`
}

type ParsedComment struct {
	Time     float64
	Type     int
	FontSize int
	Color    int
	SendTime int64
	PoolType int
	MidHash  string
	Dmid     int64
	Content  string
}

func parsePAttribute(p string) (float64, int, int, int, int64, int, string, int64) {
	parts := strings.Split(p, ",")
	if len(parts) < 8 {
		return 0, 0, 0, 0, 0, 0, "", 0
	}
	var time float64
	var danmakuType, fontSize, color, poolType int
	var sendTime int64
	var midHash string
	var dmid int64
	_, _ = fmt.Sscanf(parts[0], "%f", &time)
	_, _ = fmt.Sscanf(parts[1], "%d", &danmakuType)
	_, _ = fmt.Sscanf(parts[2], "%d", &fontSize)
	_, _ = fmt.Sscanf(parts[3], "%d", &color)
	_, _ = fmt.Sscanf(parts[4], "%d", &sendTime)
	_, _ = fmt.Sscanf(parts[5], "%d", &poolType)
	_, _ = fmt.Sscanf(parts[6], "%s", &midHash)
	_, _ = fmt.Sscanf(parts[7], "%d", &dmid)
	return time, danmakuType, fontSize, color, sendTime, poolType, midHash, dmid
}
