package OPQBot

const (
	SendTypeTextMsg                = 1
	SendTypePicMsgByUrl            = 2
	SendTypePicMsgByLocal          = 3
	SendTypePicMsgByMd5            = 4
	SendTypeVoiceByUrl             = 5
	SendTypeVoiceByLocal           = 6
	SendTypeXml                    = 7
	SendTypeJson                   = 8
	SendTypeForword                = 9
	SendTypeReplay                 = 10
	SendToTypeFriend               = 1
	SendToTypeGroup                = 2
	SendToTypePrivateChat          = 3
	EventNameOnGroupMessage        = "OnGroupMsgs"
	EventNameOnFriendMessage       = "OnFriendMsgs"
	EventNameOnBotEvent            = "OnFriendMsgs"
	EventNameOnGroupJoin           = "ON_EVENT_GROUP_JOIN"
	EventNameOnGroupAdmin          = "ON_EVENT_GROUP_ADMIN"
	EventNameOnGroupExit           = "ON_EVENT_GROUP_EXIT"
	EventNameOnGroupExitSuccess    = "ON_EVENT_GROUP_EXIT_SUCC"
	EventNameOnGroupAdminSysNotify = "ON_EVENT_GROUP_ADMINSYSNOTIFY"
	EventNameOnGroupRevoke         = "ON_EVENT_GROUP_REVOKE"
	EventNameOnGroupShut           = "ON_EVENT_GROUP_SHUT"
	EventNameOnGroupSystemNotify   = "ON_EVENT_GROUP_SYSTEMNOTIFY"
	EventNameOnConnected           = "connection"
	EventNameOnDisconnected        = "disconnection"
	EventNameOnOther               = "other"
)

type SendMsgPack struct {
	SendType   int
	SendToType int
	ToUserUid  int64
	Content    interface{}
}

type SendTypeTextMsgContent struct {
	Content string
}

type SendTypeTextMsgContentPrivateChat struct {
	Content string
	Group   int64
}

type SendTypePicMsgByUrlContent struct {
	Content string
	PicUrl  string
	Flash   bool
}

type SendTypePicMsgByUrlContentPrivateChat struct {
	Content string
	PicUrl  string
	Group   int64
	Flash   bool
}

type SendTypePicMsgByLocalContent struct {
	Content string
	Path    string
	Flash   bool
}

type SendTypePicMsgByLocalContentPrivateChat struct {
	Content string
	Path    string
	Group   int64
	Flash   bool
}

type SendTypePicMsgByMd5Content struct {
	Content string
	Md5     string
	Flash   bool
}

type SendTypePicMsgByMd5ContentPrivateChat struct {
	Content string
	Md5s    []string
	Group   int64
	Flash   bool
}

type SendTypeVoiceByUrlContent struct {
	VoiceUrl string
}

type SendTypeVoiceByUrlContentPrivateChat struct {
	VoiceUrl string
	Group    int64
}

type SendTypeVoiceByLocalContent struct {
	Path string
}

type SendTypeVoiceByLocalContentPrivateChat struct {
	Path  string
	Group int64
}

type SendTypeXmlContent struct {
	Content string
}

type SendTypeXmlContentPrivateChat struct {
	Content string
	Group   int64
}

type SendTypeJsonContent struct {
	Content string
}

type SendTypeJsonContentPrivateChat struct {
	Content string
	Group   int64
}

type SendTypeForwordContent struct {
	ForwordBuf   string
	ForwordField int
}

type SendTypeForwordContentPrivateChat struct {
	ForwordBuf   string
	ForwordField int
	Group        int64
}

type SendTypeRelayContent struct {
	ReplayInfo interface{}
}

type SendTypeRelayContentPrivateChat struct {
	ReplayInfo interface{}
	Group      int64
}

type returnPack struct {
	CurrentPacket currentPacket `json:"CurrentPacket"`
	CurrentQQ     int64         `json:"CurrentQQ"`
}
type currentPacket struct {
	Data      interface{} `json:"Data"`
	WebConnID string      `json:"WebConnId"`
}
type GroupMsgPack struct {
	Content       string      `json:"Content"`
	FromGroupID   int64       `json:"FromGroupId"`
	FromGroupName string      `json:"FromGroupName"`
	FromNickName  string      `json:"FromNickName"`
	FromUserID    int64       `json:"FromUserId"`
	MsgRandom     int         `json:"MsgRandom"`
	MsgSeq        int         `json:"MsgSeq"`
	MsgTime       int         `json:"MsgTime"`
	MsgType       string      `json:"MsgType"`
	RedBaginfo    interface{} `json:"RedBaginfo"`
}
type FriendMsgPack struct {
	Content    string      `json:"Content"`
	FromUin    int64       `json:"FromUin"`
	MsgSeq     int         `json:"MsgSeq"`
	MsgType    string      `json:"MsgType"`
	RedBaginfo interface{} `json:"RedBaginfo"`
	ToUin      int64       `json:"ToUin"`
}
type eventPack struct {
	CurrentPacket struct {
		Data      interface{} `json:"Data"`
		WebConnID string      `json:"WebConnId"`
	} `json:"CurrentPacket"`
	CurrentQQ int64 `json:"CurrentQQ"`
}

type GroupJoinPack struct {
	EventData struct {
		InviteUin int64  `json:"InviteUin"`
		UserID    int64  `json:"UserID"`
		UserName  string `json:"UserName"`
	} `json:"EventData"`
	EventMsg struct {
		FromUin    int64       `json:"FromUin"`
		ToUin      int64       `json:"ToUin"`
		MsgType    string      `json:"MsgType"`
		MsgSeq     int         `json:"MsgSeq"`
		Content    string      `json:"Content"`
		RedBaginfo interface{} `json:"RedBaginfo"`
	} `json:"EventMsg"`
}

