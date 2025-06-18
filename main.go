package main

import (
	"bufio"
	"compress/flate"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"bilibilicomments/common"

	"github.com/tidwall/gjson"
)

func getCID(input string) (int64, error) {
	cid, err := strconv.ParseInt(input, 10, 64)
	if err == nil {
		return cid, nil
	}

	vid := common.MatchVIDURL(common.ParseShortURL(input))
	if vid == "" {
		vid = input
	}

	log.Println("开始请求获取CID...")
	cidURL, err := url.Parse("https://api.bilibili.com/x/player/pagelist")
	if err != nil {
		return 0, fmt.Errorf("解析URL失败: %v", err)
	}
	query := cidURL.Query()
	switch strings.ToLower(vid)[:2] {
	case "av":
		query.Add("avid", vid)
	case "bv":
		query.Add("bvid", vid)
	default:
		return 0, fmt.Errorf("不支持的输入格式")
	}
	cidURL.RawQuery = query.Encode()

	resp, err := common.Client.Get(cidURL.String())
	if err != nil {
		return 0, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("读取响应失败: %v", err)
	}

	result := gjson.Parse(string(body))
	data := result.Get("data")
	if len(data.Array()) <= 0 {
		return 0, fmt.Errorf("未找到视频")
	}

sel:
	pid := data.Array()[0].Get("page").Int()
	if len(data.Array()) > 1 {
		data.ForEach(func(key, value gjson.Result) bool {
			log.Printf("%d, Part标题: %s, PartCID: %d",
				value.Get("page").Int(),
				value.Get("part").String(),
				value.Get("cid").Int(),
			)
			return true
		})

		inp, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return 0, fmt.Errorf("读取输入失败: %v", err)
		}
		pid, err = strconv.ParseInt(strings.TrimSpace(inp), 10, 64)
		if err != nil {
			log.Printf("无效输入 `%s`，使用默认 part %d\n", inp, pid)
			goto sel
		}
	}

	var cidVal int64
	result.Get("data").ForEach(func(key, value gjson.Result) bool {
		if value.Get("page").Int() == pid {
			cidVal = value.Get("cid").Int()
			return false
		}
		return true
	})

	if cidVal <= 0 {
		return 0, fmt.Errorf("无效的 CID")
	}
	return cidVal, nil
}

type entry struct {
	time    float64
	color   int
	content string

	midHash string

	uid   string
	ready bool
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalln(os.Args[0], "[av/bv号/url/cid]")
	}

	cid, err := getCID(os.Args[1])
	if err != nil {
		log.Fatalln("获取CID失败:", err)
	}

	log.Println("使用 CID:", cid)

	u, _ := url.Parse("https://api.bilibili.com/x/v1/dm/list.so")
	q := u.Query()
	q.Add("oid", strconv.FormatInt(cid, 10))
	u.RawQuery = q.Encode()

	resp, err := common.Client.Get(u.String())
	if err != nil {
		log.Fatalln("请求弹幕失败:", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalln("请求返回非 200:", resp.Status)
	}

	var reader io.ReadCloser
	if resp.Header.Get("Content-Encoding") == "deflate" {
		reader = flate.NewReader(resp.Body)
		defer reader.Close()
	} else {
		reader = resp.Body
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalln("读弹幕失败:", err)
	}

	var commit common.Commit
	if err := xml.Unmarshal(data, &commit); err != nil {
		log.Fatalln("XML解析失败:", err)
	}
	log.Printf("共 %d 条弹幕\n", len(commit.Comments))

	entries := make([]*entry, len(commit.Comments))
	for i, cm := range commit.Comments {
		t, _, _, color, _, _, midHash, _, _ := common.ParsePAttribute(cm.P)
		entries[i] = &entry{
			time:    t,
			color:   color,
			content: cm.Content,
			midHash: midHash,
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].time < entries[j].time
	})

	var (
		mu   sync.Mutex
		cond = sync.NewCond(&mu)
	)

	wg := sync.WaitGroup{}
	jobs := make(chan *entry)

	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for e := range jobs {
				// crack很慢
				uid := common.Crack(e.midHash)
				mu.Lock()
				e.uid = uid
				e.ready = true
				cond.Broadcast()
				mu.Unlock()
			}
		}()
	}

	go func() {
		for _, e := range entries {
			jobs <- e
		}
		close(jobs)
	}()

	mu.Lock()
	for _, e := range entries {
		for !e.ready {
			cond.Wait()
		}
		fmt.Printf("%.2fs | %s%s%s | UID: %s\n",
			e.time,
			common.RgbToAnsi(e.color),
			e.content,
			common.ResetColor(),
			e.uid)
	}
	mu.Unlock()

	wg.Wait()
}
