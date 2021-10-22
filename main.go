package main

import (
	_ "sfacg/bot"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
)

func main() {
	zero.Run(zero.Config{
		NickName:      []string{"报更姬"},
		CommandPrefix: "",
		SuperUsers:    []string{"2224825532"},

		Driver: []zero.Driver{
			&driver.WSClient{
				// OneBot 正向WS 默认使用 6700 端口
				Url:         "ws://127.0.0.1:6700",
				AccessToken: "",
			},
		},
	})

	select {}
}