type GroupAdminPack struct {
	EventData struct {
		Flag    int   `json:"Flag"`
		GroupID int64 `json:"GroupID"`
		UserID  int64 `json:"UserID"`
	} `json:"EventData"`
	EventMsg struct {
		FromUin    int64       `json:"FromUin"`
		ToUin      int64       `json:"ToUin"`
		MsgType    string      `json:"MsgType"`
		MsgSeq     int         `json:"MsgSeq"`
		Content    string      `json:"Content"`
		RedBaginfo interface{} `json:"RedBaginfo"`
	} `json:"EventMsg"`
}

type GroupExitPack struct {
	EventData struct {
		UserID int64 `json:"UserID"`
	} `json:"EventData"`
	EventMsg struct {
		FromUin    int64       `json:"FromUin"`
		ToUin      int64       `json:"ToUin"`
		MsgType    string      `json:"MsgType"`
		MsgSeq     int         `json:"MsgSeq"`
		Content    string      `json:"Content"`
		RedBaginfo interface{} `json:"RedBaginfo"`
	} `json:"EventMsg"`
}
type GroupExitSuccessPack struct {
	EventData struct {
		GroupID int64 `json:"GroupID"`
	} `json:"EventData"`
	EventMsg struct {
		FromUin    int64       `json:"FromUin"`
		ToUin      int64       `json:"ToUin"`
		MsgType    string      `json:"MsgType"`
		MsgSeq     int         `json:"MsgSeq"`
		Content    string      `json:"Content"`
		RedBaginfo interface{} `json:"RedBaginfo"`
	} `json:"EventMsg"`
}

type GroupAdminSysNotifyPack struct {
	EventData struct {
		Seq             int64  `json:"Seq"`
		Type            int    `json:"Type"`
		MsgTypeStr      string `json:"MsgTypeStr"`
		Who             int    `json:"Who"`
		WhoName         string `json:"WhoName"`
		MsgStatusStr    string `json:"MsgStatusStr"`
		Content         string `json:"Content"`
		RefuseContent   string `json:"RefuseContent"`
		Flag7           int    `json:"Flag_7"`
		Flag8           int    `json:"Flag_8"`
		GroupID         int64  `json:"GroupId"`
		GroupName       string `json:"GroupName"`
		ActionUin       int    `json:"ActionUin"`
		ActionName      string `json:"ActionName"`
		ActionGroupCard string `json:"ActionGroupCard"`
		Action          int    `json:"Action"`
	} `json:"EventData"`
	EventMsg struct {
		FromUin    int64       `json:"FromUin"`
		ToUin      int64       `json:"ToUin"`
		MsgType    string      `json:"MsgType"`
		MsgSeq     int         `json:"MsgSeq"`
		Content    string      `json:"Content"`
		RedBaginfo interface{} `json:"RedBaginfo"`
	} `json:"EventMsg"`
}

type GroupRevokePack struct {
	EventData struct {
		AdminUserID int   `json:"AdminUserID"`
		GroupID     int64 `json:"GroupID"`
		MsgRandom   int64 `json:"MsgRandom"`
		MsgSeq      int   `json:"MsgSeq"`
		UserID      int64 `json:"UserID"`
	} `json:"EventData"`
	EventMsg struct {
		FromUin    int64       `json:"FromUin"`
		ToUin      int64       `json:"ToUin"`
		MsgType    string      `json:"MsgType"`
		MsgSeq     int         `json:"MsgSeq"`
		Content    string      `json:"Content"`
		RedBaginfo interface{} `json:"RedBaginfo"`
	} `json:"EventMsg"`
}
type GroupShutPack struct {
	EventData struct {
		GroupID  int64 `json:"GroupID"`
		ShutTime int   `json:"ShutTime"`
		UserID   int64 `json:"UserID"`
	} `json:"EventData"`
	EventMsg struct {
		FromUin    int64       `json:"FromUin"`
		ToUin      int64       `json:"ToUin"`
		MsgType    string      `json:"MsgType"`
		MsgSeq     int         `json:"MsgSeq"`
		Content    string      `json:"Content"`
		RedBaginfo interface{} `json:"RedBaginfo"`
	} `json:"EventMsg"`
}
type GroupSystemNotifyPack struct {
	EventData struct {
		Content string `json:"Content"`
		GroupID int64  `json:"GroupID"`
		UserID  int64  `json:"UserID"`
	} `json:"EventData"`
	EventMsg struct {
		FromUin    int64       `json:"FromUin"`
		ToUin      int64       `json:"ToUin"`
		MsgType    string      `json:"MsgType"`
		MsgSeq     int         `json:"MsgSeq"`
		Content    string      `json:"Content"`
		RedBaginfo interface{} `json:"RedBaginfo"`
	} `json:"EventMsg"`
}

type Result struct {
	Ret int `json:"Ret"`
}
type UserInfo struct {
	Age       int    `json:"Age"`
	City      string `json:"City"`
	LikeNums  int    `json:"LikeNums"`
	LoginDays int    `json:"LoginDays"`
	NickName  string `json:"NickName"`
	Province  string `json:"Province"`
	QQLevel   int    `json:"QQLevel"`
	QQUin     int64  `json:"QQUin"`
	Sex       int    `json:"Sex"`
}
