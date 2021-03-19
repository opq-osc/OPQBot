package OPQBot

import (
	"encoding/json"
	"errors"
	"github.com/asmcos/requests"
	"github.com/goinggo/mapstructure"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

type BotManager struct {
	QQ       int64
	SendChan chan SendMsgPack
	Running  bool
	OPQUrl   string
	onEvent  map[string]reflect.Value
	delayed  int
	locker   sync.RWMutex
}

func NewBotManager(QQ int64, OPQUrl string) BotManager {
	return BotManager{QQ: QQ, OPQUrl: OPQUrl, SendChan: make(chan SendMsgPack, 1024), onEvent: make(map[string]reflect.Value), locker: sync.RWMutex{}, delayed: 1000}
}

// 设置发送消息的时延 单位毫秒 默认1000
func (b *BotManager) SetSendDelayed(Millisecond int) {
	b.delayed = Millisecond
}

// 开始连接
func (b *BotManager) Start() error {
	b.Running = true
	go b.receiveSendPack()
	c, err := gosocketio.Dial(strings.ReplaceAll(b.OPQUrl, "http://", "ws://")+"/socket.io/?EIO=3&transport=websocket", transport.GetDefaultWebsocketTransport())
	if err != nil {
		return err
	}
	_ = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		// log.Println("连接成功！")
		f, ok := b.onEvent[EventNameOnConnected]
		if ok {
			f.Call([]reflect.Value{})
		}
	})
	_ = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		// log.Println("连接断开！")
		f, ok := b.onEvent[EventNameOnDisconnected]
		if ok {
			f.Call([]reflect.Value{})
		}
	})
	_ = c.On("OnGroupMsgs", func(h *gosocketio.Channel, args returnPack) {
		if args.CurrentQQ != b.QQ {
			return
		}
		b.locker.RLock()
		defer b.locker.RUnlock()
		f, ok := b.onEvent["OnGroupMsgs"]
		if ok {
			result := GroupMsgPack{}
			err = mapstructure.Decode(args.CurrentPacket.Data, &result)
			if err != nil {
				log.Println("解析包错误")
				return
			}
			f.Call([]reflect.Value{reflect.ValueOf(args.CurrentQQ), reflect.ValueOf(result)})
		}
		//log.Println(args)
	})
	_ = c.On("OnFriendMsgs", func(h *gosocketio.Channel, args returnPack) {
		if args.CurrentQQ != b.QQ {
			return
		}
		b.locker.RLock()
		defer b.locker.RUnlock()
		f, ok := b.onEvent["OnFriendMsgs"]
		if ok {
			result := FriendMsgPack{}
			err = mapstructure.Decode(args.CurrentPacket.Data, &result)
			if err != nil {
				log.Println("解析包错误")
				return
			}
			f.Call([]reflect.Value{reflect.ValueOf(args.CurrentQQ), reflect.ValueOf(result)})
		}
		//log.Println(args)
	})
	_ = c.On("OnEvents", func(h *gosocketio.Channel, args eventPack) {
		if args.CurrentQQ != b.QQ {
			return
		}
		e, ok := args.CurrentPacket.Data.(map[string]interface{})
		if !ok {
			log.Println("解析出错")
			return
		}
		e1, ok := e["EventName"].(string)
		if !ok {
			log.Println("解析出错")
			return
		}
		switch e1 {
		case EventNameOnGroupJoin:
			b.locker.RLock()
			defer b.locker.RUnlock()
			f, ok := b.onEvent[EventNameOnGroupJoin]
			if ok {
				result := GroupJoinPack{}
				err = mapstructure.Decode(args.CurrentPacket.Data, &result)
				if err != nil {
					log.Println("解析包错误")
					return
				}
				f.Call([]reflect.Value{reflect.ValueOf(args.CurrentQQ), reflect.ValueOf(result)})
			}
		case EventNameOnGroupAdmin:
			b.locker.RLock()
			defer b.locker.RUnlock()
			f, ok := b.onEvent[EventNameOnGroupAdmin]
			if ok {
				result := GroupAdminPack{}
				err = mapstructure.Decode(args.CurrentPacket.Data, &result)
				if err != nil {
					log.Println("解析包错误")
					return
				}
				f.Call([]reflect.Value{reflect.ValueOf(args.CurrentQQ), reflect.ValueOf(result)})
			}
		case EventNameOnGroupExit:
			b.locker.RLock()
			defer b.locker.RUnlock()
			f, ok := b.onEvent[EventNameOnGroupExit]
			if ok {
				result := GroupExitPack{}
				err = mapstructure.Decode(args.CurrentPacket.Data, &result)
				if err != nil {
					log.Println("解析包错误")
					return
				}
				f.Call([]reflect.Value{reflect.ValueOf(args.CurrentQQ), reflect.ValueOf(result)})
			}
		case EventNameOnGroupExitSuccess:
			b.locker.RLock()
			defer b.locker.RUnlock()
			f, ok := b.onEvent[EventNameOnGroupExitSuccess]
			if ok {
				result := GroupExitSuccessPack{}
				err = mapstructure.Decode(args.CurrentPacket.Data, &result)
				if err != nil {
					log.Println("解析包错误")
					return
				}
				f.Call([]reflect.Value{reflect.ValueOf(args.CurrentQQ), reflect.ValueOf(result)})
			}
		case EventNameOnGroupAdminSysNotify:
			b.locker.RLock()
			defer b.locker.RUnlock()
			f, ok := b.onEvent[EventNameOnGroupAdminSysNotify]
			if ok {
				result := GroupAdminSysNotifyPack{}
				err = mapstructure.Decode(args.CurrentPacket.Data, &result)
				if err != nil {
					log.Println("解析包错误")
					return
				}
				f.Call([]reflect.Value{reflect.ValueOf(args.CurrentQQ), reflect.ValueOf(result)})
			}
		case EventNameOnGroupRevoke:
			b.locker.RLock()
			defer b.locker.RUnlock()
			f, ok := b.onEvent[EventNameOnGroupRevoke]
			if ok {
				result := GroupRevokePack{}
				err = mapstructure.Decode(args.CurrentPacket.Data, &result)
				if err != nil {
					log.Println("解析包错误")
					return
				}
				f.Call([]reflect.Value{reflect.ValueOf(args.CurrentQQ), reflect.ValueOf(result)})
			}
		case EventNameOnGroupShut:
			b.locker.RLock()
			defer b.locker.RUnlock()
			f, ok := b.onEvent[EventNameOnGroupShut]
			if ok {
				result := GroupShutPack{}
				err = mapstructure.Decode(args.CurrentPacket.Data, &result)
				if err != nil {
					log.Println("解析包错误")
					return
				}
				f.Call([]reflect.Value{reflect.ValueOf(args.CurrentQQ), reflect.ValueOf(result)})
			}
		case EventNameOnGroupSystemNotify:
			b.locker.RLock()
			defer b.locker.RUnlock()
			f, ok := b.onEvent[EventNameOnGroupSystemNotify]
			if ok {
				result := GroupSystemNotifyPack{}
				err = mapstructure.Decode(args.CurrentPacket.Data, &result)
				if err != nil {
					log.Println("解析包错误")
					return
				}
				f.Call([]reflect.Value{reflect.ValueOf(args.CurrentQQ), reflect.ValueOf(result)})
			}
		default:
			b.locker.RLock()
			defer b.locker.RUnlock()
			f, ok := b.onEvent[EventNameOnOther]
			if ok {
				f.Call([]reflect.Value{reflect.ValueOf(args.CurrentQQ), reflect.ValueOf(args)})
			}
		}
	})
	return nil
}

