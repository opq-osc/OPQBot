package events

import (
	"context"
	"encoding/json"
	"github.com/opq-osc/OPQBot/v2/apiBuilder"
)

type EventName string

const (
	EventNameNewMsg      EventName = "ON_EVENT_QQNT_NEW_MSG"
	EventNameGroupMsg    EventName = "ON_EVENT_GROUP_NEW_MSG"
	EventNameFriendMsg   EventName = "ON_EVENT_FRIEND_NEW_MSG"
	EventNameGroupJoin   EventName = "ON_EVENT_GROUP_JOIN"
	EventNameGroupExit   EventName = "ON_EVENT_GROUP_EXIT"
	EventNameGroupInvite EventName = "ON_EVENT_GROUP_EXIT"
)

type EventCallbackFunc func(ctx context.Context, event IEvent)

type IEvent interface {
	GetCurrentQQ() int64
	GetRawBytes() []byte
	GetEventName() EventName
	ParseGroupMsg() IGroupMsg
	GetApiBuilder() apiBuilder.IMainFunc
}
type IGroupMsg interface {
	ICommonMsg
	GetGroupUin() int64
	ParseTextMsg() ITextMsg
}
type ITextMsg interface {
	GetTextContent() string
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
			MsgHead struct {
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
			} `json:"MsgHead"`
			MsgBody struct {
				SubMsgType int         `json:"SubMsgType"`
				Content    string      `json:"Content"`
				AtUinLists interface{} `json:"AtUinLists"`
				Images     []struct {
					FileId   int64  `json:"FileId"`
					FileMd5  string `json:"FileMd5"`
					FileSize int    `json:"FileSize"`
					Url      string `json:"Url"`
				} `json:"Images"`
				Video interface{} `json:"Video"`
				Voice interface{} `json:"Voice"`
			} `json:"MsgBody"`
		} `json:"EventData"`
		EventName string `json:"EventName"`
	} `json:"CurrentPacket"`
	CurrentQQ int64 `json:"CurrentQQ"`
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
