package bot

import (
	"fmt"
	"sfacg/core"
	"strings"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var api core.SFAPI
var mode int

func init() {
	go sfacgTrack()

	zero.On("notice/group_increase").Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.At(ctx.Event.UserID), message.Text("欢迎新人"))
	})

	zero.OnFullMatch("切换报更模式").Handle(func(ctx *zero.Ctx) {
		switch mode {
		case 0:
			mode = 1
			ctx.Send("已切换为合并转发模式")
		case 1:
			mode = 0
			ctx.Send("已切换为普通报更模式")
		}
	})

	zero.OnFullMatch("查看登录").Handle(func(ctx *zero.Ctx) {
		ctx.Send(api.GetCookie())
	})

	zero.OnCommand("更改登录").Handle(func(ctx *zero.Ctx) {
		api.SetCookie(ctx.State["args"].(string))
	})

	zero.OnCommand("查找书号").Handle(func(ctx *zero.Ctx) {
		var novel core.Novel
		novel.Init(api.FindBookID(ctx.State["args"].(string)))

		var content = "书名: " + novel.Name +
			"\n书号: " + novel.Id +
			"\n作者: " + novel.Writer +
			"\n更新: " + novel.NewChapter.Time.Format("2006年01月02日 15时04分05秒")

		ctx.Send(content)
	})

	zero.OnCommand("测试小说").Handle(func(ctx *zero.Ctx) {
		var novel core.Novel
		novel.Init(ctx.State["args"].(string))

		img := "file:///" + novel.MakeImage()
		code, _ := novel.MakeJson()

		switch mode {
		case 0:
			ctx.SendChain(message.Image(img))
			ctx.SendChain(message.JSON(code))
		case 1:
			code = strings.ReplaceAll(code, ",", "&#44;")
			code = strings.ReplaceAll(code, "[", "&#91;")
			code = strings.ReplaceAll(code, "]", "&#93;")

			ctx.SendGroupForwardMessage(ctx.Event.GroupID, message.Message{
				message.CustomNode(ctx.Event.Sender.NickName, ctx.Event.UserID, "[CQ:image,file="+img+"]"),
				message.CustomNode(ctx.Event.Sender.NickName, ctx.Event.UserID, "[CQ:json,data="+code+"]"),
			})
		}
	})

}

func sfacgTrack() {
	var bot *zero.Ctx
	var novel core.Novel
	var config = core.LoadConfig()
	var content = strings.Join([]string{
		"======================[Sfacg-Track]======================",
		"* OneBot + ZeroBot + Golang",
		fmt.Sprintf("* And there are %d Novels", len(config)),
		"=========================================================",
	}, "\n")

	zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
		bot = ctx
		fmt.Println(content)
		return false
	})

	var img, cmt, code string
	for {
		for idx := 0; idx < len(config); idx++ {
			if config[idx].RecordUrl == api.FindChapterUrl(config[idx].BookId) {
				continue
			}

			novel.Init(config[idx].BookId)
			config[idx].RecordUrl = novel.NewChapter.Url

			img = "file:///" + novel.MakeImage()
			code, cmt = novel.MakeJson()
			code = strings.ReplaceAll(code, ",", "&#44;")
			code = strings.ReplaceAll(code, "[", "&#91;")
			code = strings.ReplaceAll(code, "]", "&#93;")

			switch mode {
			case 0:

				for _, groupID := range config[idx].GroupID {
					bot.SendGroupMessage(groupID, "[CQ:image,file="+img+"]")
					bot.SendGroupMessage(groupID, "[CQ:json,data="+code+"]")
				}

			case 1:

				msg := message.Message{
					message.CustomNode(novel.Writer, config[idx].UserID, "[CQ:image,file="+img+"]"),
					message.CustomNode(novel.Writer, config[idx].UserID, "[CQ:json,data="+code+"]"),
					message.CustomNode(novel.Writer, config[idx].UserID, "积极作者，在线更新~"),
				}

				for _, groupID := range config[idx].GroupID {
					bot.SendGroupForwardMessage(groupID, msg)
				}

			}

			if config[idx].IsSend {
				api.SendComment(novel.Id, cmt)
			}
		}
		time.Sleep(5 * time.Second)
	}

}
