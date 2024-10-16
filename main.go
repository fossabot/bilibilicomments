package main

import (
	"bufio"
	"compress/flate"
	"encoding/xml"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	log.Println("开始请求获取CID...")
	client := &http.Client{}
	cidURL := "https://api.bilibili.com/x/player/pagelist?"
	args := os.Args
	if strings.HasPrefix(strings.ToLower(args[1]), "av") {
		cidURL += fmt.Sprintf("avid=%s", url.QueryEscape(args[1]))
	} else if strings.HasPrefix(strings.ToLower(args[1]), "bv") {
		cidURL += fmt.Sprintf("bvid=%s", url.QueryEscape(args[1]))
	}
	cidReq, err := http.NewRequest("GET", cidURL, nil)
	if err != nil {
		panic(err)
		return
	}
	cidRsp, err := client.Do(cidReq)
	if err != nil {
		panic(err)
		return
	}
	cidContext, err := io.ReadAll(cidRsp.Body)
	if err != nil {
		panic(err)
		return
	}
	cidJson := gjson.Parse(string(cidContext))
	cidJson.Get("data").ForEach(func(key, value gjson.Result) bool {
		log.Printf("%d, Part标题: %s, PartCID: %d", value.Get("page").Int(), value.Get("part").String(), value.Get("cidReq").Int())
		return true
	})
	inp, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		panic(err)
		return
	}
	inpInt, err := strconv.ParseInt(strings.TrimSuffix(inp, "\r\n"), 10, 64)
	if err != nil {
		panic(err)
		return
	}
	var cid int64
	cidJson.Get("data").ForEach(func(key, value gjson.Result) bool {
		if value.Get("page").Int() == inpInt {
			cid = value.Get("cid").Int()
			return false
		}
		return true
	})
	if cid == 0 {
		log.Println("没有找到对应的CID")
		return
	}
	commitReq, err := http.NewRequest("GET", fmt.Sprintf("https://comment.bilibili.com/%d.xml", cid), nil)
	if err != nil {
		panic(err)
		return
	}
	res, err := client.Do(commitReq)
	if err != nil {
		panic(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	flateReader := flate.NewReader(res.Body)
	defer func(flateReader io.ReadCloser) {
		err := flateReader.Close()
		if err != nil {
			panic(err)
		}
	}(flateReader)
	context, err := io.ReadAll(flateReader)
	if err != nil {
		panic(err)
		return
	}

	var commit Commit
	err = xml.Unmarshal(context, &commit)
	if err != nil {
		panic(err)
		return
	}
	log.Printf("获取到 %d 条弹幕, 开始输出:\n", len(commit.Comments))

	var parsedComments []ParsedComment
	for _, comment := range commit.Comments {
		time, danmakuType, fontSize, color, sendTime, poolType, midHash, dmid := parsePAttribute(comment.P)
		parsedComments = append(parsedComments, ParsedComment{
			Time:     time,
			Type:     danmakuType,
			FontSize: fontSize,
			Color:    color,
			SendTime: sendTime,
			PoolType: poolType,
			MidHash:  midHash,
			Dmid:     dmid,
			Content:  comment.Content,
		})
	}

	sort.Slice(parsedComments, func(i, j int) bool {
		return parsedComments[i].Time < parsedComments[j].Time
	})
	for _, comment := range parsedComments {
		colorCode := rgbToAnsi(comment.Color)
		// debug
		//fmt.Printf("Time: %.2f, Type: %d, Font Size: %d, Color: %d, Send Time: %s, Pool Type: %d, MidHash: %s, Dmid: %d\n", comment.Time, comment.Type, comment.FontSize, comment.Color, time.Unix(comment.SendTime, 0).String(), comment.PoolType, crack(comment.MidHash), comment.Dmid)
		fmt.Printf("%s%s%s UID: %s\n", colorCode, comment.Content, resetColor(), crack(comment.MidHash))
	}
}
