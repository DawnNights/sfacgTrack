package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"gopkg.in/yaml.v2"
)

func LoadConfig() (config Config) {
	var sf SFAPI
	yaml.Unmarshal(FileRead("data/Config.yaml"), &config)
	for idx, _ := range config {
		config[idx].RecordUrl = sf.FindChapterUrl(config[idx].BookId)
	}
	return config
}

func (sf *SFAPI) SetCookie(cookie string) {
	FileWrite("data/Cookie.txt", []byte(cookie))
}

func (sf *SFAPI) GetCookie() string {
	return string(FileRead("data/Cookie.txt"))
}

func (nv *Novel) MakeImage() string {
	dc := gg.NewContext(904, 432)

	// 加载背景图片并画入
	im, _ := gg.LoadImage("data/Result/Base.jpg")
	dc.DrawImage(im, 0, 0)

	// 加载封面图片并画入，不存在则下载
	if !FileExist("data/Cover/" + nv.Id + ".png") {
		req, _ := http.Get(nv.CoverUrl)
		defer req.Body.Close()
		cover, _ := ioutil.ReadAll(req.Body)
		im, _, _ = image.Decode(bytes.NewReader(cover))
		im = resize.Resize(240, 335, im, resize.Lanczos3)
		gg.SavePNG("data/Cover/"+nv.Id+".png", im)
	}
	im, _ = gg.LoadImage("data/Cover/" + nv.Id + ".png")
	dc.DrawImage(im, 10, 100)
	dc.SetColor(color.White)

	// 写入小说名称
	dc.LoadFontFace("data/FZYTK.TTF", 26)
	dc.DrawString("《 "+nv.Name+" 》", 0, 27)

	// 写入更新时间和更新章节
	dc.LoadFontFace("data/FZYTK.TTF", 23)
	dc.DrawString(nv.NewChapter.Time.Format("【更新时间: 2006年01月02日 15时04分05秒】"), 0, 61)
	dc.DrawString("【"+nv.NewChapter.Title+"】", 0, 93)

	// 写入小说字数、收藏量和点击量
	dc.LoadFontFace("data/FZYTK.TTF", 20)
	dc.DrawString("字数: "+nv.WordNum, 730, 23)
	dc.DrawString("收藏: "+nv.Collection, 730, 53)
	dc.DrawString("点击: "+nv.HitNum, 730, 83)

	// 写入小说更新内容概述
	if nv.IsVip {
		dc.LoadFontFace("data/FZYTK.TTF", 28)
		for i, text := range SplitText(nv.Preview, 20) {
			dc.DrawString(text, 300, 130+float64(i*40))
		}
	} else {
		dc.LoadFontFace("data/FZYTK.TTF", 24)
		for i, text := range SplitText(nv.Preview, 24) {
			dc.DrawString(text, 300, 130+float64(i*30))
		}
	}

	path, _ := os.Getwd()
	path = strings.ReplaceAll(path, "\\", "/")
	path = path + "/data/Result/" + nv.Id + ".png"

	dc.SavePNG(path)
	return path
}

func (nv *Novel) makeCompare() Compare {
	var cm Compare
	var this, last Chapter
	this = nv.NewChapter
	last.Init(this.LastUrl)

	cm.WordGap = this.WordNum - last.WordNum
	cm.TimeGap = this.Time.Sub(last.Time)

	cm.Times = 1
	timeGap := cm.TimeGap
	for timeGap.Hours() < 10 {
		this = last
		last.Init(this.LastUrl)
		timeGap = this.Time.Sub(last.Time)
		cm.Times++
	}

	var num int
	cm.Stars = 4 + cm.Times

	for _, num = range []int{2500, 3500, 4500, 5500, 6500} {
		if nv.NewChapter.WordNum >= num {
			cm.Stars++
		}
	}

	for _, num = range []int{400, 600, 800} {
		if cm.WordGap >= num {
			cm.Stars++
		}
	}

	for _, num = range []int{-500, -700, -900} {
		if cm.WordGap <= num {
			cm.Stars--
		}
	}

	for _, num = range []int{46800, 28800, 10800} {
		if cm.TimeGap.Seconds() < float64(num) {
			cm.Stars++
		}
	}

	for _, num = range []int{108000, 144000, 180000, 216000, 252000, 288000} {
		if cm.TimeGap.Seconds() > float64(num) {
			cm.Stars--
		}
	}

	return cm
}

func (nv *Novel) MakeJson() (string, string) {
	var cord = JsonCord{}
	var cm = nv.makeCompare()
	_ = json.Unmarshal([]byte(`{"app":"com.tencent.miniapp","desc":"","view":"notification","ver":"0.0.0.1","prompt":"","appID":"","sourceName":"","actionData":"","actionData_A":"","sourceUrl":"","meta":{"notification":{"appInfo":{"appName":"","appType":4,"appid":1109659848,"iconUrl":""},"data":[{"title":"更新字数","value":""},{"title":"间隔时间","value":""},{"title":"更新状态","value":""},{"title":"本章评分","value":""},{"title":"综合评价","value":""}],"title":"本次评价如下：","button":[{"name":"请及时观看最新章节","action":""}],"emphasis_keyword":""}},"text":"","sourceAd":"","extra":""}`), &cord)

	cord.Prompt = fmt.Sprintf("小说《%s》更新啦", nv.Name)
	cord.Meta.Notification.AppInfo.AppName = nv.Name
	cord.Meta.Notification.AppInfo.IconURL = nv.HeadUrl

	wordGap := fmt.Sprintf("%d字", nv.NewChapter.WordNum)
	if cm.WordGap > 0 {
		wordGap = wordGap + fmt.Sprintf("(较上章多更%d字)", cm.WordGap)
	} else if cm.WordGap < 0 {
		wordGap = wordGap + fmt.Sprintf("(较上章少更%d字)", -cm.WordGap)
	} else {
		wordGap = wordGap + "较上章无变化"
	}
	cord.Meta.Notification.Data[0].Value = wordGap

	timeGap := cm.TimeGap.String()
	timeGap = strings.Replace(timeGap, "h", "小时", 1)
	timeGap = strings.Replace(timeGap, "m", "分钟", 1)
	timeGap = strings.Replace(timeGap, "s", "秒", 1)
	cord.Meta.Notification.Data[1].Value = timeGap

	cord.Meta.Notification.Data[2].Value = fmt.Sprintf("今日第『%d』更", cm.Times)
	cord.Meta.Notification.Data[3].Value = MakeStar(cm.Stars)

	if cm.Stars >= 9 {
		cord.Meta.Notification.Data[4].Value = "您就是码字之神在世？"
	} else if cm.Stars >= 6 {
		cord.Meta.Notification.Data[4].Value = fmt.Sprintf("感谢『%s』同学今日的积极更新", nv.Writer)
	} else if cm.Stars >= 3 {
		cord.Meta.Notification.Data[4].Value = fmt.Sprintf("『%s』同学今日的更新不太给力呢", nv.Writer)
	} else {
		cord.Meta.Notification.Data[4].Value = fmt.Sprintf("请『%s』同学尽快爆更补偿吧", nv.Writer)
	}

	body, _ := json.Marshal(&cord)
	return string(body), strings.Join([]string{
		"更新字数: " + wordGap,
		"间隔时间: " + timeGap,
		"更新状态: " + fmt.Sprintf("今日第『%d』更", cm.Times),
		"本章评分: " + MakeStar(cm.Stars),
		"综合评价: " + cord.Meta.Notification.Data[4].Value,
	}, "\n")
}