// 发送消息函数
func (b *BotManager) Send(sendMsgPack SendMsgPack) {
	select {
	case b.SendChan <- sendMsgPack:
	default:
	}
}

// 停止
func (b *BotManager) Stop() {
	if !b.Running {
		return
	}
	b.Running = false
	close(b.SendChan)
}

// 撤回消息
func (b *BotManager) ReCallMsg(GroupID, MsgRandom int64, MsgSeq int) error {
	res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=PbMessageSvc.PbMsgWithDraw&qq="+strconv.FormatInt(b.QQ, 10), map[string]interface{}{"GroupID": GroupID, "MsgSeq": MsgSeq, "MsgRandom": MsgRandom})
	if err != nil {
		// log.Println(err.Error())
		return err
	}
	var result Result
	err = res.Json(&result)
	if err != nil {
		return err
	}
	if result.Ret != 0 {
		return errors.New("Error ")
	} else {
		return nil
	}
}

// 刷新Key
func (b *BotManager) RefreshKey() error {
	res, err := requests.Get(b.OPQUrl + "/v1/RefreshKeys?qq=" + strconv.FormatInt(b.QQ, 10))
	if err != nil {
		// log.Println(err.Error())
		return err
	}
	var result Result
	err = res.Json(&result)
	if err != nil {
		return err
	}
	if result.Ret != 0 {
		return errors.New("Error ")
	} else {
		return nil
	}
}

