package main
import (
	_ "sfacg/bot"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
)


func main() {
	zero.Run(zero.Config{
		NickName:      []string{"bot"},
		CommandPrefix: "",
		SuperUsers:    []string{},
		Driver: []zero.Driver{
			driver.NewWebSocketClient("127.0.0.1", "6700", ""),
		},
	})
	select {}
}