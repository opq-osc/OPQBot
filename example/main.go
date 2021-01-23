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
		if packet.FromUserID != opqBot.QQ {
			if packet.Content == "#菜单" {
				opqBot.Send(OPQBot.SendMsgPack{
					SendType:   OPQBot.SendTypeTextMsg,
					SendToType: OPQBot.SendToTypeGroup,
					ToUserUid:  int64(packet.FromGroupID),
					Content:    OPQBot.SendTypeTextMsgContent{Content: `HelloWorld`},
				})
			}
		}
		log.Println(botQQ, packet)
	})
	if err != nil {
		log.Println(err.Error())
	}
	err = opqBot.AddEvent(OPQBot.EventNameOnFriendMessage, func(botQQ int64, packet OPQBot.FriendMsgPack) {
		if packet.FromUin != opqBot.QQ {
			if packet.Content == "1" {
				opqBot.Send(OPQBot.SendMsgPack{
					SendType:   OPQBot.SendTypeTextMsg,
					SendToType: OPQBot.SendToTypeFriend,
					ToUserUid:  int64(packet.FromUin),
					Content:    OPQBot.SendTypeTextMsgContent{Content: "你好"},
				})
			}
		}
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
	err = opqBot.AddEvent(OPQBot.EventNameOnConnected, func() {
		log.Println("连接成功！！！")
	})
	if err != nil {
		log.Println(err.Error())
	}
	err = opqBot.AddEvent(OPQBot.EventNameOnDisconnected, func() {
		log.Println("连接断开！！")
	})
	if err != nil {
		log.Println(err.Error())
	}
	//opqBot.Send(OPQBot.SendMsgPack{
	//	SendType:   OPQBot.SendTypePicMsgByUrl,
	//	SendToType: OPQBot.SendToTypeFriend,
	//	ToUserUid:  2435932516,
	//	Content:    OPQBot.SendTypePicMsgByUrlContent{Content: "你好", PicUrl: "https://img-home.csdnimg.cn/images/20201124032511.png"},
	//})
	time.Sleep(1 * time.Hour)
}
