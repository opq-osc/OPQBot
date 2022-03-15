package OPQBot

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/big"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/goinggo/mapstructure"
	gosocketio "github.com/mcoo/OPQBot/golang-socketio-edit"
	"github.com/mcoo/OPQBot/golang-socketio-edit/transport"
	"github.com/mcoo/OPQBot/session"
	_ "github.com/mcoo/OPQBot/session/provider"
	"github.com/mcoo/requests"
)

type BotManager struct {
	QQ             int64
	SendChan       chan SendMsgPack
	Running        bool
	OPQUrl         string
	MaxRetryCount  int
	Done           chan int
	wg             sync.WaitGroup
	myRecord       map[string]MyRecord
	myRecordLocker sync.RWMutex
	onEvent        map[string][][]reflect.Value
	middleware     []middleware
	delayed        int
	locker         sync.RWMutex
	restart        chan int
	Session        *session.Manager
}

type middleware struct {
	priority int
	fun      func(m map[string]interface{}) map[string]interface{}
}

func (b *BotManager) SetMaxRetryCount(maxRetryCount int) {
	b.MaxRetryCount = maxRetryCount
}

var interrupt chan os.Signal

func ParserGroupAtMsg(pack GroupMsgPack) (a AtMsg, e error) {
	if pack.MsgType != "AtMsg" {
		e = errors.New("Not AtMsg. ")
		return
	}
	e = json.Unmarshal([]byte(pack.Content), &a)
	if e != nil {
		return
	}
	return
}
func (a AtMsg) Clean() AtMsg {
	for _, v := range a.UserExt {
		a.Content = strings.TrimSpace(strings.ReplaceAll(a.Content, "@"+v.QQNick, ""))
	}
	return a
}
func ParserGroupReplyMsg(pack GroupMsgPack) (a Reply, e error) {
	if pack.MsgType != "AtMsg" {
		e = errors.New("Not ReplyMsg. ")
		return
	}
	e = json.Unmarshal([]byte(pack.Content), &a)
	if e != nil {
		return
	}
	return
}
func ParserGroupPicMsg(pack GroupMsgPack) (a PicMsg, e error) {
	if pack.MsgType != "PicMsg" {
		e = errors.New("Not PicMsg. ")
		return
	}
	e = json.Unmarshal([]byte(pack.Content), &a)
	if e != nil {
		return
	}
	return
}
func ParserGroupFileMsg(pack GroupMsgPack) (a GroupFileMsg, e error) {
	if pack.MsgType != "GroupFileMsg" {
		e = errors.New("Not GroupFileMsg. ")
		return
	}
	e = json.Unmarshal([]byte(pack.Content), &a)
	if e != nil {
		return
	}
	return
}
func ParserVideoMsg(pack GroupMsgPack) (a VideoMsg, e error) {
	if pack.MsgType != "VideoMsg" {
		e = errors.New("Not VideoMsg. ")
		return
	}
	e = json.Unmarshal([]byte(pack.Content), &a)
	if e != nil {
		return
	}
	return
}
func (b *BotManager) Wait() {
home:
	b.wg.Wait()
	if b.MaxRetryCount > 0 {
		for i := 0; i < b.MaxRetryCount; i++ {
			log.Info("等待重试,要终止请按下Ctrl+C")
			select {
			case <-b.Done:

				b.Running = false
				log.Info("Bot结束")
				return

			case <-time.After(5 * time.Second):
				log.Warn("正在重连")
			}
			log.Warningf("重连尝试第%d/%d次\n", i+1, b.MaxRetryCount)
			err := b.Start()
			if err != nil {
				log.Error(err)
			} else {
				goto home
			}
		}
	}
	b.Running = false
	log.Info("Bot结束")
}

// VoiceMp3ToSilk Mp3转Silk mp3->silk Output: base64 String
func VoiceMp3ToSilk(mp3Path string) (string, error) {
	n, _ := rand.Int(rand.Reader, big.NewInt(100000))
	name := n.String()
	pcmFile := name + ".tmp"
	silkFile := name + ".silk"
	cmd := exec.Command("./ffmpeg", "-i", mp3Path, "-ac", "1", "-ar", "24000", "-f", "s16le", pcmFile)
	var stderr bytes.Buffer
	_, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	defer os.Remove(pcmFile)
	cmd = exec.Command("./encoder", pcmFile, silkFile, "-tencent")
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil { //获取输出对象，可以从该对象中读取输出结果
		return "", err
	}
	defer os.Remove(silkFile)
	tresult, _ := ioutil.ReadFile(silkFile)
	return base64.StdEncoding.EncodeToString(tresult), nil
}

// VoiceSilkToMp3 Silk转Mp3 silk->mp3 Output: []byte
func VoiceSilkToMp3(base64EncodedSilk string) ([]byte, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(base64EncodedSilk)
	if err != nil {
		return decodeBytes, errors.New("转码失败! ")
	}
	n, _ := rand.Int(rand.Reader, big.NewInt(100000))
	name := n.String()
	mp3file := "./" + name + ".mp3"
	pcmfile := "./" + name + ".pcm"
	silk := "./" + name + ".silk"
	err = ioutil.WriteFile(silk, decodeBytes, os.FileMode(0777))
	if err != nil {
		return decodeBytes, errors.New("写文件失败!")
	}
	defer os.Remove(silk)
	var stderr bytes.Buffer
	cmd := exec.Command("./decoder", silk, pcmfile)
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return decodeBytes, errors.New("silk转pcm失败! ")
	}
	cmd = exec.Command("./ffmpeg", "-f", "s16le", "-ar", "24000", "-i", pcmfile, "-ac", "1", mp3file)

	_, err = cmd.CombinedOutput()

	if err != nil {
		return decodeBytes, errors.New("pcm转mp3失败! ")
	}
	defer os.Remove(mp3file)
	tresult, _ := ioutil.ReadFile(mp3file)
	defer os.Remove(mp3file)
	return tresult, nil
}

