package OPQBot

import (
	"reflect"
)

const (
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
	SendToType   int
	ToUserUid    int64
	Content      interface{}
	CallbackFunc func(Code int, Info string, record MyRecord)
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

type SendTypePicMsgByBase64Content struct {
	Content string
	Base64  string
	Flash   bool
}
type SendTypePicMsgByBase64ContentPrivateChat struct {
	Content string
	Base64  string
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
	Content      string
}

type SendTypeForwordContentPrivateChat struct {
	ForwordBuf   string
	Content      string
	ForwordField int
	Group        int64
}

type SendTypeReplyContent struct {
	ReplayInfo struct {
		MsgSeq     int    `json:"MsgSeq"`
		MsgTime    int    `json:"MsgTime"`
		UserID     int64  `json:"UserID"`
		RawContent string `json:"RawContent"`
	} `json:"ReplayInfo"`
	Content string
}

type SendTypeReplyContentPrivateChat struct {
	ReplayInfo struct {
		MsgSeq     int    `json:"MsgSeq"`
		MsgTime    int    `json:"MsgTime"`
		UserID     int64  `json:"UserID"`
		RawContent string `json:"RawContent"`
	} `json:"ReplayInfo"`
	Content string
	Group   int64
}

type returnPack struct {
	CurrentPacket currentPacket `json:"CurrentPacket"`
	CurrentQQ     int64         `json:"CurrentQQ"`
}
type currentPacket struct {
	Data      interface{} `json:"Data"`
	WebConnID string      `json:"WebConnId"`
}
type MyRecord struct {
	FromGroupID int64  `json:"FromGroupId"`
	MsgRandom   int64  `json:"MsgRandom"`
	MsgSeq      int    `json:"MsgSeq"`
	MsgTime     int    `json:"MsgTime"`
	MsgType     string `json:"MsgType"`
	Content     string `json:"Content"`
}
type Context struct {
	Ban      bool
	NowIndex int
	MaxIndex int
	f        []reflect.Value
	Bot      *BotManager
}

func (ctx *Context) Next(currentQQ int64, result interface{}) {
	if ctx.Ban {
		return
	}
	if ctx.NowIndex >= ctx.MaxIndex {
		return
	}
	ctx.NowIndex += 1
	//r := reflect.ValueOf(result)
	//r.Field(0).Field(0).Elem().SetInt(int64(ctx.NowIndex))
	//log.Println(ctx.NowIndex)
	ctx.f[ctx.NowIndex].Call([]reflect.Value{reflect.ValueOf(currentQQ), reflect.ValueOf(result)})
}

type GroupMsgPack struct {
	Context
	Content       string      `json:"Content"`
	FromGroupID   int64       `json:"FromGroupId"`
	FromGroupName string      `json:"FromGroupName"`
	FromNickName  string      `json:"FromNickName"`
	FromUserID    int64       `json:"FromUserId"`
	MsgRandom     int64       `json:"MsgRandom"`
	MsgSeq        int         `json:"MsgSeq"`
	MsgTime       int         `json:"MsgTime"`
	MsgType       string      `json:"MsgType"`
	RedBaginfo    interface{} `json:"RedBaginfo"`
}
type FriendMsgPack struct {
	Context
	Content    string      `json:"Content"`
	FromUin    int64       `json:"FromUin"`
	MsgSeq     int         `json:"MsgSeq"`
	MsgType    string      `json:"MsgType"`
	RedBaginfo interface{} `json:"RedBaginfo"`
	ToUin      int64       `json:"ToUin"`
}
type eventPack struct {
	Context
	CurrentPacket struct {
		Data      interface{} `json:"Data"`
		WebConnID string      `json:"WebConnId"`
	} `json:"CurrentPacket"`
	CurrentQQ int64 `json:"CurrentQQ"`
}

type GroupJoinPack struct {
	Context
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
	Context
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
	Context
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
	Context
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
	Context
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
		ActionUin       int64  `json:"ActionUin"`
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
	Context
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
	Context
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
	Context
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
	Msg string `json:"Msg"`
	Ret int    `json:"Ret"`
}
type UserCardInfo struct {
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
type FriendFileResult struct {
	FileName string `json:"FileName"`
	FileSize int    `json:"FileSize"`
	FromUin  int64  `json:"FromUin"`
	URL      string `json:"Url"`
}
type GroupFileResult struct {
	Ret int    `json:"Ret"`
	URL string `json:"Url"`
}
type GroupMemberList struct {
	Count      int   `json:"Count"`
	GroupUin   int64 `json:"GroupUin"`
	LastUin    int64 `json:"LastUin"`
	MemberList []struct {
		Age           int    `json:"Age"`
		AutoRemark    string `json:"AutoRemark"`
		CreditLevel   int    `json:"CreditLevel"`
		Email         string `json:"Email"`
		FaceID        int    `json:"FaceId"`
		Gender        int    `json:"Gender"`
		GroupAdmin    int    `json:"GroupAdmin"`
		GroupCard     string `json:"GroupCard"`
		JoinTime      int    `json:"JoinTime"`
		LastSpeakTime int    `json:"LastSpeakTime"`
		MemberLevel   int    `json:"MemberLevel"`
		MemberUin     int64  `json:"MemberUin"`
		Memo          string `json:"Memo"`
		NickName      string `json:"NickName"`
		ShowName      string `json:"ShowName"`
		SpecialTitle  string `json:"SpecialTitle"`
		Status        int    `json:"Status"`
	} `json:"MemberList"`
}

// FriendList 获取好友列表表单
type FriendList struct {
	FriendCount int `json:"Friend_count"`
	Friendlist  []struct {
		FriendUin int64  `json:"FriendUin"`
		IsRemark  bool   `json:"IsRemark"`
		NickName  string `json:"NickName"`
		OnlineStr string `json:"OnlineStr"`
		Remark    string `json:"Remark"`
		Status    int    `json:"Status"`
	} `json:"Friendlist"`
	GetfriendCount    int `json:"GetfriendCount"`
	StartIndex        int `json:"StartIndex"`
	TotoalFriendCount int `json:"Totoal_friend_count"`
}

// GroupList 获取群列表表单
type GroupList struct {
	Count     int    `json:"Count"`
	NextToken string `json:"NextToken"`
	TroopList []struct {
		GroupID          int64  `json:"GroupId"`
		GroupMemberCount int64  `json:"GroupMemberCount"`
		GroupName        string `json:"GroupName"`
		GroupNotice      string `json:"GroupNotice"`
		GroupOwner       int64  `json:"GroupOwner"`
		GroupTotalCount  int    `json:"GroupTotalCount"`
	} `json:"TroopList"`
}

// UserInfo 用户信息表单
type UserInfo struct {
	Code int `json:"code"`
	Data struct {
		Astro         int    `json:"astro"`
		AvatarURL     string `json:"avatarUrl"`
		Bitmap        string `json:"bitmap"`
		Bluevip       int    `json:"bluevip"`
		Commfrd       int    `json:"commfrd"`
		Friendship    int    `json:"friendship"`
		From          string `json:"from"`
		Gender        int    `json:"gender"`
		Greenvip      int    `json:"greenvip"`
		IntimacyScore int    `json:"intimacyScore"`
		IsFriend      int    `json:"isFriend"`
		Logolabel     string `json:"logolabel"`
		Nickname      string `json:"nickname"`
		Publicwalfare int    `json:"publicwalfare"`
		Qqvip         int    `json:"qqvip"`
		Qzone         int    `json:"qzone"`
		Realname      string `json:"realname"`
		Smartname     string `json:"smartname"`
		Uin           int64  `json:"uin"`
	} `json:"data"`
	Default int    `json:"default"`
	Message string `json:"message"`
	Subcode int    `json:"subcode"`
}

type Cookie struct {
	ClientKey string `json:"ClientKey"`
	Cookies   string `json:"Cookies"`
	Gtk       string `json:"Gtk"`
	Gtk32     string `json:"Gtk32"`
	PSkey     struct {
		Connect     string `json:"connect"`
		Docs        string `json:"docs"`
		Docx        string `json:"docx"`
		Game        string `json:"game"`
		Gamecenter  string `json:"gamecenter"`
		Imgcache    string `json:"imgcache"`
		MTencentCom string `json:"m.tencent.com"`
		Mail        string `json:"mail"`
		Mma         string `json:"mma"`
		Now         string `json:"now"`
		Office      string `json:"office"`
		Openmobile  string `json:"openmobile"`
		Qqweb       string `json:"qqweb"`
		Qun         string `json:"qun"`
		Qzone       string `json:"qzone"`
		QzoneCom    string `json:"qzone.com"`
		TenpayCom   string `json:"tenpay.com"`
		Ti          string `json:"ti"`
		Vip         string `json:"vip"`
		Weishi      string `json:"weishi"`
	} `json:"PSkey"`
	Skey string `json:"Skey"`
}
type AtMsg struct {
	Content string `json:"Content"`
	UserExt []struct {
		QQNick string `json:"QQNick"`
		QQUID  int64  `json:"QQUid"`
	} `json:"UserExt"`
	UserID []int64 `json:"UserID"`
}
type Reply struct {
	Content    string  `json:"Content"`
	SrcContent string  `json:"SrcContent"`
	MsgSeq     int     `json:"MsgSeq"`
	Tips       string  `json:"Tips"`
	UserID     []int64 `json:"UserID"`
}

type PicMsg struct {
	Content  string `json:"Content"`
	GroupPic []struct {
		FileId       int64  `json:"FileId"`
		FileMd5      string `json:"FileMd5"`
		FileSize     int    `json:"FileSize"`
		ForwordBuf   string `json:"ForwordBuf"`
		ForwordField int    `json:"ForwordField"`
		Url          string `json:"Url"`
	} `json:"GroupPic"`
	Tips    string `json:"Tips"`
	UserExt []struct {
		QQNick string `json:"QQNick"`
		QQUid  int64  `json:"QQUid"`
	} `json:"UserExt"`
	UserID []int64 `json:"UserID"`
}

type GroupFileMsg struct {
	FileID   string `json:"FileID"`
	FileName string `json:"FileName"`
	FileSize int    `json:"FileSize"`
	Tips     string `json:"Tips"`
}

type VideoMsg struct {
	ForwordBuf   string `json:"ForwordBuf"`
	ForwordField int    `json:"ForwordField"`
	VideoMd5     string `json:"VideoMd5"`
	VideoSize    string `json:"VideoSize"`
	VideoUrl     string `json:"VideoUrl"`
	Tips         string `json:"Tips"`
}
