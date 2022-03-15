package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mcoo/OPQBot"
	"github.com/mcoo/requests"
)

var ZanNote = map[int64]int{}

func main() {
	// log.Println(OPQBot.DecodeFaceFromSentences("[表情123]啦啦啦[表情0]","[表情:%s]"))
	if len(os.Args) != 2 {
		return
	}
	opqBot := OPQBot.NewBotManager(2629326992, os.Args[1])
	// 设置发送队列每次发送的间隔时间 默认1000ms
	// 设置最大重试次数
	opqBot.SetMaxRetryCount(5)
	err := opqBot.Start()
	if err != nil {
		log.Println(err.Error())
	}
	defer opqBot.Stop()

	//log.Println(opqBot.RegMiddleware(1, func(m map[string]interface{}) map[string]interface{} {
	//	//m["Content"] = "测试"
	//	m = map[string]interface{}{"reason": "消息违规"}
	//	return m
	//}))
	//ck, _ := opqBot.GetUserCookie()
	//qz := qzone.NewQzoneManager(opqBot.QQ, ck)
	//f, _ := ioutil.ReadFile("./head.PNG")
	//u, _ := qz.UploadPic(base64.StdEncoding.EncodeToString(f))
	//bo, rich, _ := qzone.GetPicBoAndRichVal(u)
	//log.Println(qz.SendShuoShuoWithPic("发送图文测试", bo, rich))
	//lists,_ :=qz.GetShuoShuoList()
	//infoReg,_ := regexp.Compile(`<div class="f-info">(.*?)</div>`)
	//for _,v := range lists.Data.Data {
	//	if m := infoReg.FindStringSubmatch(v["html"].(string));len(m) == 2 {
	//		log.Println(m[1])
	//	}
	//}
	//qz.SendShuoShuo("发送文字测试")
	//log.Println(infoReg.FindStringSubmatch(lists.Data.Data[0]["html"].(string))[1])

	//log.Println(ck.PSkey.Qzone)
	var cancel func()
	cancel, err = opqBot.AddEvent(OPQBot.EventNameOnGroupMessage, VerifyBlackList, func(botQQ int64, packet *OPQBot.GroupMsgPack) {
		if packet.FromUserID != opqBot.QQ {
			s := opqBot.Session.SessionStart(packet.FromUserID)
			//last, _ := s.GetString("last")
			//if last != "" {
			//	opqBot.Send(OPQBot.SendMsgPack{
			//		SendToType: OPQBot.SendToTypeGroup,
			//		ToUserUid:  packet.FromGroupID,
			//		Content:    OPQBot.SendTypeTextMsgContent{Content: OPQBot.MacroAt([]int64{packet.FromUserID}) + "你上次发言为" + last},
			//	})
			//}
			s.Set("last", packet.Content)
			if packet.Content == "#上传测试" {
				//b,_ := ioutil.ReadFile("./1.mp3")base64.StdEncoding.EncodeToString(b)
				log.Println(opqBot.UploadFileWithBase64("1.mp3", "MTIzMTIzMTIzMjEz", packet.FromGroupID, true))
			}
			if packet.Content == "取消事件监听" {
				cancel()
			}
			if packet.Content == "#silk" {
				b, err := OPQBot.VoiceMp3ToSilk("./secret base ~君がくれたもの~ (10 years after Ver.).mp3")
				if err != nil {
					log.Println(err.Error())
					return
				}
				opqBot.OldSendVoice(packet.FromGroupID, 2, b)

			}
			if packet.Content == "#公告测试" {
				fmt.Println(opqBot.Announce("公告测试", "内容", 0, 10, packet.FromGroupID))
				return
			}
			if packet.Content == "#刷新" {
				_ = opqBot.RefreshKey()
				return
			}
			if packet.Content == "#Base64图片测试" {
				pic, _ := ioutil.ReadFile("./test.jpg")
				opqBot.Send(OPQBot.SendMsgPack{
					SendToType: OPQBot.SendToTypeGroup,
					ToUserUid:  packet.FromGroupID,
					Content:    OPQBot.SendTypePicMsgByBase64Content{Content: "图片", Base64: base64.StdEncoding.EncodeToString(pic)},
				})
				return
			}
			// 只有消息内容中含有宏OPQBot.MacroId() record 中才有消息的值，才能去用于撤回消息！
			if packet.Content == "#撤回测试" {
				opqBot.Send(OPQBot.SendMsgPack{
					SendToType: OPQBot.SendToTypeGroup,
					ToUserUid:  packet.FromGroupID,
					Content:    OPQBot.SendTypeTextMsgContent{Content: OPQBot.MacroAt([]int64{packet.FromUserID}) + "20s撤回测试！\n" + OPQBot.MacroId()},
					CallbackFunc: func(Code int, Info string, record OPQBot.MyRecord) {
						time.Sleep(20 * time.Second)
						_ = opqBot.ReCallMsg(record.FromGroupID, record.MsgRandom, record.MsgSeq)
					},
				})
				opqBot.Send(OPQBot.SendMsgPack{
					SendToType: OPQBot.SendToTypeGroup,
					ToUserUid:  packet.FromGroupID,
					Content:    OPQBot.SendTypeTextMsgContent{Content: OPQBot.MacroAt([]int64{packet.FromUserID}) + "20s撤回测试！\n" + OPQBot.MacroId()},
					CallbackFunc: func(Code int, Info string, record OPQBot.MyRecord) {
						time.Sleep(20 * time.Second)
						_ = opqBot.ReCallMsg(record.FromGroupID, record.MsgRandom, record.MsgSeq)
					},
				})
				opqBot.Send(OPQBot.SendMsgPack{
					SendToType: OPQBot.SendToTypeGroup,
					ToUserUid:  packet.FromGroupID,
					Content:    OPQBot.SendTypeTextMsgContent{Content: OPQBot.MacroAt([]int64{packet.FromUserID}) + "20s撤回测试！\n" + OPQBot.MacroId()},
					CallbackFunc: func(Code int, Info string, record OPQBot.MyRecord) {
						time.Sleep(20 * time.Second)
						_ = opqBot.ReCallMsg(record.FromGroupID, record.MsgRandom, record.MsgSeq)
					},
				})
			}
			if packet.Content == "#赞我" {
				i, ok := ZanNote[packet.FromUserID]
				if ok {
					if i == time.Now().Day() {
						opqBot.Send(OPQBot.SendMsgPack{
							SendToType: OPQBot.SendToTypeGroup,
							ToUserUid:  packet.FromGroupID,
							Content:    OPQBot.SendTypeTextMsgContent{Content: "今日已赞!"},
							CallbackFunc: func(Code int, Info string, record OPQBot.MyRecord) {
								log.Println(record)
								if Code == 0 {
									log.Println("发送成功")
								} else {
									log.Println("发送失败 错误代码", Code, Info)
								}
							},
						})
						return
					}
				}
				opqBot.Send(OPQBot.SendMsgPack{
					SendToType: OPQBot.SendToTypeGroup,
					ToUserUid:  packet.FromGroupID,
					Content:    OPQBot.SendTypeTextMsgContent{Content: "正在赞请稍后！"},
					CallbackFunc: func(Code int, Info string, record OPQBot.MyRecord) {
						log.Println(record)
						if Code == 0 {
							log.Println("发送成功")
						} else {
							log.Println("发送失败 错误代码", Code, Info)
						}
					},
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
					ToUserUid:  packet.FromGroupID,
					Content:    OPQBot.SendTypePicMsgByUrlContent{Content: "随机", PicUrl: pixivPic.Imgurl},
					CallbackFunc: func(Code int, Info string, record OPQBot.MyRecord) {
						if Code == 0 {
							log.Println("发送成功")
						} else {
							log.Println("发送失败 错误代码", Code, Info)
						}
					},
				})
				return
			}
			if packet.Content == "刷新" && packet.FromUserID == 2435932516 {
				err := opqBot.RefreshKey()
				if err != nil {
					log.Println(err.Error())
				}
			}
			if packet.Content == "戳我" {
				log.Println("chuo")
				err := opqBot.Chuo(1, packet.FromGroupID, packet.FromUserID)
				if err != nil {
					log.Println(err.Error())
				}
				log.Println(err)
			}
		}
		log.Println(botQQ, packet)
	})
	if err != nil {
		log.Println(err.Error())
	}
	_, err = opqBot.AddEvent(OPQBot.EventNameOnFriendMessage, func(botQQ int64, packet *OPQBot.FriendMsgPack) {
		if packet.Content == "赞我" {
			i, ok := ZanNote[packet.FromUin]
			if ok {
				// Todo 添加回调函数??
				// CallbackFunc: func(Code int,Info string) {
				//
				//	},
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
				user, err := opqBot.GetUserCardInfo(qq)
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
	_, err = opqBot.AddEvent(OPQBot.EventNameOnGroupShut, func(botQQ int64, packet *OPQBot.GroupShutPack) {
		log.Println(botQQ, packet)
	})
	if err != nil {
		log.Println(err.Error())
	}
	_, err = opqBot.AddEvent(OPQBot.EventNameOnConnected, func() {
		log.Println("连接成功！！！")
	})
	if err != nil {
		log.Println(err.Error())
	}
	_, err = opqBot.AddEvent(OPQBot.EventNameOnDisconnected, func() {
		log.Println("连接断开！！")
	})
	if err != nil {
		log.Println(err.Error())
	}
	_, err = opqBot.AddEvent(OPQBot.EventNameOnOther, func(botQQ int64, e interface{}) {
		log.Println(e)
	})
	_, err = opqBot.AddEvent(OPQBot.EventNameOnGroupSystemNotify, func(botQQ int64, e *OPQBot.GroupSystemNotifyPack) {
		log.Println(e)
	})
	_, err = opqBot.AddEvent(OPQBot.EventNameOnGroupRevoke, func(botQQ int64, e *OPQBot.GroupRevokePack) {
		log.Println(e)
	})
	_, err = opqBot.AddEvent(OPQBot.EventNameOnGroupJoin, func(botQQ int64, e *OPQBot.GroupJoinPack) {
		log.Println(e)
	})
	_, err = opqBot.AddEvent(OPQBot.EventNameOnGroupAdmin, func(botQQ int64, e *OPQBot.GroupAdminPack) {
		log.Println(e)
	})
	_, err = opqBot.AddEvent(OPQBot.EventNameOnGroupExit, func(botQQ int64, e *OPQBot.GroupExitPack) {
		log.Println(e)
	})
	_, err = opqBot.AddEvent(OPQBot.EventNameOnGroupExitSuccess, func(botQQ int64, e *OPQBot.GroupExitSuccessPack) {
		log.Println(e)
	})
	_, err = opqBot.AddEvent(OPQBot.EventNameOnGroupAdminSysNotify, func(botQQ int64, e *OPQBot.GroupAdminSysNotifyPack) {
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
	opqBot.Wait()
}
func VerifyBlackList(botQQ int64, packet *OPQBot.GroupMsgPack) {
	if packet.FromUserID == 123123123 {
		log.Println("触发黑名单")
		return
	}
	packet.Next(botQQ, packet)
}

type Pic struct {
	Code   string `json:"code"`
	Imgurl string `json:"imgurl"`
	Width  string `json:"width"`
	Height string `json:"height"`
}
