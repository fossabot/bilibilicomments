package common

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
	Level    int
	UID      string
	Content  string
}

func ParsePAttribute(p string) (time float64, danmakuType, fontSize, color int, sendTime int64, poolType int, midHash string, dmid int64, level int) {
	parts := strings.SplitN(p, ",", 9)
	if len(parts) < 8 {
		return
	}
	if len(parts) == 8 {
		parts = append(parts, "0")
	}

	_, _ = fmt.Sscanf(parts[0], "%f", &time)
	_, _ = fmt.Sscanf(parts[1], "%d", &danmakuType)
	_, _ = fmt.Sscanf(parts[2], "%d", &fontSize)
	_, _ = fmt.Sscanf(parts[3], "%d", &color)
	_, _ = fmt.Sscanf(parts[4], "%d", &sendTime)
	_, _ = fmt.Sscanf(parts[5], "%d", &poolType)
	midHash = parts[6]
	_, _ = fmt.Sscanf(parts[7], "%d", &dmid)
	_, _ = fmt.Sscanf(parts[8], "%d", &level)
	return time, danmakuType, fontSize, color, sendTime, poolType, midHash, dmid, level
}
