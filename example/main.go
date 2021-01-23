package main

import (
	"github.com/mcoo/OPQBot"
	"log"
	"time"
)

func main() {
	opqBot := OPQBot.NewBotManager(2629326992, "http://192.168.2.2:8899")
	err := opqBot.Start()
	if err != nil {
		log.Println(err.Error())
	}
	defer opqBot.Stop()
	err = opqBot.AddEvent(OPQBot.EventNameOnGroupMessage, func(botQQ int64, packet OPQBot.GroupMsgPack) {
		log.Println(botQQ, packet)
	})
	if err != nil {
		log.Println(err.Error())
	}
	err = opqBot.AddEvent(OPQBot.EventNameOnFriendMessage, func(botQQ int64, packet OPQBot.FriendMsgPack) {
		log.Println(botQQ, packet)
	})
	if err != nil {
		log.Println(err.Error())
	}
	err = opqBot.AddEvent(OPQBot.EventNameOnGroupShut, func(botQQ int64, packet OPQBot.GroupShutPack) {
		log.Println(botQQ, packet)
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
	time.Sleep(1 * time.Hour)
}
