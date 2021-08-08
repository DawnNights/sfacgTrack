package core

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func strmid(pre string,suf string,str string) string {
	n := strings.Index(str, pre)
	if n == -1 {n = 0} else {n = n + len(pre)}
	str = string([]byte(str)[n:])
	m := strings.Index(str, suf)
	if m == -1 {m = len(str)}
	return string([]byte(str)[:m])
}

func (sf *SFAPI) SetCookie(cookie string) {
	FileWrite("resource/Cookie.txt",[]byte(cookie))
}

func (sf *SFAPI) GetCookie() string {
	return string(FileRead("resource/Cookie.txt"))
}

func (sf *SFAPI) Comment(bookid string, content string) string {
	client := &http.Client{}
	params := fmt.Sprintf("nid=%s&commentid=0&content=%s",bookid,url.QueryEscape(content))
	req,_ := http.NewRequest("POST","https://book.sfacg.com/ajax/ashx/Common.ashx?op=addCmt",strings.NewReader(params))
	req.Header = http.Header{
		"User-Agent": []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:84.0) Gecko/20100101 Firefox/84.0"},
		"Content-Type": []string{"application/x-www-form-urlencoded"},
		"Cookie": []string{sf.GetCookie()},
	}
	rsp,_ := client.Do(req)
	defer rsp.Body.Close()
	body,_ := ioutil.ReadAll(rsp.Body)
	return string(body)
}

