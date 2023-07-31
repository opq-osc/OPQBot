package apiBuilder

type IGroupManager interface {
	DoApi
	GetGroupLists() IGroupManager
	GetGroupMemberLists(uin int64, lastBuffer string) IGroupManager
	GroupSystemMsgAction(MsgType int, MsgSeq, GroupCode int64) IGroupSystemMsgAction
	RevokeMsg() IGroupManager
	ToGUin(Uin int64) IGroupManager
	MsgSeq(MsgSeq int64) IGroupManager
	MsgRandom(MsgRandom int64) IGroupManager
	ProhibitedUser() IGroupManager
	ToUid(Uid string) IGroupManager
	ShutTime(ShutTime int) IGroupManager
	RemoveUser() IGroupManager
	RenameUserNickName(NickName string) IGroupManager
}
type IGroupSystemMsgAction interface {
	DoApi
	OpAccept() IGroupSystemMsgAction
	OpReject() IGroupSystemMsgAction
	OpIgnore() IGroupSystemMsgAction
}

func (b *Builder) OpAccept() IGroupSystemMsgAction {
	op := 1
	b.CgiRequest.OpCode = &op
	return b
}
func (b *Builder) OpReject() IGroupSystemMsgAction {
	op := 2
	b.CgiRequest.OpCode = &op
	return b
}
func (b *Builder) OpIgnore() IGroupSystemMsgAction {
	op := 3
	b.CgiRequest.OpCode = &op
	return b
}
func (b *Builder) GroupManager() IGroupManager {
	if b.CgiRequest == nil {
		b.CgiRequest = &CgiRequest{}
	}
	return b
}

type GroupLists struct {
	GroupLists []struct {
		CreateTime int    `json:"CreateTime"`
		GroupCnt   int    `json:"GroupCnt"`
		GroupCode  int    `json:"GroupCode"`
		GroupName  string `json:"GroupName"`
		MemberCnt  int    `json:"MemberCnt"`
	} `json:"GroupLists"`
}
type GroupMemberLists struct {
	LastBuffer  string `json:"LastBuffer"`
	MemberLists []struct {
		CreditLevel   int    `json:"CreditLevel"`
		JoinTime      int    `json:"JoinTime"`
		LastSpeakTime int    `json:"LastSpeakTime"`
		Level         int    `json:"Level"`
		MemberFlag    int    `json:"MemberFlag"`
		Nick          string `json:"Nick"`
		Uid           string `json:"Uid"`
		Uin           int    `json:"Uin"`
	} `json:"MemberLists"`
}

func (b *Builder) GetGroupLists() IGroupManager {
	cmd := "GetGroupLists"
	b.CgiCmd = &cmd
	return b
}
func (b *Builder) GetGroupMemberLists(uin int64, lastBuffer string) IGroupManager {
	cmd := "GetGroupMemberLists"
	b.CgiCmd = &cmd
	b.CgiRequest.Uin = &uin
	b.CgiRequest.LastBuffer = &lastBuffer
	return b
}
func (b *Builder) GroupSystemMsgAction(MsgType int, MsgSeq, GroupCode int64) IGroupSystemMsgAction {
	b.CgiRequest.MsgType = &MsgType
	b.CgiRequest.MsgSeq = &MsgSeq
	b.CgiRequest.GroupCode = &GroupCode
	return b
}

func (b *Builder) RevokeMsg() IGroupManager {
	cmd := "GroupRevokeMsg"
	b.CgiCmd = &cmd
	return b
}

func (b *Builder) ToGUin(Uin int64) IGroupManager {
	b.CgiRequest.Uin = &Uin
	return b
}

func (b *Builder) ToUid(Uid string) IGroupManager {
	b.CgiRequest.Uid = &Uid
	return b
}

func (b *Builder) ShutTime(ShutTime int) IGroupManager {
	b.CgiRequest.BanTime = &ShutTime
	return b
}

func (b *Builder) MsgSeq(MsgSeq int64) IGroupManager {
	b.CgiRequest.MsgSeq = &MsgSeq
	return b
}

func (b *Builder) MsgRandom(MsgRandom int64) IGroupManager {
	b.CgiRequest.MsgRandom = &MsgRandom
	return b
}

func (b *Builder) ProhibitedUser() IGroupManager {
	cmd := "SsoGroup.Op"
	code := 4691
	b.CgiCmd = &cmd
	b.CgiRequest.OpCode = &code
	return b
}

func (b *Builder) RemoveUser() IGroupManager {
	cmd := "SsoGroup.Op"
	code := 2208
	b.CgiCmd = &cmd
	b.CgiRequest.OpCode = &code
	return b
}
func (b *Builder) RenameUserNickName(NickName string) IGroupManager {
	cmd := "SsoGroup.Op"
	code := 2300
	b.CgiCmd = &cmd
	b.CgiRequest.OpCode = &code
	b.CgiRequest.Nick = &NickName

	return b
}
