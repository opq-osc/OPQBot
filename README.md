# OPQBot
![GitHub go.mod Go version (branch & subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/mcoo/OPQBot/master?filename=go.mod&style=for-the-badge) ![GitHub release (latest SemVer including pre-releases)](https://img.shields.io/github/v/release/mcoo/OPQBot?include_prereleases&style=for-the-badge) ![GitHub tag (latest SemVer pre-release)](https://img.shields.io/github/v/tag/mcoo/OPQBot?include_prereleases&style=for-the-badge) ![GitHub](https://img.shields.io/github/license/mcoo/OPQBot?style=for-the-badge)
### 功能
|功能|是否实现|
|-|-|
|群消息处理事件|是|
|好友消息处理事件|是|
|机器人事件处理|是|
|所有支持的消息发送|是|
|At|是|
|表情|是|
|撤回|是|
### Install
`go get github.com/mcoo/OPQBot`
### Example
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

### 没人看的更新历史
```
20210318    简化发送代码
20210319    将宏移出BotManager,添加对发送队列每次发送时间的控制
20210322    添加发送函数的中间件
```
