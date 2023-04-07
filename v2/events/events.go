package events

import (
	"context"
	"encoding/json"
	"github.com/opq-osc/OPQBot/v2/apiBuilder"
	"strings"
)

type EventName string

const (
	EventNameNewMsg        EventName = "ON_EVENT_QQNT_NEW_MSG"
	EventNameGroupMsg      EventName = "ON_EVENT_GROUP_NEW_MSG"
	EventNameFriendMsg     EventName = "ON_EVENT_FRIEND_NEW_MSG"
	EventNameGroupJoin     EventName = "ON_EVENT_GROUP_JOIN"
	EventNameGroupExit     EventName = "ON_EVENT_GROUP_EXIT"
	EventNameGroupInvite   EventName = "ON_EVENT_GROUP_EXIT"
	EventNameLoginSuccess  EventName = "ON_EVENT_LOGIN_SUCCESS"
	EventNameNetworkChange EventName = "ON_EVENT_NETWORK_CHANGE"
)

type EventCallbackFunc func(ctx context.Context, event IEvent)

type IEvent interface {
	GetCurrentQQ() int64
	GetRawBytes() []byte
	GetEventName() EventName
	ParseGroupMsg() IGroupMsg
	ParseLoginSuccessEvent() ILoginSuccess
	GetApiBuilder() apiBuilder.IMainFunc
	ExcludeBot() IEvent
}
type IGroupMsg interface {
	ICommonMsg
	ExcludeAtInfo() IGroupMsg
	AtBot() bool
	GetGroupUin() int64
	GetSenderUin() int64
	ParseTextMsg() ITextMsg
}
type ITextMsg interface {
	GetTextContent() string
}
type ILoginSuccess interface {
	GetLoginSuccessBot() (nick string, uin int64)
}
type ICommonMsg interface {
	GetMsgUid() int64
	GetMsgType() int
	GetMsgTime() int64
}

func New(apiUrl string, data []byte) (IEvent, error) {
	event := &EventStruct{apiUrl: apiUrl}
	event.rawEvent = data
	return event, json.Unmarshal(data, event)
}

type EventStruct struct {
	apiUrl        string
	rawEvent      []byte
	CurrentPacket struct {
		EventData struct {
			Nick    *string `json:"Nick,omitempty"`
			Uin     *int64  `json:"Uin,omitempty"`
			MsgHead *struct {
				FromUin    int64  `json:"FromUin"`
				ToUin      int64  `json:"ToUin"`
				FromType   int    `json:"FromType"`
				SenderUin  int64  `json:"SenderUin"`
				SenderNick string `json:"SenderNick"`
				MsgType    int    `json:"MsgType"`
				C2CCmd     int    `json:"C2cCmd"`
				MsgSeq     int64  `json:"MsgSeq"`
				MsgTime    int64  `json:"MsgTime"`
				MsgRandom  int64  `json:"MsgRandom"`
				MsgUid     int64  `json:"MsgUid"`
				GroupInfo  struct {
					GroupCard    string `json:"GroupCard"`
					GroupCode    int    `json:"GroupCode"`
					GroupInfoSeq int    `json:"GroupInfoSeq"`
					GroupLevel   int    `json:"GroupLevel"`
					GroupRank    int    `json:"GroupRank"`
					GroupType    int    `json:"GroupType"`
					GroupName    string `json:"GroupName"`
				} `json:"GroupInfo"`
				C2CTempMessageHead interface{} `json:"C2CTempMessageHead"`
			} `json:"MsgHead,omitempty"`
			MsgBody *struct {
				SubMsgType int    `json:"SubMsgType"`
				Content    string `json:"Content"`
				Images     []struct {
					FileId   int64  `json:"FileId"`
					FileMd5  string `json:"FileMd5"`
					FileSize int    `json:"FileSize"`
					Url      string `json:"Url"`
				} `json:"Images"`
				AtUinLists []struct {
					QQNick string `json:"QQNick"`
					QQUid  int64  `json:"QQUid"`
				} `json:"AtUinLists"`
				Video interface{} `json:"Video"`
				Voice interface{} `json:"Voice"`
			} `json:"MsgBody,omitempty"`
		} `json:"EventData"`
		EventName string `json:"EventName"`
	} `json:"CurrentPacket"`
	CurrentQQ int64 `json:"CurrentQQ"`
}

func (e *EventStruct) ParseLoginSuccessEvent() ILoginSuccess {
	return e
}
func (e *EventStruct) GetAtList() (list []int64) {
	for _, v := range e.CurrentPacket.EventData.MsgBody.AtUinLists {
		list = append(list, v.QQUid)
	}
	return list
}
func (e *EventStruct) AtBot() (at bool) {
	for _, v := range e.CurrentPacket.EventData.MsgBody.AtUinLists {
		if v.QQUid == e.CurrentQQ {
			at = true
			break
		}
	}
	return at
}
func (e *EventStruct) ExcludeAtInfo() IGroupMsg {
	for _, v := range e.CurrentPacket.EventData.MsgBody.AtUinLists {
		e.CurrentPacket.EventData.MsgBody.Content = strings.ReplaceAll(e.CurrentPacket.EventData.MsgBody.Content, "@"+v.QQNick, "")
	}
	e.CurrentPacket.EventData.MsgBody.Content = strings.TrimSpace(e.CurrentPacket.EventData.MsgBody.Content)
	return e
}

func (e *EventStruct) ExcludeBot() IEvent {
	if e.CurrentQQ == e.CurrentPacket.EventData.MsgHead.FromUin {
		return nil
	} else {
		return e
	}
}
func (e *EventStruct) GetSenderUin() int64 {
	return e.CurrentPacket.EventData.MsgHead.SenderUin
}
func (e *EventStruct) GetApiBuilder() apiBuilder.IMainFunc {
	return apiBuilder.NewApi(e.apiUrl, e.CurrentQQ)
}
func (e *EventStruct) GetMsgType() int {
	return e.CurrentPacket.EventData.MsgHead.MsgType
}
func (e *EventStruct) GetMsgTime() int64 {
	return e.CurrentPacket.EventData.MsgHead.MsgTime
}

func (e *EventStruct) GetTextContent() string {
	return e.CurrentPacket.EventData.MsgBody.Content
}
func (e *EventStruct) ParseTextMsg() ITextMsg {
	return e
}
func (e *EventStruct) ParseGroupMsg() IGroupMsg {
	return e
}
func (e *EventStruct) GetMsgUid() int64 {
	return e.CurrentPacket.EventData.MsgHead.MsgUid
}
func (e *EventStruct) GetGroupUin() int64 {
	return e.CurrentPacket.EventData.MsgHead.FromUin
}
func (e *EventStruct) GetCurrentQQ() int64 {
	return e.CurrentQQ
}
func (e *EventStruct) GetRawBytes() []byte {
	return e.rawEvent
}
func (e *EventStruct) GetEventName() EventName {
	return EventName(e.CurrentPacket.EventName)
}
func (e *EventStruct) GetLoginSuccessBot() (nick string, uin int64) {
	nick = *e.CurrentPacket.EventData.Nick
	uin = *e.CurrentPacket.EventData.Uin
	return
}
