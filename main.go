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
	locker   sync.RWMutex
}

func NewBotManager(QQ int64, OPQUrl string) BotManager {
	return BotManager{QQ: QQ, OPQUrl: OPQUrl, SendChan: make(chan SendMsgPack, 1024), onEvent: make(map[string]reflect.Value), locker: sync.RWMutex{}}
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
func (b *BotManager) At(qqs []int64) string {
	var qqs_str []string
	for i := range qqs_str {
		qqs_str = append(qqs_str, qqs_str[i])
	}
	return "[ATUSER(" + strings.Join(qqs_str, ",") + ")]"
}
func (b *BotManager) AtAll() string {
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
	default:
		return errors.New("Unknown EventName ")
	}

	if fVal.Type().NumIn() == 0 && okStruck == "ok" {
		b.locker.Lock()
		defer b.locker.Unlock()
		b.onEvent[EventName] = fVal
		return nil
	}
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
		switch sendMsgPack.SendType {
		case SendTypeTextMsg:
			sendJsonPack["SendMsgType"] = "TextMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			switch content := sendMsgPack.Content.(type) {
			case SendTypeTextMsgContent:
				sendJsonPack["Content"] = content.Content
			case SendTypeTextMsgContentPrivateChat:
				sendJsonPack["Content"] = content.Content
				sendJsonPack["GroupID"] = content.Group
			default:
				log.Println("类型不匹配")
				continue OuterLoop
			}
		case SendTypePicMsgByUrl:
			sendJsonPack["SendMsgType"] = "PicMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			switch Content := sendMsgPack.Content.(type) {
			case SendTypePicMsgByUrlContent:
				sendJsonPack["PicUrl"] = Content.PicUrl
				sendJsonPack["Content"] = Content.Content
				sendJsonPack["FlashPic"] = Content.Flash
			case SendTypePicMsgByUrlContentPrivateChat:
				sendJsonPack["PicUrl"] = Content.PicUrl
				sendJsonPack["Content"] = Content.Content
				sendJsonPack["GroupID"] = Content.Group
				sendJsonPack["FlashPic"] = Content.Flash
			default:
				log.Println("类型不匹配")
				continue OuterLoop
			}
		case SendTypePicMsgByLocal:
			sendJsonPack["SendMsgType"] = "PicMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			switch Content := sendMsgPack.Content.(type) {
			case SendTypePicMsgByLocalContent:
				sendJsonPack["PicPath"] = Content.Path
				sendJsonPack["Content"] = Content.Content
				sendJsonPack["FlashPic"] = Content.Flash
			case SendTypePicMsgByLocalContentPrivateChat:
				sendJsonPack["PicPath"] = Content.Path
				sendJsonPack["Content"] = Content.Content
				sendJsonPack["GroupID"] = Content.Group
				sendJsonPack["FlashPic"] = Content.Flash
			default:
				log.Println("类型不匹配")
				continue OuterLoop
			}
		case SendTypePicMsgByMd5:
			sendJsonPack["SendMsgType"] = "PicMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			switch Content := sendMsgPack.Content.(type) {
			case SendTypePicMsgByMd5Content:
				sendJsonPack["PicMd5s"] = Content.Md5
				sendJsonPack["Content"] = Content.Content
				sendJsonPack["FlashPic"] = Content.Flash
			case SendTypePicMsgByMd5ContentPrivateChat:
				sendJsonPack["PicMd5s"] = Content.Md5s
				sendJsonPack["Content"] = Content.Content
				sendJsonPack["GroupID"] = Content.Group
				sendJsonPack["FlashPic"] = Content.Flash
			default:
				log.Println("类型不匹配")
				continue OuterLoop
			}
		case SendTypeVoiceByUrl:
			sendJsonPack["SendMsgType"] = "VoiceMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			switch Content := sendMsgPack.Content.(type) {
			case SendTypeVoiceByUrlContent:
				sendJsonPack["VoiceUrl"] = Content.VoiceUrl
			case SendTypeVoiceByUrlContentPrivateChat:
				sendJsonPack["VoiceUrl"] = Content.VoiceUrl
				sendJsonPack["GroupID"] = Content.Group
			default:
				log.Println("类型不匹配")
				continue OuterLoop
			}
		case SendTypeVoiceByLocal:
			sendJsonPack["SendMsgType"] = "VoiceMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			switch Content := sendMsgPack.Content.(type) {
			case SendTypeVoiceByLocalContent:
				sendJsonPack["VoiceUrl"] = Content.Path
			case SendTypeVoiceByLocalContentPrivateChat:
				sendJsonPack["VoiceUrl"] = Content.Path
				sendJsonPack["GroupID"] = Content.Group
			default:
				log.Println("类型不匹配")
				continue OuterLoop
			}
		case SendTypeXml:
			sendJsonPack["SendMsgType"] = "XmlMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			switch Content := sendMsgPack.Content.(type) {
			case SendTypeXmlContent:
				sendJsonPack["Content"] = Content.Content
			case SendTypeXmlContentPrivateChat:
				sendJsonPack["Content"] = Content.Content
				sendJsonPack["GroupID"] = Content.Group
			default:
				log.Println("类型不匹配")
				continue OuterLoop
			}
		case SendTypeJson:
			sendJsonPack["SendMsgType"] = "XmlMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			switch Content := sendMsgPack.Content.(type) {
			case SendTypeJsonContent:
				sendJsonPack["Content"] = Content.Content
			case SendTypeJsonContentPrivateChat:
				sendJsonPack["Content"] = Content.Content
				sendJsonPack["GroupID"] = Content.Group
			default:
				log.Println("类型不匹配")
				continue OuterLoop
			}
		case SendTypeForword:
			sendJsonPack["SendMsgType"] = "ForwordMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			switch Content := sendMsgPack.Content.(type) {
			case SendTypeForwordContent:
				sendJsonPack["ForwordBuf"] = Content.ForwordBuf
				sendJsonPack["ForwordField"] = Content.ForwordField
			case SendTypeForwordContentPrivateChat:
				sendJsonPack["ForwordBuf"] = Content.ForwordBuf
				sendJsonPack["ForwordField"] = Content.ForwordField
				sendJsonPack["GroupID"] = Content.Group
			default:
				log.Println("类型不匹配")
				continue OuterLoop
			}
		case SendTypeReplay:
			sendJsonPack["SendMsgType"] = "ReplayMsg"
			sendJsonPack["SendToType"] = sendMsgPack.SendToType
			switch Content := sendMsgPack.Content.(type) {
			case SendTypeRelayContent:
				sendJsonPack["ReplayInfo"] = Content.ReplayInfo
			case SendTypeRelayContentPrivateChat:
				sendJsonPack["ReplayInfo"] = Content.ReplayInfo
				sendJsonPack["GroupID"] = Content.Group
			default:
				log.Println("类型不匹配")
				continue OuterLoop
			}
		}
		tmp, _ := json.Marshal(sendJsonPack)
		log.Println(string(tmp))
		res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=SendMsgV2&qq="+strconv.FormatInt(b.QQ, 10), sendJsonPack)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		log.Println(res.Text())
		time.Sleep(1 * time.Second)
	}
}