var log *logrus.Entry

func init() {
	interrupt = make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	l := logrus.New()
	l.SetLevel(logrus.InfoLevel)
	l.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log = l.WithField("Core", "OPQBOT")
}
func SetLog(l *logrus.Entry) {
	log = l
}
func NewBotManager(QQ int64, OPQUrl string) *BotManager {

	s, err := session.NewManager("qq", 3600)
	if err != nil {
		panic(err)
	}
	go s.GC()
	b := &BotManager{restart: make(chan int, 1), Session: s, Done: make(chan int, 10), MaxRetryCount: 10, wg: sync.WaitGroup{}, QQ: QQ, OPQUrl: OPQUrl, SendChan: make(chan SendMsgPack, 1024), onEvent: make(map[string][][]reflect.Value), myRecord: map[string]MyRecord{}, myRecordLocker: sync.RWMutex{}, locker: sync.RWMutex{}, delayed: 1000}
	go func() {
		for {
			select {
			case <-interrupt:
				log.Info("程序被用户终止,正在进行释放资源操作!")
				b.MaxRetryCount = 0
				b.Done <- 0
				b.Done <- 0
				b.Done <- 0
			case <-b.restart:
				log.Warn("程序重连尝试!")
				b.Done <- 1
				b.Done <- 2
			}
		}

	}()
	return b
}

// SetSendDelayed 设置发送消息的时延 单位毫秒 默认1000
func (b *BotManager) SetSendDelayed(Millisecond int) {
	b.delayed = Millisecond
}

