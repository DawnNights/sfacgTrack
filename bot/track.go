package bot

import (
	"fmt"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"sfacg/core"
	"strings"
)

var sf core.SFAPI
var IsNormal bool = true

func init() {

	zero.OnFullMatch("切换报更模式").Handle(func(ctx *zero.Ctx) {
		IsNormal = !IsNormal
		if IsNormal{
			ctx.Send("已切换为普通卡片模式")
		}else {
			ctx.Send("已切换为合并转发模式")
		}
	})

	zero.On("notice/group_increase").Handle(func(ctx *zero.Ctx) {
		ctx.Send(fmt.Sprintf("[CQ:at,qq=%d]欢迎新人",ctx.Event.UserID))
	})

	zero.OnCommand("测试小说").Handle(func(ctx *zero.Ctx) {
		var info core.NovelInfo
		info.Init(ctx.State["args"].(string))

		xmlText := info.MakeXmlCord()
		jsonText,_ := info.MakeJsonCord()

		if IsNormal{
			ctx.Send(xmlText)
			ctx.Send(jsonText)
		}else {
			var msg message.Message = []message.MessageSegment{
				message.CustomNode("测试人员", ctx.Event.SelfID, xmlText),
				message.CustomNode("测试人员", ctx.Event.SelfID, jsonText),
			}
			ctx.SendGroupForwardMessage(ctx.Event.GroupID,msg)

		}
	})

	zero.OnCommand("查找书号").Handle(func(ctx *zero.Ctx) {
		var info core.NovelInfo
		info.Init(sf.FindBookID(ctx.State["args"].(string)))

		ctx.Send(fmt.Sprintf(
			"书名: %s\n书号: %s\n作者: %s\n更新时间: %s",
			info.Name,info.Id,info.NewChapter.Writer,info.Time))
	})

	zero.OnFullMatch("查看登录").Handle(func(ctx *zero.Ctx) {
		ctx.Send(sf.GetCookie())
	})

	zero.OnCommand("更改登录").Handle(func(ctx *zero.Ctx) {
		sf.SetCookie(ctx.State["args"].(string))
		ctx.Send("本地Cookie更新成功")
	})

	go sfTrack()
}

func sfTrack()  {
	var info core.NovelInfo
	var ctx *zero.Ctx
	var groupId int64
	var xmlText,jsonText,cmtText,record string
	config := core.LoadConfig()

	zero.RangeBot(func(id int64, bot *zero.Ctx) bool {
		ctx = bot
		fmt.Println("\n===================================================================")
		fmt.Println("* Version 1.1.0 - 2021-08-08 09:10:29 +0800 CST")
		fmt.Println("* Project: https://github.com/DawnNights/sfacgTrack")
		fmt.Println("* Config: Read",len(config),"novels with local configuration")
		fmt.Println("===================================================================\n")
		return false
	})

	for {for idx, _ := range config {
		if sf.FindChapterUrl(config[idx].BookId) != config[idx].RecordUrl{
			info.Init(config[idx].BookId)
			config[idx].RecordUrl = info.NewChapter.Url
			xmlText = info.MakeXmlCord()
			jsonText, cmtText = info.MakeJsonCord()
			record = fmt.Sprintf(
				"小说书名: %s\n%s最新章节: %s\n评论状态: ",
				info.Name,cmtText[0:strings.Index(cmtText,"本章评分")],info.NewChapter.Title)


			for _, groupId = range config[idx].GroupId{
				if IsNormal {
					ctx.SendGroupMessage(groupId,xmlText)
					ctx.SendGroupMessage(groupId,jsonText)
				}else {
					var msg message.Message = []message.MessageSegment{
						message.CustomNode(info.NewChapter.Writer, config[idx].WriterId, xmlText),
						message.CustomNode(info.NewChapter.Writer, config[idx].WriterId, jsonText),
					}
					ctx.SendGroupForwardMessage(groupId, msg)
				}
			}

			if config[idx].IsSend {
				record = record + sf.Comment(config[idx].BookId,cmtText)
			}else {
				record = record + "禁止评论"
			}

			ctx.SendGroupMessage(522245324,record)
		}
	}}
}
