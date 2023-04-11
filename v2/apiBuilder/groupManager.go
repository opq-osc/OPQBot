package apiBuilder

type IGroupManager interface {
	DoApi
	GetGroupLists() IGroupManager
	GroupSystemMsgAction(MsgType int, MsgSeq, GroupCode int64) IGroupSystemMsgAction
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
	} `json:"MembrLists"`
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
