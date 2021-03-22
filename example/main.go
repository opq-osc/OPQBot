package main

import (
	"encoding/base64"
	"fmt"
	"github.com/asmcos/requests"
	"github.com/mcoo/OPQBot"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

var ZanNote = map[int64]int{}

func main() {
	opqBot := OPQBot.NewBotManager(2629326992, "http://192.168.2.2:8899")
	// 设置发送队列每次发送的间隔时间 默认1000ms
	opqBot.SetSendDelayed(1000)
	err := opqBot.Start()
	if err != nil {
		log.Println(err.Error())
	}
	defer opqBot.Stop()
	log.Println(opqBot.RegMiddleware(1, func(m map[string]interface{}) {
		m["Content"] = "XXX"
	}))
	err = opqBot.AddEvent(OPQBot.EventNameOnGroupMessage, func(botQQ int64, packet OPQBot.GroupMsgPack) {
		if packet.FromUserID != opqBot.QQ {
			if packet.Content == "Base64图片测试" {
				pic, _ := ioutil.ReadFile("./test.jpg")
				opqBot.Send(OPQBot.SendMsgPack{
					SendToType: OPQBot.SendToTypeGroup,
					ToUserUid:  packet.FromGroupID,
					Content:    OPQBot.SendTypePicMsgByBase64Content{Content: "图片", Base64: base64.StdEncoding.EncodeToString(pic)},
				})
				return
			}
			if packet.Content == "赞我" {
				i, ok := ZanNote[packet.FromUserID]
				if ok {
					if i == time.Now().Day() {
						opqBot.Send(OPQBot.SendMsgPack{
							SendToType: OPQBot.SendToTypeGroup,
							ToUserUid:  packet.FromGroupID,
							Content:    OPQBot.SendTypeTextMsgContent{Content: "今日已赞!"},
						})
						return
					}
				}
				opqBot.Send(OPQBot.SendMsgPack{
					SendToType: OPQBot.SendToTypeGroup,
					ToUserUid:  packet.FromGroupID,
					Content:    OPQBot.SendTypeTextMsgContent{Content: "正在赞请稍后！"},
				})
				success := opqBot.Zan(packet.FromUserID, 50)
				opqBot.Send(OPQBot.SendMsgPack{
					SendToType: OPQBot.SendToTypeGroup,
					ToUserUid:  packet.FromGroupID,
					Content:    OPQBot.SendTypeTextMsgContent{Content: "成功赞了" + strconv.Itoa(success) + "次"},
				})
				ZanNote[packet.FromUserID] = time.Now().Day()
				return
			}
			if packet.Content == "二次元图片" {
				res, err := requests.Get("http://www.dmoe.cc/random.php?return=json")
				if err != nil {
					return
				}
				var pixivPic Pic
				_ = res.Json(&pixivPic)
				opqBot.Send(OPQBot.SendMsgPack{
					SendToType: OPQBot.SendToTypeGroup,
					ToUserUid:  int64(packet.FromGroupID),
					Content:    OPQBot.SendTypePicMsgByUrlContent{Content: "随机", PicUrl: pixivPic.Imgurl},
				})
				return
			}
			if packet.Content == "刷新" && packet.FromUserID == 2435932516 {
				err := opqBot.RefreshKey()
				if err != nil {
					log.Println(err.Error())
				}
			}
		}
		log.Println(botQQ, packet)
	})
	if err != nil {
		log.Println(err.Error())
	}
	err = opqBot.AddEvent(OPQBot.EventNameOnFriendMessage, func(botQQ int64, packet OPQBot.FriendMsgPack) {
		if packet.Content == "赞我" {
			i, ok := ZanNote[packet.FromUin]
			if ok {
				if i == time.Now().Day() {
					opqBot.Send(OPQBot.SendMsgPack{
						SendToType: OPQBot.SendToTypeFriend,
						ToUserUid:  packet.FromUin,
						Content:    OPQBot.SendTypeTextMsgContent{Content: OPQBot.MacroAt([]int64{packet.FromUin}) + "今日已赞!"},
					})
					return
				}
			}
			opqBot.Send(OPQBot.SendMsgPack{
				SendToType: OPQBot.SendToTypeFriend,
				ToUserUid:  packet.FromUin,
				Content:    OPQBot.SendTypeTextMsgContent{Content: "正在赞请稍后！"},
			})
			success := opqBot.Zan(packet.FromUin, 50)
			opqBot.Send(OPQBot.SendMsgPack{
				SendToType: OPQBot.SendToTypeFriend,
				ToUserUid:  packet.FromUin,
				Content:    OPQBot.SendTypeTextMsgContent{Content: "成功赞了" + strconv.Itoa(success) + "次"},
			})
			ZanNote[packet.FromUin] = time.Now().Day()
			return
		}
		if c := strings.Split(packet.Content, " "); len(c) >= 2 {
			if c[0] == "#查询" {
				log.Println(c[1])
				qq, err := strconv.ParseInt(c[1], 10, 64)
				if err != nil {
					opqBot.Send(OPQBot.SendMsgPack{
						SendToType: OPQBot.SendToTypeFriend,
						ToUserUid:  packet.FromUin,
						Content:    OPQBot.SendTypeTextMsgContent{Content: err.Error()},
					})
				}
				user, err := opqBot.GetUserInfo(qq)
				log.Println(user)
				if err != nil {
					opqBot.Send(OPQBot.SendMsgPack{
						SendToType: OPQBot.SendToTypeFriend,
						ToUserUid:  packet.FromUin,
						Content:    OPQBot.SendTypeTextMsgContent{Content: err.Error()},
					})
				} else {
					var sex string
					if user.Sex == 1 {
						sex = "女"
					} else {
						sex = "男"
					}
					opqBot.Send(OPQBot.SendMsgPack{
						SendToType: OPQBot.SendToTypeFriend,
						ToUserUid:  packet.FromUin,
						Content:    OPQBot.SendTypeTextMsgContent{Content: fmt.Sprintf("用户:%s[%d]%s\n来自:%s\n年龄:%d\n被赞了:%d次\n", user.NickName, user.QQUin, sex, user.Province+user.City, user.Age, user.LikeNums)},
					})
				}
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
	err = opqBot.AddEvent(OPQBot.EventNameOnOther, func(botQQ int64, e interface{}) {
		log.Println(e)
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

type Pic struct {
	Code   string `json:"code"`
	Imgurl string `json:"imgurl"`
	Width  string `json:"width"`
	Height string `json:"height"`
}