// Start 开始连接
func (b *BotManager) Start() error {
	b.Running = true
	b.wg.Add(2)
	go b.receiveSendPack()
	go func() {
		for {
			select {
			case <-b.Done:
				b.wg.Done()
				return
			case <-time.After(10 * time.Second):
				go func() {
					if len(b.myRecord) > 50 {
						b.myRecordLocker.Lock()
						for i, v := range b.myRecord {
							if time.Since(time.Unix(int64(v.MsgTime), 0)) > time.Second*180 {
								delete(b.myRecord, i)
							}
						}
						b.myRecordLocker.Unlock()
					}
				}()
			}
		}
	}()
	c, err := gosocketio.Dial(strings.ReplaceAll(b.OPQUrl, "http://", "ws://")+"/socket.io/?EIO=3&transport=websocket", transport.GetDefaultWebsocketTransport())
	if err != nil {
		b.restart <- 1
		return err
	}
	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		// log.Println("连接成功！")
		f, ok := b.onEvent[EventNameOnConnected]
		if ok && len(f) >= 1 {
			f[0][0].Call([]reflect.Value{})
		}
	})
	if err != nil {
		b.restart <- 1
		return err
	}
	err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		// log.Println("连接断开！")
		f, ok := b.onEvent[EventNameOnDisconnected]
		if ok && len(f) >= 1 {
			f[0][0].Call([]reflect.Value{})
		}
		b.restart <- 1
	})
	if err != nil {
		b.restart <- 1
		return err
	}
	err = c.On("OnGroupMsgs", func(h *gosocketio.Channel, args returnPack) {
		//log.Println(args)
		if args.CurrentQQ != b.QQ {
			return
		}
		b.locker.RLock()
		defer b.locker.RUnlock()
		f, ok := b.onEvent["OnGroupMsgs"]
		if ok && len(f) >= 1 {
			result := GroupMsgPack{}
			err = mapstructure.Decode(args.CurrentPacket.Data, &result)
			if err != nil {
				log.Error("解析包错误")
				return
			}
			reg1, _ := regexp.Compile(`\[([0-9]{1,5})\]`)
			id := reg1.FindStringSubmatch(result.Content)
			if result.FromUserID == b.QQ && len(id) > 1 {
				go func() {
					record := MyRecord{
						FromGroupID: result.FromGroupID,
						MsgRandom:   result.MsgRandom,
						MsgSeq:      result.MsgSeq,
						MsgTime:     result.MsgTime,
						MsgType:     result.MsgType,
						Content:     result.Content,
					}
					b.myRecordLocker.Lock()
					b.myRecord[id[1]] = record
					b.myRecordLocker.Unlock()
				}()
			}
			for _, v := range f {
				if result.Ban {
					return
				}
				result.Bot = b
				result.f = v
				result.NowIndex = 0
				result.MaxIndex = len(v) - 1
				v[0].Call([]reflect.Value{reflect.ValueOf(args.CurrentQQ), reflect.ValueOf(&result)})
			}

		}
		//log.Println(args)
	})
	if err != nil {
		return err
	}
	err = c.On("OnFriendMsgs", func(h *gosocketio.Channel, args returnPack) {
		if args.CurrentQQ != b.QQ {
			return
		}
		b.locker.RLock()
		defer b.locker.RUnlock()
		f, ok := b.onEvent["OnFriendMsgs"]
		if ok && len(f) >= 1 {
			result := FriendMsgPack{}
			err = mapstructure.Decode(args.CurrentPacket.Data, &result)
			if err != nil {
				log.Error("解析包错误")
				return
			}
			for _, v := range f {
				if result.Ban {
					return
				}
				result.Bot = b
				result.f = v
				result.NowIndex = 0
				result.MaxIndex = len(v) - 1
				v[0].Call([]reflect.Value{reflect.ValueOf(args.CurrentQQ), reflect.ValueOf(&result)})
			}
		}
		//log.Println(args)
	})
	if err != nil {
		b.restart <- 1
		return err
	}
	err = c.On("OnEvents", func(h *gosocketio.Channel, args eventPack) {
		if args.CurrentQQ != b.QQ {
			return
		}
		e, ok := args.CurrentPacket.Data.(map[string]interface{})
		if !ok {
			log.Error("解析出错")
			return
		}
		e1, ok := e["EventName"].(string)
		if !ok {
			log.Error("解析出错")
			return
		}
		switch e1 {
		case EventNameOnGroupJoin:
			b.locker.RLock()
			defer b.locker.RUnlock()
			f, ok := b.onEvent[EventNameOnGroupJoin]
			if ok && len(f) >= 1 {
				result := GroupJoinPack{}
				err = mapstructure.Decode(args.CurrentPacket.Data, &result)
				if err != nil {
					log.Error("解析包错误")
					return
				}
				for _, v := range f {
					if result.Ban {
						return
					}
					result.Bot = b
					result.f = v
					result.NowIndex = 0
					result.MaxIndex = len(v) - 1
					v[0].Call([]reflect.Value{reflect.ValueOf(args.CurrentQQ), reflect.ValueOf(&result)})
				}
			}
		case EventNameOnGroupAdmin:
			b.locker.RLock()
			defer b.locker.RUnlock()
			f, ok := b.onEvent[EventNameOnGroupAdmin]
			if ok && len(f) >= 1 {
				result := GroupAdminPack{}
				err = mapstructure.Decode(args.CurrentPacket.Data, &result)
				if err != nil {
					log.Error("解析包错误")
					return
				}
				for _, v := range f {
					if result.Ban {
						return
					}
					result.Bot = b
					result.f = v
					result.NowIndex = 0
					result.MaxIndex = len(v) - 1
					v[0].Call([]reflect.Value{reflect.ValueOf(args.CurrentQQ), reflect.ValueOf(&result)})
				}
			}
		case EventNameOnGroupExit:
			b.locker.RLock()
			defer b.locker.RUnlock()
			f, ok := b.onEvent[EventNameOnGroupExit]
			if ok && len(f) >= 1 {
				result := GroupExitPack{}
				err = mapstructure.Decode(args.CurrentPacket.Data, &result)
				if err != nil {
					log.Error("解析包错误")
					return
				}
				for _, v := range f {
					if result.Ban {
						return
					}
					result.Bot = b
					result.f = v
					result.NowIndex = 0
					result.MaxIndex = len(v) - 1
					v[0].Call([]reflect.Value{reflect.ValueOf(args.CurrentQQ), reflect.ValueOf(&result)})
				}
			}
		case EventNameOnGroupExitSuccess:
			b.locker.RLock()
			defer b.locker.RUnlock()
			f, ok := b.onEvent[EventNameOnGroupExitSuccess]
			if ok && len(f) >= 1 {
				result := GroupExitSuccessPack{}
				err = mapstructure.Decode(args.CurrentPacket.Data, &result)
				if err != nil {
					log.Error("解析包错误")
					return
				}
				for _, v := range f {
					if result.Ban {
						return
					}
					result.Bot = b
					result.f = v
					result.NowIndex = 0
					result.MaxIndex = len(v) - 1
					v[0].Call([]reflect.Value{reflect.ValueOf(args.CurrentQQ), reflect.ValueOf(&result)})
				}
			}
		case EventNameOnGroupAdminSysNotify:
			b.locker.RLock()
			defer b.locker.RUnlock()
			f, ok := b.onEvent[EventNameOnGroupAdminSysNotify]
			if ok && len(f) >= 1 {
				result := GroupAdminSysNotifyPack{}
				err = mapstructure.Decode(args.CurrentPacket.Data, &result)
				if err != nil {
					log.Error("解析包错误")
					return
				}
				for _, v := range f {
					if result.Ban {
						return
					}
					result.Bot = b
					result.f = v
					result.NowIndex = 0
					result.MaxIndex = len(v) - 1
					v[0].Call([]reflect.Value{reflect.ValueOf(args.CurrentQQ), reflect.ValueOf(&result)})
				}
			}
		case EventNameOnGroupRevoke:
			b.locker.RLock()
			defer b.locker.RUnlock()
			f, ok := b.onEvent[EventNameOnGroupRevoke]
			if ok && len(f) >= 1 {
				result := GroupRevokePack{}
				err = mapstructure.Decode(args.CurrentPacket.Data, &result)
				if err != nil {
					log.Error("解析包错误")
					return
				}
				for _, v := range f {
					if result.Ban {
						return
					}
					result.Bot = b
					result.f = v
					result.NowIndex = 0
					result.MaxIndex = len(v) - 1
					v[0].Call([]reflect.Value{reflect.ValueOf(args.CurrentQQ), reflect.ValueOf(&result)})
				}
			}
		case EventNameOnGroupShut:
			b.locker.RLock()
			defer b.locker.RUnlock()
			f, ok := b.onEvent[EventNameOnGroupShut]
			if ok && len(f) >= 1 {
				result := GroupShutPack{}
				err = mapstructure.Decode(args.CurrentPacket.Data, &result)
				if err != nil {
					log.Error("解析包错误")
					return
				}
				for _, v := range f {
					if result.Ban {
						return
					}
					result.Bot = b
					result.f = v
					result.NowIndex = 0
					result.MaxIndex = len(v) - 1
					v[0].Call([]reflect.Value{reflect.ValueOf(args.CurrentQQ), reflect.ValueOf(&result)})
				}
			}
		case EventNameOnGroupSystemNotify:
			b.locker.RLock()
			defer b.locker.RUnlock()
			f, ok := b.onEvent[EventNameOnGroupSystemNotify]
			if ok && len(f) >= 1 {
				result := GroupSystemNotifyPack{}
				err = mapstructure.Decode(args.CurrentPacket.Data, &result)
				if err != nil {
					log.Error("解析包错误")
					return
				}
				for _, v := range f {
					if result.Ban {
						return
					}
					result.Bot = b
					result.f = v
					result.NowIndex = 0
					result.MaxIndex = len(v) - 1
					v[0].Call([]reflect.Value{reflect.ValueOf(args.CurrentQQ), reflect.ValueOf(&result)})
				}
			}
		default:
			b.locker.RLock()
			defer b.locker.RUnlock()
			f, ok := b.onEvent[EventNameOnOther]
			if ok && len(f) >= 1 {
				for _, v := range f {
					v[0].Call([]reflect.Value{reflect.ValueOf(args.CurrentQQ), reflect.ValueOf(args)})
				}

			}
		}
	})
	if err != nil {
		b.restart <- 1
		return err
	}
	return nil
}

