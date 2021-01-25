# OPQBot
OPQBot

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
		SendType:   OPQBot.SendTypePicMsgByUrl,
		SendToType: OPQBot.SendToTypeFriend,
		ToUserUid:  2435932516,
		Content:    OPQBot.SendTypePicMsgByUrlContent{Content: "你好", PicUrl: "https://img-home.csdnimg.cn/images/20201124032511.png"},
	})
	time.Sleep(1*time.Hour)
}
```

更多请看 [wiki](https://github.com/mcoo/OPQBot/wiki)