// 获取用户信息
func (b *BotManager) GetUserInfo(qq int64) (UserInfo, error) {
	var result UserInfo
	res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=SummaryCard.ReqSummaryCard&qq="+strconv.FormatInt(b.QQ, 10), map[string]int64{"UserID": qq})
	if err != nil {
		// log.Println(err.Error())
		return result, err
	}
	log.Println(res.Text())
	err = res.Json(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// QQ赞 次数
func (b *BotManager) Zan(qq int64, num int) int {
	var result Result
	success := 0
	for i := 0; i < num; i++ {
		res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=OidbSvc.0x7e5_4&qq="+strconv.FormatInt(b.QQ, 10), map[string]int64{"UserID": qq})
		if err == nil {
			err = res.Json(&result)
			if err != nil {
				break
			}
			if result.Ret == 0 {
				success += 1
			}
		}
		time.Sleep(500 * time.Microsecond)
	}
	return success
}

// At宏
func MacroAt(qqs []int64) string {
	var qqsStr []string
	for i := range qqs {
		qqsStr = append(qqsStr, strconv.FormatInt(qqs[i], 10))
	}
	return "[ATUSER(" + strings.Join(qqsStr, ",") + ")]"
}
func MacroAtAll() string {
	return "[ATALL()]"
}

func (b *BotManager) AddEvent(EventName string, f interface{}) error {
	fVal := reflect.ValueOf(f)
	if fVal.Kind() != reflect.Func {
		return errors.New("NotFuncError")
	}
	var okStruck string
	switch EventName {
	case EventNameOnFriendMessage:
		okStruck = "OPQBot.FriendMsgPack"
	case EventNameOnGroupMessage:
		okStruck = "OPQBot.GroupMsgPack"
	case EventNameOnGroupJoin:
		okStruck = "OPQBot.GroupJoinPack"
	case EventNameOnGroupAdmin:
		okStruck = "OPQBot.GroupAdminPack"
	case EventNameOnGroupExit:
		okStruck = "OPQBot.GroupExitPack"
	case EventNameOnGroupExitSuccess:
		okStruck = "OPQBot.GroupExitSuccessPack"
	case EventNameOnGroupAdminSysNotify:
		okStruck = "OPQBot.GroupAdminSysNotifyPack"
	case EventNameOnGroupRevoke:
		okStruck = "OPQBot.GroupRevokePack"
	case EventNameOnGroupShut:
		okStruck = "OPQBot.GroupShutPack"
	case EventNameOnGroupSystemNotify:
		okStruck = "OPQBot.GroupSystemNotifyPack"
	case EventNameOnDisconnected:
		okStruck = "ok"
	case EventNameOnConnected:
		okStruck = "ok"
	case EventNameOnOther:
		okStruck = "interface {}"
	default:
		return errors.New("Unknown EventName ")
	}

	if fVal.Type().NumIn() == 0 && okStruck == "ok" {
		b.locker.Lock()
		defer b.locker.Unlock()
		b.onEvent[EventName] = fVal
		return nil
	}
	//log.Println( fVal.Type().In(0).String())
	if fVal.Type().NumIn() != 2 || fVal.Type().In(1).String() != okStruck {
		return errors.New("FuncError, Your Function Should Have " + okStruck)
	}

	b.locker.Lock()
	defer b.locker.Unlock()
	b.onEvent[EventName] = fVal
	return nil

}

func (b *BotManager) receiveSendPack() {
	log.Println("QQ发送信息通道开启")
OuterLoop:
	for {
		if !b.Running {
			break
		}
		sendMsgPack := <-b.SendChan
		sendJsonPack := make(map[string]interface{})
		sendJsonPack["ToUserUid"] = sendMsgPack.ToUserUid
		switch content := sendMsgPack.Content.(type) {
		case SendTypeTextMsgContent:
			sendJsonPack["SendMsgType"] = "TextMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			sendJsonPack["Content"] = content.Content
		case SendTypeTextMsgContentPrivateChat:
			sendJsonPack["SendMsgType"] = "TextMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			sendJsonPack["Content"] = content.Content
			sendJsonPack["GroupID"] = content.Group
		case SendTypePicMsgByUrlContent:
			sendJsonPack["SendMsgType"] = "PicMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			sendJsonPack["PicUrl"] = content.PicUrl
			sendJsonPack["Content"] = content.Content
			sendJsonPack["FlashPic"] = content.Flash
		case SendTypePicMsgByUrlContentPrivateChat:
			sendJsonPack["SendMsgType"] = "PicMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			sendJsonPack["PicUrl"] = content.PicUrl
			sendJsonPack["Content"] = content.Content
			sendJsonPack["FlashPic"] = content.Flash
			sendJsonPack["GroupID"] = content.Group
		case SendTypePicMsgByLocalContent:
			sendJsonPack["SendMsgType"] = "PicMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			sendJsonPack["PicPath"] = content.Path
			sendJsonPack["Content"] = content.Content
			sendJsonPack["FlashPic"] = content.Flash
		case SendTypePicMsgByLocalContentPrivateChat:
			sendJsonPack["SendMsgType"] = "PicMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			sendJsonPack["PicPath"] = content.Path
			sendJsonPack["Content"] = content.Content
			sendJsonPack["FlashPic"] = content.Flash
			sendJsonPack["GroupID"] = content.Group
		case SendTypePicMsgByMd5Content:
			sendJsonPack["SendMsgType"] = "PicMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			sendJsonPack["PicMd5s"] = content.Md5
			sendJsonPack["Content"] = content.Content
			sendJsonPack["FlashPic"] = content.Flash
		case SendTypeVoiceByUrlContent:
			sendJsonPack["SendMsgType"] = "VoiceMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			sendJsonPack["VoiceUrl"] = content.VoiceUrl
		case SendTypeVoiceByUrlContentPrivateChat:
			sendJsonPack["SendMsgType"] = "VoiceMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			sendJsonPack["VoiceUrl"] = content.VoiceUrl
			sendJsonPack["GroupID"] = content.Group
		case SendTypeVoiceByLocalContent:
			sendJsonPack["SendMsgType"] = "VoiceMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			sendJsonPack["VoiceUrl"] = content.Path
		case SendTypeVoiceByLocalContentPrivateChat:
			sendJsonPack["SendMsgType"] = "VoiceMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			sendJsonPack["VoiceUrl"] = content.Path
			sendJsonPack["GroupID"] = content.Group
		case SendTypeXmlContent:
			sendJsonPack["SendMsgType"] = "XmlMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			sendJsonPack["Content"] = content.Content
		case SendTypeXmlContentPrivateChat:
			sendJsonPack["SendMsgType"] = "XmlMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			sendJsonPack["Content"] = content.Content
			sendJsonPack["GroupID"] = content.Group
		case SendTypeJsonContent:
			sendJsonPack["SendMsgType"] = "JsonMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			sendJsonPack["Content"] = content.Content
		case SendTypeJsonContentPrivateChat:
			sendJsonPack["SendMsgType"] = "JsonMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			sendJsonPack["Content"] = content.Content
			sendJsonPack["GroupID"] = content.Group
		case SendTypeForwordContent:
			sendJsonPack["SendMsgType"] = "ForwordMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			sendJsonPack["ForwordBuf"] = content.ForwordBuf
			sendJsonPack["ForwordField"] = content.ForwordField
		case SendTypeForwordContentPrivateChat:
			sendJsonPack["SendMsgType"] = "ForwordMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			sendJsonPack["ForwordBuf"] = content.ForwordBuf
			sendJsonPack["ForwordField"] = content.ForwordField
			sendJsonPack["GroupID"] = content.Group
		case SendTypeRelayContent:
			sendJsonPack["ReplayInfo"] = content.ReplayInfo
		case SendTypeRelayContentPrivateChat:
			sendJsonPack["SendMsgType"] = "ReplayMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			sendJsonPack["ReplayInfo"] = content.ReplayInfo
			sendJsonPack["GroupID"] = content.Group
		case SendTypePicMsgByBase64Content:
			sendJsonPack["SendMsgType"] = "PicMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			sendJsonPack["PicBase64Buf"] = content.Base64
			sendJsonPack["Content"] = content.Content
			sendJsonPack["FlashPic"] = content.Flash
		case SendTypePicMsgByBase64ContentPrivateChat:
			sendJsonPack["SendMsgType"] = "PicMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			sendJsonPack["PicBase64Buf"] = content.Base64
			sendJsonPack["Content"] = content.Content
			sendJsonPack["GroupID"] = content.Group
			sendJsonPack["FlashPic"] = content.Flash
		default:
			log.Println("未知发送的类型")
			continue OuterLoop
		}
		tmp, _ := json.Marshal(sendJsonPack)
		log.Println(string(tmp))
		res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=SendMsgV2&qq="+strconv.FormatInt(b.QQ, 10), sendJsonPack)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		log.Println(res.Text())
		time.Sleep(time.Duration(b.delayed) * time.Millisecond)
	}
}
