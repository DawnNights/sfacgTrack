package core

import (
	"bytes"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func parseText(content string, lineLen int) []string {
	var i,idx int
	text := []rune(content)
	result := []string{}

	for idx, _ = range text {
		if idx != 0 && idx % lineLen == 0{
			result = append(result,string(text[i:idx]))
			i = idx
		}
	}
	if i != idx {
		result = append(result,string(text[i:idx]))
	}
	return result
}

func netFile(url string,filePath string) []byte {
	//发送POST文件请求获取返回数据
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)
	paths := strings.Split(filePath,"/");body_writer.CreateFormFile("files", paths[len(paths)-1])

	file,_ := os.Open(filePath)
	boundary := body_writer.Boundary()
	close_buf := bytes.NewBufferString("\r\n--"+boundary+"--\r\n")
	request_reader := io.MultiReader(body_buf, file, close_buf)
	defer file.Close()
	client := &http.Client{}
	request,_ := http.NewRequest("POST",url,request_reader)
	request.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)

	res,_ := client.Do(request)
	defer res.Body.Close()
	result,_ := ioutil.ReadAll(res.Body)
	return result
}

func (self *NovelInfo) Dawn() string {
	dc := gg.NewContext(904, 432)

	// 加载背景图片并画入
	im, _ := gg.LoadImage("resource/BackGround.jpg")
	dc.DrawImage(im, 0, 0)

	// 加载封面图片并画入，不存在则下载
	if !FileExist("resource/Cover/"+self.Id+".png"){
		req,_ := http.Get(self.CoverUrl)
		defer req.Body.Close()
		cover,_ := ioutil.ReadAll(req.Body)
		im, _, _ = image.Decode(bytes.NewReader(cover))
		im = resize.Resize(240,335, im, resize.Lanczos3)
		gg.SavePNG("resource/Cover/"+self.Id+".png",im)
	}
	im, _ = gg.LoadImage("resource/Cover/"+self.Id+".png")
	dc.DrawImage(im, 40, 100)

	dc.SetColor(color.White)
	// 写入小说名称
	dc.LoadFontFace("resource/FZYTK.TTF",26)
	dc.DrawString("《 " + self.Name + " 》",30,27)

	// 写入更新时间和更新章节
	dc.LoadFontFace("resource/FZYTK.TTF",23)
	dc.DrawString("【更新时间: " + self.Time + "】",30,61)
	dc.DrawString("【" + self.NewChapter.Title + "】",30,93)

	// 写入小说字数、收藏量和点击量
	dc.LoadFontFace("resource/FZYTK.TTF",20)
	dc.DrawString("字数: " + strconv.Itoa(self.Words),678,23)
	dc.DrawString("收藏: " + strconv.Itoa(self.Collection),678,53)
	dc.DrawString("点击: " + self.Hits,678,83)

	var viewText []string
	if self.IsVip {
		dc.LoadFontFace("resource/FZYTK.TTF",28)
		viewText = parseText(self.Preview,20)
		for i, text := range viewText {
			dc.DrawString(text,300,130+float64(i*40))
		}
	}else {
		dc.LoadFontFace("resource/FZYTK.TTF",24)
		viewText = parseText(self.Preview,24)
		for i, text := range viewText {
			dc.DrawString(text,300,130+float64(i*30))
		}
	}

	dc.SavePNG("resource/"+self.Id+".png")
	rsp := netFile("https://pic.sogou.com/pic/upload_pic.jsp","resource/"+self.Id+".png")
	os.Remove("resource/"+self.Id+".png")
	return string(rsp)
}