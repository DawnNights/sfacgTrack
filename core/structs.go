package core

import (
	"io/ioutil"
	"os"
)

type (
	SFAPI uint8

	NovelInfo struct {
		Id string	// 小说书号
		Url string	// 小说网址
		Name string	// 小说书名
		Hits string	// 小说点击量
		Time string	// 小说更新时间
		Words int	// 小说字数
		IsVip bool	// 是否为上架小说
		Preview string	// 最新章节预览
		CoverUrl string	// 小说封面网址
		Collection int	// 小说收藏量
		NewChapter NovelChapter	// 最新章节信息
	}

	NovelChapter struct {
		Url string	// 章节网址
		Time string	// 章节更新时间
		Title string	// 章节名称
		Writer string	// 作者昵称
		Words int	// 章节字数
		LastUrl string	// 上一章网址
		NextUrl string	// 下一章网址
	}

	TrackInfo struct {
		WordsGap string	// 较上章字数差
		TimeGap string	// 较上章时间差
		Times string	// 更新次数
		Stars string	// 更新评星
		Comment	string	// 更新评语
	}

	TrackConfig []struct {
		BookId string	// 报更小说书号
		IsSend bool	// 是否需要发送至评论区
		GroupId []int64	// 需要发送的群号
		RecordUrl string	// 当前更新网址记录
	}
)

// 读取文件数据
func FileRead(path string) []byte {
	res,_ := os.Open(path)
	defer res.Close()
	data,_ := ioutil.ReadAll(res)
	return data
}

// 写入文件数据
func FileWrite(path string,content []byte) int {
	ioutil.WriteFile(path,content,0644)
	return len(content)
}

// 判断文件是否存在
func FileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}