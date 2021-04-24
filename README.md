# OPQBot 🎉
![GitHub go.mod Go version (branch & subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/mcoo/OPQBot/master?filename=go.mod&style=for-the-badge&logo=go) ![GitHub tag (latest SemVer pre-release)](https://img.shields.io/github/v/tag/mcoo/OPQBot?include_prereleases&style=for-the-badge&logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAA1klEQVQ4T6XTvQ3CMBCG4fdKRgCxBQvQMgAZggUQFdDBCAxAAz1INGQAKCnYgR0O2VjBSc5JBCn98+Tusy2qegXG1L8HcAF2IvI05v2QqKqmJsP4C5iJyNFaFwM5sAeGwNJYnFlIDGxEZOE2quoJmHRBYmAtIqsA3IBRorVSJdUMzkAvEWrsFYgDmv7WlK9HHODKtkJrORw/nTmgD7gqBl12VNbcJYQ2BQ4/ALkH/kCyAvgB+YRYLVtVu7TzPUar7xakfJFSwSWQ2nuotRCDAZmHsa31mN4A6l46o4qtxAAAAABJRU5ErkJggg==) ![GitHub](https://img.shields.io/github/license/mcoo/OPQBot?style=for-the-badge&logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAABMklEQVQ4T5WTsS5EQRSGv18EhcgmOkqJVlQ6sZWGgkQhUW9hS4XS7jOICIUKkUi20CjxBrKdBxCFgifwy6y7nL3mbjjNnZz7zzfn/GdGhLC9B6yF1IOkVtSU1+onbO8CR8B9EK0AC5K6VZAISCcdANeFeApYBeqSInSAFQEN4CRz0iPwnsl3JB1+A5LA9l36Sqrbfis2JUA5RoBlYP8vgFRZP6bD+hhoxxa2gNMgmARGh0zgQtJOzsTYb+0/JlaNcVFSzodecWUPzoHZUPYr0JCUm8IPwHbatARslgDPQOr1duhFsn0JbANPwEsQzwDzQEtSOweR7TNgY9h9Bz6AK0nNX2/BtituWllbkzTgWc/EApDK60rq2J4r3kD6PwaMAxNFG5WAL0elBLwB1rP9Zir4BJmUbAFx6PbeAAAAAElFTkSuQmCC)
### 功能 😄
|功能|是否实现|
|-|-|
|群消息处理事件|是|
|好友消息处理事件|是|
|机器人事件处理|是|
|所有支持的消息发送|是|
|At|是|
|表情|是|
|撤回|是|
|下载文件|是|
|禁言|是|
|群公告|是|

### 安装 💡
` github.com/mcoo/OPQBot@latest`
[使用教程](https://mcenjoy.cn/330/)
### 例子 👆
```golang
package main

import (
	"github.com/mcoo/OPQBot"
	"log"
	"time"
)
func main()  {
	opqBot := OPQBot.NewBotManager(2629326992,"http://192.168.2.2:8899")
	err := opqBot.Start()
	if err != nil {
		log.Println(err.Error())
	}
	defer opqBot.Stop()
	err = opqBot.AddEvent(OPQBot.EventNameOnGroupMessage, func(botQQ int64, packet OPQBot.GroupMsgPack) {
		log.Println(botQQ,packet)
	})
	if err != nil {
		log.Println(err.Error())
	}
	err = opqBot.AddEvent(OPQBot.EventNameOnFriendMessage, func(botQQ int64, packet OPQBot.FriendMsgPack) {
		log.Println(botQQ,packet)
	})
	if err != nil {
		log.Println(err.Error())
	}
	err = opqBot.AddEvent(OPQBot.EventNameOnGroupShut, func(botQQ int64, packet OPQBot.GroupShutPack) {
		log.Println(botQQ,packet)
	})
	if err != nil {
		log.Println(err.Error())
	}
	opqBot.Send(OPQBot.SendMsgPack{
		SendToType: OPQBot.SendToTypeFriend,
		ToUserUid:  2435932516,
		Content:    OPQBot.SendTypePicMsgByUrlContent{Content: "你好", PicUrl: "https://img-home.csdnimg.cn/images/20201124032511.png"},
	})
	time.Sleep(1*time.Hour) // 可以用WaitGroup替代
}
```

更多请看 [wiki](https://github.com/mcoo/OPQBot/wiki)

以example文件为准 [example](https://github.com/opq-osc/OPQBot/blob/main/example/main.go)

### 没人看的更新历史 ✏️
```
20210318    简化发送代码
20210319    将宏移出BotManager,添加对发送队列每次发送时间的控制
20210322    添加发送函数的中间件
20210403    增加发送回调和优化中间件,基础功能完善
20210405    添加对撤回功能的完善和支持 注意看一下 example
20210406    戳一戳功能，example 即是文档
20210407    删除多余log，完善戳戳
20210420    添加Mp3转Silk功能和一些其他的功能
20210424    添加事件的中间件，向下兼容以前的代码，使用看example，完善silk功能
```
