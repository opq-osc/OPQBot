package OPQBot

const (
	SendTypeTextMsg          = 1
	SendTypePicMsgByUrl      = 2
	SendTypePicMsgByLocal    = 3
	SendTypePicMsgByMd5      = 4
	SendTypeVoiceByUrl       = 5
	SendTypeVoiceByLocal     = 6
	SendTypeXml              = 7
	SendTypeJson             = 8
	SendTypeForword          = 9
	SendTypeReplay           = 10
	SendToTypeFriend         = 1
	SendToTypeGroup          = 2
	SendToTypePrivateChat    = 3
	EventNameOnGroupMessage  = "OnGroupMsgs"
	EventNameOnFriendMessage = "OnFriendMsgs"
	EventNameOnBotEvent      = "OnFriendMsgs"
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
	CurrentPacket CurrentPacket `json:"CurrentPacket"`
	CurrentQQ     int64         `json:"CurrentQQ"`
}
type CurrentPacket struct {
	Data struct {
		Content       string      `json:"Content"`
		FromGroupID   int         `json:"FromGroupId"`
		FromGroupName string      `json:"FromGroupName"`
		FromNickName  string      `json:"FromNickName"`
		FromUserID    int64       `json:"FromUserId"`
		MsgRandom     int         `json:"MsgRandom"`
		MsgSeq        int         `json:"MsgSeq"`
		MsgTime       int         `json:"MsgTime"`
		MsgType       string      `json:"MsgType"`
		RedBaginfo    interface{} `json:"RedBaginfo"`
	} `json:"Data"`
	WebConnID string `json:"WebConnId"`
}
