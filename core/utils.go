package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func Atoi(str string) int {
	num, _ := strconv.Atoi(str)
	return num
}

func GetMidText(pre string, suf string, str string) string {
	n := strings.Index(str, pre)
	if n == -1 {
		n = 0
	} else {
		n = n + len(pre)
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, suf)
	if m == -1 {
		m = len(str)
	}
	return string([]byte(str)[:m])
}

func SplitText(content string, length int) []string {
	reg, _ := regexp.Compile(fmt.Sprintf(".{1,%d}", length))
	result := reg.FindAllString(content, -1)
	return result
}

func FileRead(path string) []byte {
	res, _ := os.Open(path)
	defer res.Close()
	data, _ := ioutil.ReadAll(res)
	return data
}

func FileWrite(path string, content []byte) int {
	ioutil.WriteFile(path, content, 0644)
	return len(content)
}

func FileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func MakeStar(num int) string {
	str := ""
	for i := 0; i < 10; i++ {
		if num >= i {
			str = str + "★"
		} else {
			str = str + "☆"
		}
	}
	return str
}
