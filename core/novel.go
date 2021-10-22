package core

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func (nv *Novel) Init(bookId string) {
	nv.Id = bookId
	nv.Url = "https://book.sfacg.com/Novel/" + bookId
	doc, _ := goquery.NewDocument(nv.Url)

	nv.Name = doc.Find("h1.title").Find("span.text").Text()
	nv.Writer = doc.Find("div.author-name").Find("span").Text()
	nv.HeadUrl, _ = doc.Find("div.author-mask").Find("img").Attr("src")

	textRow := doc.Find("div.text-row").Find("span")
	nv.HitNum = textRow.Eq(2).Text()[9:]
	nv.WordNum = textRow.Eq(1).Text()
	nv.WordNum = nv.WordNum[9 : len(nv.WordNum)-14]

	nv.CoverUrl, _ = doc.Find("div.figure").Find("img").Eq(0).Attr("src")
	nv.Collection = doc.Find("#BasicOperation").Find("a").Eq(2).Text()[7:]

	nv.Preview = doc.Find("div.chapter-info").Find("p").Text()
	nv.Preview = strings.Replace(nv.Preview, " ", "", -1)
	nv.Preview = strings.Replace(nv.Preview, "\n", "", -1)
	nv.Preview = strings.Replace(nv.Preview, "\r", "", -1)
	nv.Preview = strings.Replace(nv.Preview, "　", "", -1)

	nvUrl, _ := doc.Find("div.chapter-info").Find("h3").Find("a").Attr("href")
	nv.IsVip = strings.Contains(nvUrl, "vip")
	nv.NewChapter.Init("https://book.sfacg.com" + nvUrl)
}

func (cp *Chapter) Init(url string) {
	cp.Url = url
	loc, _ := time.LoadLocation("Local")
	doc, _ := goquery.NewDocument(cp.Url)

	desc := doc.Find("div.article-desc").Find("span")
	cp.WordNum = Atoi(desc.Eq(2).Text()[9:])
	cp.Time, _ = time.ParseInLocation("2006/1/2 15:04:05", desc.Eq(1).Text()[15:], loc)
	cp.Title = doc.Find("h1.article-title").Text()

	cp.LastUrl, _ = doc.Find("div.fn-btn").Eq(-1).Find("a").Eq(0).Attr("href")
	cp.NextUrl, _ = doc.Find("div.fn-btn").Eq(-1).Find("a").Eq(1).Attr("href")
	cp.LastUrl = "https://book.sfacg.com" + cp.LastUrl
	cp.NextUrl = "https://book.sfacg.com" + cp.NextUrl
}

func (sf *SFAPI) FindBookID(keyword string) string {
	doc, _ := goquery.NewDocument("http://s.sfacg.com/?Key=" + keyword + "&S=1&SS=0")
	href, _ := doc.Find("#SearchResultList1___ResultList_LinkInfo_0").Attr("href")

	return href[28:]
}

func (sf *SFAPI) FindChapterUrl(bookid string) string {
	req, _ := http.Get("http://book.sfacg.com/Novel/" + bookid)
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)

	result := string(body)
	result = GetMidText("<h3 class=\"chapter-title\">", "</h3>", result)
	result = GetMidText("<a href=\"", "\" class=", result)

	return "https://book.sfacg.com" + result
}

func (sf *SFAPI) SendComment(bookid string, content string) string {
	client := &http.Client{}
	params := fmt.Sprintf("nid=%s&commentid=0&content=%s", bookid, url.QueryEscape(content))

	req, _ := http.NewRequest("POST", "https://book.sfacg.com/ajax/ashx/Common.ashx?op=addCmt", strings.NewReader(params))
	req.Header = http.Header{
		"User-Agent":   []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:84.0) Gecko/20100101 Firefox/84.0"},
		"Content-Type": []string{"application/x-www-form-urlencoded"},
		"Cookie":       []string{sf.GetCookie()},
	}

	rsp, _ := client.Do(req)
	defer rsp.Body.Close()
	body, _ := ioutil.ReadAll(rsp.Body)

	return string(body)
}
