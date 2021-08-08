package bot

import (
	"fmt"
	zero "github.com/wdvxdr1123/ZeroBot"
	"sfacg/core"
	"strings"
)

var sf core.SFAPI

func init() {
	zero.RegisterPlugin(&sfacg{})  //注册插件
	go sfTrack() //启动监控线程
}

type sfacg struct{}

func (_ *sfacg) GetPluginInfo() zero.PluginInfo {
        return zero.PluginInfo{
                Author:     "夜不语",
                PluginName: "SF报更",
                Version:    "1.1.0",
                Details:    "SF报更插件，Go语言版",
        }
}

func (_ *sfacg) Start() { // 插件主体
        zero.OnNotice().Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			if event.NoticeType == "group_increase"{
                zero.Send(event, fmt.Sprintf("[CQ:at,qq=%d]欢迎新人",event.UserID))
        }
        return zero.FinishResponse
		})  //欢迎新人
		
        zero.OnCommand("测试小说").Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			var info core.NovelInfo
			info.Init(state["args"].(string))
			xmlText := info.MakeXmlCord()
			jsonText,_ := info.MakeJsonCord()
			zero.Send(event,xmlText)
			zero.Send(event,jsonText)
			return zero.FinishResponse
		})
        
        zero.OnCommand("查找书号").Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			var info core.NovelInfo
			bookid := sf.FindBookID(state["args"].(string))
			info.Init(bookid)
			zero.Send(event,fmt.Sprintf(
			"书名: %s\n书号: %s\n作者: %s\n更新时间: %s",
			info.Name,info.Id,info.NewChapter.Writer,info.Time))
			return zero.FinishResponse
		})
        
        zero.OnFullMatch("查看登录").Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			zero.Send(event,sf.GetCookie())
			return zero.FinishResponse
		})

        zero.OnCommand("更改登录").Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			sf.SetCookie(state["args"].(string))
			zero.Send(event,"本地Cookie更新成功")
        	return zero.FinishResponse
		})
}

func sfTrack()  {
	var info core.NovelInfo
	var groupId int64
	var xmlText,jsonText,cmtText,record string
	config := core.LoadConfig()

	fmt.Println("\n===================================================================")
	fmt.Println("* Version 1.1.0 - 2021-08-08 09:10:29 +0800 CST")
	fmt.Println("* Project: https://github.com/DawnNights/sfacgTrack")
	fmt.Println("* Config: Read",len(config),"novels with local configuration")
	config.Print()
	fmt.Println("===================================================================\n")

	for {
		for idx, _ := range config {
			if sf.FindChapterUrl(config[idx].BookId) != config[idx].RecordUrl{
				info.Init(config[idx].BookId)
				config[idx].RecordUrl = info.NewChapter.Url
				xmlText = info.MakeXmlCord()
				jsonText, cmtText = info.MakeJsonCord()

				for _, groupId = range config[idx].GroupId{
					zero.SendGroupMessage(groupId,xmlText)
					zero.SendGroupMessage(groupId,jsonText)
				}

				if config[idx].IsSend {
					record = record + sf.Comment(config[idx].BookId,cmtText)
				}else {
					record = record + "禁止评论"
				}

				record = fmt.Sprintf(
					"小说书名: %s\n%s最新章节: %s\n评论状态: ",
					info.Name,cmtText[0:strings.Index(cmtText,"本章评分")],info.NewChapter.Title)
				zero.SendGroupMessage(522245324,record)
			}
		}
	}

}