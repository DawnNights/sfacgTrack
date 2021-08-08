package bot

import (
	"fmt"
	zero "github.com/wdvxdr1123/ZeroBot"
	"sfacg/core"
	"strings"
)

var sf core.SFAPI

func init()  {
	go sfTrack()

	zero.OnFullMatch("查看登录").Handle(func(ctx *zero.Ctx) {
		ctx.Send(sf.GetCookie())
	})

	zero.OnCommand("更改登录").Handle(func(ctx *zero.Ctx) {
		sf.SetCookie(ctx.State["args"].(string))
		ctx.Send("本地Cookie更新成功")
	})

	zero.OnCommand("查找书号").Handle(func(ctx *zero.Ctx) {
		var info core.NovelInfo
		bookid := sf.FindBookID(ctx.State["args"].(string))
		info.Init(bookid)
		ctx.Send(fmt.Sprintf(
			"书名: %s\n书号: %s\n作者: %s\n更新时间: %s",
			info.Name,info.Id,info.NewChapter.Writer,info.Time))
	})

	zero.OnCommand("测试小说").Handle(func(ctx *zero.Ctx) {
		var info core.NovelInfo
		info.Init(ctx.State["args"].(string))
		xmlText := info.MakeXmlCord()
		jsonText,_ := info.MakeJsonCord()
		ctx.Send(xmlText)
		ctx.Send(jsonText)
	})
}

func sfTrack()  {
	var info core.NovelInfo
	var bot *zero.Ctx
	var groupId int64
	var xmlText,jsonText,cmtText,record string
	config := core.LoadConfig()

	zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
		bot = ctx
		fmt.Println("\n===================================================================")
		fmt.Println("* Version 1.1.0 - 2021-08-08 09:10:29 +0800 CST")
		fmt.Println("* Project: https://github.com/DawnNights/sfacgTrack")
		fmt.Println("* Config: Read",len(config),"novels with local configuration")
		fmt.Println("===================================================================\n")
		return false
	})

	for {
		for idx, _ := range config {
			if sf.FindChapterUrl(config[idx].BookId) != config[idx].RecordUrl{
				info.Init(config[idx].BookId)
				config[idx].RecordUrl = info.NewChapter.Url
				xmlText = info.MakeXmlCord()
				jsonText, cmtText = info.MakeJsonCord()
				record = cmtText[0:strings.Index(cmtText,"本章评分")] + "评论状态: "


				for _, groupId = range config[idx].GroupId{
					bot.SendGroupMessage(groupId,xmlText)
					bot.SendGroupMessage(groupId,jsonText)
				}

				if config[idx].IsSend {
					record = record + sf.Comment(config[idx].BookId,cmtText)
				}else {
					record = record + "禁止评论"
				}

				bot.SendGroupMessage(578562889,record)
			}
		}
	}

}