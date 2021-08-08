package main
import (
        "github.com/wdvxdr1123/ZeroBot"
        _ "sfacg/bot"
)


func main() {

        zero.Run(zero.Option{
                Host:          "127.0.0.1", // cqhttp的ip地址
                Port:          "6700", // cqhttp的端口
                AccessToken:   "",
                NickName:      []string{"机器人"},
                CommandPrefix: "", //指令前缀
                SuperUsers:    []string{"2224825532"},
        })
        select {}

}