// Send 发送消息函数
func (b *BotManager) Send(sendMsgPack SendMsgPack) {
	select {
	case b.SendChan <- sendMsgPack:
	default:
	}
}

// Stop 停止
func (b *BotManager) Stop() {
	if !b.Running {
		return
	}
	b.Running = false
	close(b.SendChan)
}

// ReCallMsg 撤回消息
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
		return errors.New("Error 撤回失败")
	} else {
		return nil
	}
}

// RefreshKey 刷新Key
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
		return errors.New("Error 刷新Key失败")
	} else {
		return nil
	}
}

// Announce 发公告 Pinned 1为置顶,0为普通公告 announceType 发布类型(10为使用弹窗公告,20为发送给新成员,其他暂未知)
func (b *BotManager) Announce(title, text string, pinned, announceType int, groupID int64) error {
	var result Result
	res, err := requests.PostJson(b.OPQUrl+"/v1/Group/Announce?qq="+strconv.FormatInt(b.QQ, 10), map[string]interface{}{"GroupID": groupID, "Title": title, "Text": text, "Pinned": pinned, "Type": announceType})
	if err != nil {
		return err
	}
	err = res.Json(&result)
	if err != nil {
		return err
	}
	if result.Ret != 0 {
		return errors.New(result.Msg)
	}
	return nil
}

// UploadFileWithBase64 上传群文件
func (b *BotManager) UploadFileWithBase64(FileName, FileBase64 string, ToUserUid int64, Notify bool) error {
	var result Result
	res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=SendMsgV2&qq="+strconv.FormatInt(b.QQ, 10), map[string]interface{}{"ToUserUid": ToUserUid, "Notify": Notify, "FileName": FileName, "FileBase64": FileBase64, "SendMsgType": "UploadGroupFile"})
	if err != nil {
		return err
	}
	err = res.Json(&result)
	if err != nil {
		return err
	}
	if result.Ret != 0 {
		return errors.New(result.Msg)
	}
	return nil
}

// UploadFileWithFileUrl 上传群文件
func (b *BotManager) UploadFileWithFileUrl(FileName, FileUrl string, ToUserUid int64, Notify bool) error {
	var result Result
	res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=SendMsgV2&qq="+strconv.FormatInt(b.QQ, 10), map[string]interface{}{"ToUserUid": ToUserUid, "Notify": Notify, "FileName": FileName, "FileUrl": FileUrl, "SendMsgType": "UploadGroupFile"})
	if err != nil {
		return err
	}
	err = res.Json(&result)
	if err != nil {
		return err
	}
	if result.Ret != 0 {
		return errors.New(result.Msg)
	}
	return nil
}

// UploadFileWithFilePath 上传群文件
func (b *BotManager) UploadFileWithFilePath(FilePath string, ToUserUid int64, Notify bool) error {
	var result Result
	res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=SendMsgV2&qq="+strconv.FormatInt(b.QQ, 10), map[string]interface{}{"ToUserUid": ToUserUid, "Notify": Notify, "FilePath": FilePath, "SendMsgType": "UploadGroupFile"})
	if err != nil {
		return err
	}
	err = res.Json(&result)
	if err != nil {
		return err
	}
	if result.Ret != 0 {
		return errors.New(result.Msg)
	}
	return nil
}

// Chuo 戳戳 sendType  0戳好友 1戳群友 sendType=0 时可以不填此字段 sendType=1 时不能为空
func (b *BotManager) Chuo(sendType int, groupID, userId int64) error {
	var result Result
	res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=OidbSvc.0xed3_1&qq="+strconv.FormatInt(b.QQ, 10), map[string]interface{}{"GroupID": groupID, "UserID": userId, "type": sendType})
	if err != nil {
		return err
	}
	err = res.Json(&result)
	if err != nil {
		return err
	}
	if result.Ret != 0 {
		return errors.New(result.Msg)
	}
	return nil
}

// SetAdmin 设置管理员 flag 1为设置管理员 2为取消管理员
func (b *BotManager) SetAdmin(flag int, groupID, userId int64) error {
	var result Result
	res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=OidbSvc.0x55c_1&qq="+strconv.FormatInt(b.QQ, 10), map[string]interface{}{"GroupID": groupID, "UserID": userId, "Flag": flag})
	if err != nil {
		return err
	}
	err = res.Json(&result)
	if err != nil {
		return err
	}
	if result.Ret != 0 {
		return errors.New(result.Msg)
	}
	return nil
}

