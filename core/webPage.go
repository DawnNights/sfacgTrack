package core

import (
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func atoi(numStr string) int {
	num, _ := strconv.Atoi(numStr)
	return num
}

func (self *NovelInfo) Init(bookId string) {
	self.Id = bookId
	self.Url = "https://book.sfacg.com/Novel/" + bookId

	doc,_ := goquery.NewDocument(self.Url)
	textRow := doc.Find("div.text-row").Find("span")

	Words := textRow.Eq(1).Text()
	self.Words = atoi(Words[9:len(Words)-14])
	self.Hits = textRow.Eq(2).Text()[9:]
	self.Time = textRow.Eq(3).Text()[9:]

	self.Name = doc.Find("h1.title").Find("span.text").Text()
	self.CoverUrl,_ = doc.Find("div.figure").Find("img").Eq(0).Attr("src")
	self.Collection = atoi(doc.Find("#BasicOperation").Find("a").Eq(2).Text()[7:])
	self.Preview = doc.Find("div.chapter-info").Find("p").Text()
	self.Preview = strings.Replace(self.Preview," ","", -1)
	self.Preview = strings.Replace(self.Preview,"\n","", -1)
	self.Preview = strings.Replace(self.Preview,"\r","", -1)
	self.Preview = strings.Replace(self.Preview,"　","", -1)

	cpUrl,_ :=  doc.Find("div.chapter-info").Find("h3").Find("a").Attr("href")
	self.IsVip = strings.Index(cpUrl,"vip") != -1
	self.NewChapter.Init("https://book.sfacg.com" + cpUrl)
}

func (self *NovelChapter) Init(url string) {
	self.Url = url
	doc,_ := goquery.NewDocument(self.Url)
	desc := doc.Find("div.article-desc").Find("span")

	self.Title = doc.Find("h1.article-title").Text()
	self.Writer = desc.Eq(0).Text()[9:]
	self.Time = desc.Eq(1).Text()[15:]
	self.Time = strings.Replace(self.Time,"/","-", -1)
	self.Words = atoi(desc.Eq(2).Text()[9:])

	self.LastUrl,_ = doc.Find("div.fn-btn").Eq(-1).Find("a").Eq(0).Attr("href")
	self.NextUrl,_ = doc.Find("div.fn-btn").Eq(-1).Find("a").Eq(1).Attr("href")
	self.LastUrl = "https://book.sfacg.com" + self.LastUrl
	self.NextUrl = "https://book.sfacg.com" + self.NextUrl
}

func (sf *SFAPI) FindChapterUrl(bookid string) string {
	req,_ := http.Get("http://book.sfacg.com/Novel/" + bookid)
	defer req.Body.Close()
	body,_ := ioutil.ReadAll(req.Body)
	result := string(body)
	result = strmid("<h3 class=\"chapter-title\">","</h3>",result)
	result = strmid("<a href=\"","\" class=",result)
	return "https://book.sfacg.com" + result
}

func (sf *SFAPI) FindBookID(keyword string) string {
	doc,_ := goquery.NewDocument("http://s.sfacg.com/?Key="+keyword+"&S=1&SS=0")
	rsp,_ := doc.Find("#SearchResultList1___ResultList_LinkInfo_0").Attr("href")
	return rsp[28:len(rsp)]
}