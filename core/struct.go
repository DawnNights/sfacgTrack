package core

import "time"

type (
	SFAPI uint8

	Novel struct {
		Id         string  // 小说书号
		Url        string  // 小说网址
		Name       string  // 小说书名
		IsVip      bool    // 是否上架
		Writer     string  // 作者昵称
		HitNum     string  // 小说点击
		WordNum    string  // 小说字数
		Preview    string  // 章节预览
		HeadUrl    string  // 头像网址
		CoverUrl   string  // 封面网址
		Collection string  // 小说收藏
		NewChapter Chapter // 章节信息
	}

	Chapter struct {
		Url     string    // 章节网址
		Time    time.Time // 更新时间
		Title   string    // 章节名称
		WordNum int       // 章节字数
		LastUrl string    // 上章网址
		NextUrl string    // 下章网址
	}

	Compare struct {
		Times   int           // 更新次数
		Stars   int           // 更新评星
		WordGap int           // 更新字数差
		TimeGap time.Duration // 更新时间差
		Comment string        // 更新评语
	}

	Config []struct {
		BookId    string  // 报更书号
		IsSend    bool    // 是否评论
		UserID    int64   // 作者帐号
		GroupID   []int64 // 书友群号
		RecordUrl string  // 更新记录
	}

	JsonCord struct {
		App         string `json:"app"`
		Desc        string `json:"desc"`
		View        string `json:"view"`
		Ver         string `json:"ver"`
		Prompt      string `json:"prompt"`
		AppID       string `json:"appID"`
		SourceName  string `json:"sourceName"`
		ActionData  string `json:"actionData"`
		ActionDataA string `json:"actionData_A"`
		SourceURL   string `json:"sourceUrl"`
		Meta        Meta   `json:"meta"`
		Text        string `json:"text"`
		SourceAd    string `json:"sourceAd"`
		Extra       string `json:"extra"`
	}

	AppInfo struct {
		AppName string `json:"appName"`
		AppType int    `json:"appType"`
		Appid   int    `json:"appid"`
		IconURL string `json:"iconUrl"`
	}

	Data struct {
		Title string `json:"title"`
		Value string `json:"value"`
	}

	Button struct {
		Name   string `json:"name"`
		Action string `json:"action"`
	}

	Notification struct {
		AppInfo         AppInfo  `json:"appInfo"`
		Data            []Data   `json:"data"`
		Title           string   `json:"title"`
		Button          []Button `json:"button"`
		EmphasisKeyword string   `json:"emphasis_keyword"`
	}

	Meta struct {
		Notification Notification `json:"notification"`
	}
)
