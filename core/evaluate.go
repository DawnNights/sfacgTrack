package core

import (
	"fmt"
	"strings"
	"time"
)

func (new *NovelChapter) TakeNotes() TrackInfo {
	var note TrackInfo
	var last NovelChapter
	last.Init(new.LastUrl)
	wordGap := new.Words - last.Words
	if wordGap > 0{
		note.WordsGap = fmt.Sprintf("%d字（较上章多更%d字）",new.Words,wordGap)
	}else if wordGap < 0{
		note.WordsGap = fmt.Sprintf("%d字（较上章少更%d字）",new.Words,-wordGap)
	}else {
		note.WordsGap = fmt.Sprintf("%d字（较上章无变化）",new.Words)
	}

	loc, _ := time.LoadLocation("Local")	//获取本地时区
	timeNow, _:= time.ParseInLocation("2006-1-2 15:04:05", new.Time, loc)
	timeLast, _:= time.ParseInLocation("2006-1-2 15:04:05", last.Time, loc)
	timeGap := timeNow.Sub(timeLast)
	note.TimeGap = timeGap.String()
	note.TimeGap = strings.Replace(note.TimeGap,"h","小时",1)
	note.TimeGap = strings.Replace(note.TimeGap,"m","分钟",1)
	note.TimeGap = strings.Replace(note.TimeGap,"s","秒",1)

	var times int = 1
	var now NovelChapter
	for timeGap.Hours() < 10 {
		now = last
		last.Init(now.LastUrl)
		timeNow, _ = time.ParseInLocation("2006-1-2 15:04:05", now.Time, loc)
		timeLast, _ = time.ParseInLocation("2006-1-2 15:04:05", last.Time, loc)
		timeGap = timeNow.Sub(timeLast)
		times++
	}
	note.Times = fmt.Sprintf("今日第『%d』更",times)

	var num int
	var stars int = 5
	for _, num = range []int{2500,3500,4500,5500,6500} {
		if new.Words >= num{
			stars++
		}
	}

	for _, num = range []int{400,600,800} {
		if wordGap >= num{
			stars++
		}
	}
	for _, num = range []int{-500,-700,-900} {
		if wordGap <= num{
			stars--
		}
	}

	for _, num = range []int{46800,28800,10800} {
		if timeGap.Seconds() < float64(num) {
			stars++
		}
	}
	for _, num = range []int{108000,144000,180000,216000,252000,288000} {
		if timeGap.Seconds() > float64(num) {
			stars--
		}
	}

	for _, num = range []int{2,3,4,5,6,7} {
		if times >= num {
			stars++
		}
	}

	if stars > 10{
		stars = 10
	}else if stars < 1 {
		stars = 1
	}

	for i := 1; i < 11; i++ {
		if stars >= i{
			note.Stars = note.Stars + "★"
		}else {
			note.Stars = note.Stars + "☆"
		}
	}

	comments := map[int]string{1: "『%s』同学又是没好好更新的一天呢", 2: "『%s』同学又是没好好更新的一天呢", 3: "『%s』同学又是没好好更新的一天呢", 4: "『%s』同学今天的更新有待提升呢", 5: "请大家愉快阅读，积极催更『%s』同学", 6: "『%s』同学今日更新辛苦了", 7: "感谢『%s』同学今日的积极更新", 8: "今天又是『%s』同学积极更新的一天", 9: "您就是码字之神在世？？？", 10: "您就是码字之神在世？？？"}
	if stars < 9{
		note.Comment = fmt.Sprintf(comments[stars],new.Writer)
	}else {
		note.Comment = comments[stars]
	}
	return note
}

func (self *NovelInfo) MakeXmlCord() string {
	xmlText := fmt.Sprintf(
		`[CQ:xml,data=<?xml version='1.0' encoding='UTF-8' standalone='yes' ?><msg serviceID="83" templateID="1" action="web" brief="《%s》更新啦" sourceMsgId="0" url="https://post.mp.qq.com/" flag="0" adverSign="0" multiMsgFlag="0"><item layout="2" advertiser_id="0" aid="0"><picture cover="%s" w="0" h="0" /><title>%s</title><summary>%s</summary></item><source name="SF互动传媒网" icon="http://p.qpic.cn/qqconnect_app_logo/5mxuSU5RGhZEMaeicmQqdjrXBek1ZicOKpZFDOZWOVvow/0?date=20200607" url="http://url.cn/5w4mpO3" action="app" a_actionData="com.sfacg" i_actionData="tencent10112681://" appid="10112681" /></msg>]`,
		self.Name,self.Dawn(),self.Name,self.NewChapter.Title)
	return xmlText
}

func (self *NovelInfo) MakeJsonCord() (string,string) {
	note := self.NewChapter.TakeNotes()
	paramsA := fmt.Sprintf(
		`{"title":"更新字数"&#44;"value":"%s"}&#44;{"title":"间隔时间"&#44;"value":"%s"}&#44;{"title":"更新状态"&#44;"value":"%s"}&#44;{"title":"本章评分"&#44;"value":"%s"}&#44;{"title":"综合评价"&#44;"value":"%s"}`,
		note.WordsGap,note.TimeGap,note.Times,note.Stars,note.Comment)
	cmtText := fmt.Sprintf(
		"更新字数: %s\n间隔时间: %s\n更新状态: %s\n本章评分: %s\n综合评价: %s",
		note.WordsGap,note.TimeGap,note.Times,note.Stars,note.Comment)
	jsonText := fmt.Sprintf(
		`[CQ:json,data={"app":"com.tencent.miniapp"&#44;"desc":""&#44;"view":"notification"&#44;"ver":"0.0.0.1"&#44;"prompt":"小说《%s》更新啦"&#44;"appID":""&#44;"sourceName":""&#44;"actionData":""&#44;"actionData_A":""&#44;"sourceUrl":""&#44;"meta":{"notification":{"appInfo":{"appName":"%s"&#44;"appType":4&#44;"appid":1109659848&#44;"iconUrl":"http:\/\/img03.sogoucdn.com\/app\/a\/100520146\/05b054de56047cf4ff58ecc2c0cf7a57"}&#44;"data":&#91;%s&#93;&#44;"title":"本次评价如下："&#44;"button":&#91;{"name":"请及时观看最新章节"&#44;"action":"http:\/\/q1.qlogo.cn\/g?b=qq&amp;nk=1543366909&amp;s=640"}&#93;&#44;"emphasis_keyword":""}}&#44;"text":""&#44;"sourceAd":""&#44;"extra":""}]`,
		self.Name,self.Name,paramsA)
	return jsonText,cmtText
}