// GetUserInfo 获取用户信息
func (b *BotManager) GetUserInfo(qq int64) (UserInfo, error) {
	var result UserInfo
	res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=GetUserInfo&qq="+strconv.FormatInt(b.QQ, 10), map[string]int64{"UserID": qq})
	if err != nil {
		// log.Println(err.Error())
		return result, err
	}
	err = res.Json(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetUserCookie 获取QQ相关ck
func (b *BotManager) GetUserCookie() (Cookie, error) {
	var result Cookie
	res, err := requests.Get(b.OPQUrl + "/v1/LuaApiCaller?funcname=GetUserCook&qq=" + strconv.FormatInt(b.QQ, 10))
	if err != nil {
		// log.Println(err.Error())
		return result, err
	}
	err = res.Json(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// KickGroupMember 踢出群成员
func (b *BotManager) KickGroupMember(groupID, userId int64) error {
	var result Result
	res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=GroupMgr&qq="+strconv.FormatInt(b.QQ, 10), map[string]interface{}{"GroupID": groupID, "ActionUserID": userId, "ActionType": 3, "Content": ""})
	if err != nil {
		return err
	}
	err = res.Json(&result)
	if err != nil {
		return err
	}
	if result.Ret != 0 {
		return errors.New(result.Msg)
	}
	return nil
}

// GetGroupMemberList 获取群成员列表
func (b *BotManager) GetGroupMemberList(groupID, LastUin int64) (GroupMemberList, error) {
	var result GroupMemberList
	res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=friendlist.GetTroopMemberListReq&qq="+strconv.FormatInt(b.QQ, 10), map[string]interface{}{"GroupUin": groupID, "LastUin": LastUin})
	if err != nil {
		return result, err
	}
	err = res.Json(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// SetGroupNewNick 设置群名片
func (b *BotManager) SetGroupNewNick(newNick string, groupID, userId int64) error {
	var result Result
	res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=friendlist.ModifyGroupCardReq&qq="+strconv.FormatInt(b.QQ, 10), map[string]interface{}{"GroupID": groupID, "UserID": userId, "NewNick": newNick})
	if err != nil {
		return err
	}
	err = res.Json(&result)
	if err != nil {
		return err
	}
	if result.Ret != 0 {
		return errors.New(result.Msg)
	}
	return nil
}

// SetGroupUniqueTitle 设置群头衔
func (b *BotManager) SetGroupUniqueTitle(newNick string, groupID, userId int64) error {
	var result Result
	res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=OidbSvc.0x8fc_2&qq="+strconv.FormatInt(b.QQ, 10), map[string]interface{}{"GroupID": groupID, "UserID": userId, "NewNick": newNick})
	if err != nil {
		return err
	}
	err = res.Json(&result)
	if err != nil {
		return err
	}
	if result.Ret != 0 {
		return errors.New(result.Msg)
	}
	return nil
}

// GetFriendList 获取好友列表
func (b *BotManager) GetFriendList(startIndex int) (FriendList, error) {
	var result FriendList
	res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=friendlist.GetFriendListReq&timeout=10&qq="+strconv.FormatInt(b.QQ, 10), map[string]interface{}{"StartIndex": startIndex})
	if err != nil {
		return result, err
	}
	err = res.Json(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetGroupList 获取群列表
func (b *BotManager) GetGroupList(nextToken string) (GroupList, error) {
	var result GroupList
	res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=friendlist.GetTroopListReqV2&timeout=10&qq="+strconv.FormatInt(b.QQ, 10), map[string]interface{}{"NextToken": nextToken})
	if err != nil {
		return result, err
	}
	err = res.Json(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// SetForbidden 设置禁言 flag 0为设置全体禁言 1为设置某人禁言 ShutTime 0为取消禁言 >0为禁言分钟数 全体禁言>0为开启禁言
func (b *BotManager) SetForbidden(flag, ShutTime int, groupID, userId int64) error {
	var result Result
	if flag == 0 {
		Switch := 0
		if ShutTime > 0 {
			Switch = 1
		}
		res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=OidbSvc.0x89a_0&qq="+strconv.FormatInt(b.QQ, 10), map[string]interface{}{"GroupID": groupID, "Switch": Switch})
		if err != nil {
			return err
		}
		err = res.Json(&result)
		if err != nil {
			return err
		}
		if result.Ret != 0 {
			return errors.New(result.Msg)
		}
	} else {
		res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=OidbSvc.0x570_8&qq="+strconv.FormatInt(b.QQ, 10), map[string]interface{}{"GroupID": groupID, "ShutUpUserID": userId, "ShutTime": ShutTime})
		if err != nil {
			return err
		}
		err = res.Json(&result)
		if err != nil {
			return err
		}
		if result.Ret != 0 {
			return errors.New(result.Msg)
		}
	}
	return nil
}

// GetFile 下载文件 groupId 为0 是下载好友分享文件
func (b *BotManager) GetFile(fileId string, groupID int64) (FriendFileResult, GroupFileResult, error) {
	var friendFileResult FriendFileResult
	var groupFileResult GroupFileResult
	if groupID == 0 {
		res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=OfflineFilleHandleSvr.pb_ftn_CMD_REQ_APPLY_DOWNLOAD-1200&qq="+strconv.FormatInt(b.QQ, 10), map[string]interface{}{"FileID": fileId})
		if err != nil {
			return friendFileResult, groupFileResult, err
		}
		err = res.Json(&friendFileResult)
		if err != nil {
			return friendFileResult, groupFileResult, err
		}
	} else {
		res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=OidbSvc.0x6d6_2&qq="+strconv.FormatInt(b.QQ, 10), map[string]interface{}{"FileID": fileId, "GroupID": groupID})
		if err != nil {
			return friendFileResult, groupFileResult, err
		}
		err = res.Json(&groupFileResult)
		if err != nil {
			return friendFileResult, groupFileResult, err
		}
	}
	return friendFileResult, groupFileResult, nil
}

// GetUserCardInfo 获取用户信息
func (b *BotManager) GetUserCardInfo(qq int64) (UserCardInfo, error) {
	var result UserCardInfo
	res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=SummaryCard.ReqSummaryCard&qq="+strconv.FormatInt(b.QQ, 10), map[string]int64{"UserID": qq})
	if err != nil {
		// log.Println(err.Error())
		return result, err
	}
	err = res.Json(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// OldSendVoice 发送语音 旧版 将被移出
func (b *BotManager) OldSendVoice(userID int64, sendToType int, data string) error {
	//var result Result
	_, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=SendMsg&qq="+strconv.FormatInt(b.QQ, 10), map[string]interface{}{"toUser": userID, "sendToType": sendToType, "sendMsgType": "VoiceMsg", "content": "",
		"groupid":        0,
		"atUser":         0,
		"voiceUrl":       "",
		"voiceBase64Buf": data,
	})
	if err != nil {
		// log.Println(err.Error())
		return err
	}
	//err = res.Json(&result)
	//if err != nil {
	//	return result, err
	//}
	return nil
}

// Zan QQ赞 次数
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
func MacroId() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(100000))
	keyRecord := n.String()
	return "[" + keyRecord + "]"
}

// MacroAt At宏
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

func (b *BotManager) AddEvent(EventName string, f ...interface{}) (func(), error) {
	var events []reflect.Value
	if len(f) == 0 {
		return nil, errors.New("调用错误")
	}
	for _, v := range f {
		fVal := reflect.ValueOf(v)
		if fVal.Kind() != reflect.Func {
			return nil, errors.New("NotFuncError")
		}
		var okStruck string
		switch EventName {
		case EventNameOnFriendMessage:
			okStruck = "*OPQBot.FriendMsgPack"
		case EventNameOnGroupMessage:
			okStruck = "*OPQBot.GroupMsgPack"
		case EventNameOnGroupJoin:
			okStruck = "*OPQBot.GroupJoinPack"
		case EventNameOnGroupAdmin:
			okStruck = "*OPQBot.GroupAdminPack"
		case EventNameOnGroupExit:
			okStruck = "*OPQBot.GroupExitPack"
		case EventNameOnGroupExitSuccess:
			okStruck = "*OPQBot.GroupExitSuccessPack"
		case EventNameOnGroupAdminSysNotify:
			okStruck = "*OPQBot.GroupAdminSysNotifyPack"
		case EventNameOnGroupRevoke:
			okStruck = "*OPQBot.GroupRevokePack"
		case EventNameOnGroupShut:
			okStruck = "*OPQBot.GroupShutPack"
		case EventNameOnGroupSystemNotify:
			okStruck = "*OPQBot.GroupSystemNotifyPack"
		case EventNameOnDisconnected:
			okStruck = "ok"
		case EventNameOnConnected:
			okStruck = "ok"
		case EventNameOnOther:
			okStruck = "interface {}"
		default:
			return nil, errors.New("Unknown EventName ")
		}

		if fVal.Type().NumIn() == 0 && okStruck == "ok" {
			events = append(events, fVal)
			continue
		}
		if fVal.Type().NumIn() != 2 || fVal.Type().In(1).String() != okStruck {
			return nil, errors.New(EventName + ": FuncError, Your Function  Should Have " + okStruck + " Your Struct is " + fVal.Type().In(1).String())
		}

		events = append(events, fVal)
	}
	b.locker.Lock()
	defer b.locker.Unlock()
	b.onEvent[EventName] = append(b.onEvent[EventName], events)
	return func() {
		b.locker.Lock()
		defer b.locker.Unlock()
		for i, v := range b.onEvent[EventName] {
			if len(v) > 0 && v[0] == reflect.ValueOf(f[0]) {
				if len(b.onEvent[EventName]) == 1 {
					delete(b.onEvent, EventName)
					break
				}
				b.onEvent[EventName] = append(b.onEvent[EventName][:i], b.onEvent[EventName][i+1:]...)
			}
		}
	}, nil

}

// RegSendMiddleware 注册 发送函数的中间件 2为最先执行 0为最后执行
func (b *BotManager) RegSendMiddleware(priority int, f func(m map[string]interface{}) map[string]interface{}) error {
	fVal := reflect.ValueOf(f)
	if fVal.Kind() != reflect.Func {
		return errors.New("NotFuncError")
	}
	if priority < 0 || priority > 2 {
		return errors.New("priority should >= 0 and <= 2 ")
	}
	if fVal.Type().NumIn() != 1 {
		return errors.New("Error ")
	}
	middle := middleware{
		priority: priority,
		fun:      f,
	}
	b.locker.Lock()
	defer b.locker.Unlock()
	b.middleware = append(b.middleware, middle)
	return nil
}
func (b *BotManager) CallFunc(FuncName string, funcStruct string) ([]byte, error) {
	res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname="+FuncName+"&qq="+strconv.FormatInt(b.QQ, 10), funcStruct)
	if err != nil {
		return nil, err
	}
	return res.Content(), nil
}
func (b *BotManager) receiveSendPack() {
	log.Info("QQ发送信息通道开启")
OuterLoop:
	for {
		select {
		case <-b.Done:
			log.Info("发信通道关闭")
			b.wg.Done()
			return
		case sendMsgPack := <-b.SendChan:
			record := MyRecord{
				FromGroupID: 0,
				MsgRandom:   0,
				MsgSeq:      0,
				MsgTime:     0,
				MsgType:     "",
				Content:     "",
			}
			sendJsonPack := make(map[string]interface{})
			sendJsonPack["ToUserUid"] = sendMsgPack.ToUserUid
			record.FromGroupID = sendMsgPack.ToUserUid
			switch content := sendMsgPack.Content.(type) {
			case SendTypeTextMsgContent:
				sendJsonPack["SendMsgType"] = "TextMsg"
				sendJsonPack["SendToType"] = sendMsgPack.SendToType
				sendJsonPack["Content"] = content.Content
				record.Content = content.Content
				record.MsgType = "TextMsg"
			case SendTypeTextMsgContentPrivateChat:
				sendJsonPack["SendMsgType"] = "TextMsg"
				sendJsonPack["SendToType"] = sendMsgPack.SendToType
				sendJsonPack["Content"] = content.Content
				sendJsonPack["GroupID"] = content.Group
				record.Content = content.Content
				record.MsgType = "TextMsg"
			case SendTypePicMsgByUrlContent:
				sendJsonPack["SendMsgType"] = "PicMsg"
				sendJsonPack["SendToType"] = sendMsgPack.SendToType
				sendJsonPack["PicUrl"] = content.PicUrl
				sendJsonPack["Content"] = content.Content
				sendJsonPack["FlashPic"] = content.Flash
				record.Content = content.Content
				record.MsgType = "PicMsg"
			case SendTypePicMsgByUrlContentPrivateChat:
				sendJsonPack["SendMsgType"] = "PicMsg"
				sendJsonPack["SendToType"] = sendMsgPack.SendToType
				sendJsonPack["PicUrl"] = content.PicUrl
				sendJsonPack["Content"] = content.Content
				sendJsonPack["FlashPic"] = content.Flash
				sendJsonPack["GroupID"] = content.Group
				record.Content = content.Content
				record.MsgType = "PicMsg"
			case SendTypePicMsgByLocalContent:
				sendJsonPack["SendMsgType"] = "PicMsg"
				sendJsonPack["SendToType"] = sendMsgPack.SendToType
				sendJsonPack["PicPath"] = content.Path
				sendJsonPack["Content"] = content.Content
				sendJsonPack["FlashPic"] = content.Flash
				record.Content = content.Content
				record.MsgType = "PicMsg"
			case SendTypePicMsgByLocalContentPrivateChat:
				sendJsonPack["SendMsgType"] = "PicMsg"
				sendJsonPack["SendToType"] = sendMsgPack.SendToType
				sendJsonPack["PicPath"] = content.Path
				sendJsonPack["Content"] = content.Content
				sendJsonPack["FlashPic"] = content.Flash
				sendJsonPack["GroupID"] = content.Group
				record.Content = content.Content
				record.MsgType = "PicMsg"
			case SendTypePicMsgByMd5Content:
				sendJsonPack["SendMsgType"] = "PicMsg"
				sendJsonPack["SendToType"] = sendMsgPack.SendToType
				sendJsonPack["PicMd5s"] = content.Md5
				sendJsonPack["Content"] = content.Content
				sendJsonPack["FlashPic"] = content.Flash
				record.Content = content.Content
				record.MsgType = "PicMsg"
			case SendTypeVoiceByUrlContent:
				sendJsonPack["SendMsgType"] = "VoiceMsg"
				sendJsonPack["SendToType"] = sendMsgPack.SendToType
				sendJsonPack["VoiceUrl"] = content.VoiceUrl
				record.MsgType = "VoiceMsg"
			case SendTypeVoiceByUrlContentPrivateChat:
				sendJsonPack["SendMsgType"] = "VoiceMsg"
				sendJsonPack["SendToType"] = sendMsgPack.SendToType
				sendJsonPack["VoiceUrl"] = content.VoiceUrl
				sendJsonPack["GroupID"] = content.Group
				record.MsgType = "VoiceMsg"
			case SendTypeVoiceByLocalContent:
				sendJsonPack["SendMsgType"] = "VoiceMsg"
				sendJsonPack["SendToType"] = sendMsgPack.SendToType
				sendJsonPack["VoiceUrl"] = content.Path
				record.MsgType = "VoiceMsg"
			case SendTypeVoiceByLocalContentPrivateChat:
				sendJsonPack["SendMsgType"] = "VoiceMsg"
				sendJsonPack["SendToType"] = sendMsgPack.SendToType
				sendJsonPack["VoiceUrl"] = content.Path
				sendJsonPack["GroupID"] = content.Group
				record.MsgType = "VoiceMsg"
			case SendTypeXmlContent:
				sendJsonPack["SendMsgType"] = "XmlMsg"
				sendJsonPack["SendToType"] = sendMsgPack.SendToType
				sendJsonPack["Content"] = content.Content
				record.Content = content.Content
				record.MsgType = "XmlMsg"
			case SendTypeXmlContentPrivateChat:
				sendJsonPack["SendMsgType"] = "XmlMsg"
				sendJsonPack["SendToType"] = sendMsgPack.SendToType
				sendJsonPack["Content"] = content.Content
				sendJsonPack["GroupID"] = content.Group
				record.Content = content.Content
				record.MsgType = "XmlMsg"
			case SendTypeJsonContent:
				sendJsonPack["SendMsgType"] = "JsonMsg"
				sendJsonPack["SendToType"] = sendMsgPack.SendToType
				sendJsonPack["Content"] = content.Content
			case SendTypeJsonContentPrivateChat:
				sendJsonPack["SendMsgType"] = "JsonMsg"
				sendJsonPack["SendToType"] = sendMsgPack.SendToType
				sendJsonPack["Content"] = content.Content
				sendJsonPack["GroupID"] = content.Group
				record.Content = content.Content
				record.MsgType = "JsonMsg"
			case SendTypeForwordContent:
				sendJsonPack["SendMsgType"] = "ForwordMsg"
				sendJsonPack["SendToType"] = sendMsgPack.SendToType
				sendJsonPack["ForwordBuf"] = content.ForwordBuf
				sendJsonPack["ForwordField"] = content.ForwordField
				sendJsonPack["Content"] = content.Content
				record.MsgType = "ForwordMsg"
			case SendTypeForwordContentPrivateChat:
				sendJsonPack["SendMsgType"] = "ForwordMsg"
				sendJsonPack["SendToType"] = sendMsgPack.SendToType
				sendJsonPack["ForwordBuf"] = content.ForwordBuf
				sendJsonPack["ForwordField"] = content.ForwordField
				sendJsonPack["GroupID"] = content.Group
				sendJsonPack["Content"] = content.Content
				record.MsgType = "ForwordMsg"

			case SendTypeReplyContent:
				sendJsonPack["SendMsgType"] = "ReplayMsg"
				sendJsonPack["ReplayInfo"] = content.ReplayInfo
				sendJsonPack["SendToType"] = sendMsgPack.SendToType
				sendJsonPack["Content"] = content.Content
				record.MsgType = "ReplayMsg"
			case SendTypeReplyContentPrivateChat:
				sendJsonPack["SendMsgType"] = "ReplayMsg"
				sendJsonPack["Content"] = content.Content
				sendJsonPack["SendToType"] = sendMsgPack.SendToType
				sendJsonPack["ReplayInfo"] = content.ReplayInfo
				sendJsonPack["GroupID"] = content.Group
				record.MsgType = "ReplayMsg"
			case SendTypePicMsgByBase64Content:
				sendJsonPack["SendMsgType"] = "PicMsg"
				sendJsonPack["SendToType"] = sendMsgPack.SendToType
				sendJsonPack["PicBase64Buf"] = content.Base64
				sendJsonPack["Content"] = content.Content
				sendJsonPack["FlashPic"] = content.Flash
				record.Content = content.Content
				record.MsgType = "PicMsg"
			case SendTypePicMsgByBase64ContentPrivateChat:
				sendJsonPack["SendMsgType"] = "PicMsg"
				sendJsonPack["SendToType"] = sendMsgPack.SendToType
				sendJsonPack["PicBase64Buf"] = content.Base64
				sendJsonPack["Content"] = content.Content
				sendJsonPack["GroupID"] = content.Group
				sendJsonPack["FlashPic"] = content.Flash
				record.Content = content.Content
				record.MsgType = "PicMsg"
			default:
				log.Warn("未知发送的类型")
				continue OuterLoop
			}
			for i := 2; i >= 0; i-- {
				for _, v := range b.middleware {
					if len(sendJsonPack) == 0 {
						break
					}
					if v.priority == i {
						sendJsonPack = v.fun(sendJsonPack)
						//v.fun.Call([]reflect.Value{reflect.ValueOf(sendJsonPack)})
					}
					r, ok := sendJsonPack["reason"].(string)
					if len(sendJsonPack) == 1 && ok {
						if r != "" {
							log.Info("消息被拦截！拦截原因 " + r)
						} else {
							log.Info("消息被拦截！无拦截原因")
						}
						continue OuterLoop
					}
				}
			}

			//tmp, _ := json.Marshal(sendJsonPack)
			//log.Println(string(tmp))
			res, err := requests.PostJson(b.OPQUrl+"/v1/LuaApiCaller?funcname=SendMsgV2&qq="+strconv.FormatInt(b.QQ, 10), sendJsonPack)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			var result Result
			err = res.Json(&result)
			if err != nil {
				log.Warn("发送失败！ ", err.Error())
				continue
			}
			reg1, _ := regexp.Compile(`\[([0-9]{1,5})\]`)
			id := reg1.FindStringSubmatch(record.Content)
			if sendMsgPack.CallbackFunc != nil {
				go func() {
					ch := make(chan MyRecord, 1)
					stop := make(chan bool, 1)
					go func() {
						if sendMsgPack.SendToType == SendToTypeFriend || len(id) <= 1 {
							ch <- MyRecord{}
							return
						}

						for {
							select {
							case <-stop:
								return
							default:
								b.myRecordLocker.Lock()
								if v, ok := b.myRecord[id[1]]; ok {
									ch <- v
									delete(b.myRecord, id[1])
								}
								b.myRecordLocker.Unlock()
							}
							time.Sleep(100 * time.Millisecond)

						}
					}()
					select {
					case myRecordPack := <-ch:
						sendMsgPack.CallbackFunc(result.Ret, result.Msg, myRecordPack)
						stop <- true

					case <-time.After(10 * time.Second):
						sendMsgPack.CallbackFunc(result.Ret, result.Msg, MyRecord{})
						stop <- true
					}

				}()
			}
			time.Sleep(time.Duration(b.delayed) * time.Millisecond)
		}
	}
}

// SendFriendTextMsg 发送文字信息给好友
func (b *BotManager) SendFriendTextMsg(FriendUin int64, Content string) {
	b.Send(SendMsgPack{
		SendToType: SendToTypeFriend,
		ToUserUid:  FriendUin,
		Content: SendTypeTextMsgContent{
			Content: Content,
		},
	})
}

// SendFriendPicMsg 发送图片信息给好友
func (b *BotManager) SendFriendPicMsg(FriendUin int64, Content string, Pic []byte) {
	b.Send(SendMsgPack{
		SendToType: SendToTypeFriend,
		ToUserUid:  FriendUin,
		Content: SendTypePicMsgByBase64Content{
			Content: Content,
			Base64:  base64.StdEncoding.EncodeToString(Pic),
			Flash:   false,
		},
	})
}

// SendGroupTextMsg 发送文字信息给群
func (b *BotManager) SendGroupTextMsg(GroupUin int64, Content string) {
	b.Send(SendMsgPack{
		SendToType: SendToTypeGroup,
		ToUserUid:  GroupUin,
		Content: SendTypeTextMsgContent{
			Content: Content,
		},
	})
}

// SendGroupPicMsg 发送图片信息给群
func (b *BotManager) SendGroupPicMsg(GroupUin int64, Content string, Pic []byte) {
	b.Send(SendMsgPack{
		SendToType: SendToTypeGroup,
		ToUserUid:  GroupUin,
		Content: SendTypePicMsgByBase64Content{
			Content: Content,
			Base64:  base64.StdEncoding.EncodeToString(Pic),
			Flash:   false,
		},
	})
}

// SendGroupJsonMsg 发送JSON信息给群
func (b *BotManager) SendGroupJsonMsg(GroupUin int64, Content string) {
	b.Send(SendMsgPack{
		SendToType: SendToTypeGroup,
		ToUserUid:  GroupUin,
		Content: SendTypeJsonContent{
			Content: Content,
		},
	})
}

// SendGroupXmlMsg 发送Xml信息给群
func (b *BotManager) SendGroupXmlMsg(GroupUin int64, Content string) {
	b.Send(SendMsgPack{
		SendToType: SendToTypeGroup,
		ToUserUid:  GroupUin,
		Content: SendTypeXmlContent{
			Content: Content,
		},
	})
}
