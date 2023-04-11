package apiBuilder

type IFriendManager interface {
	DoApi
	GetFriendLists(LastUin int64) IFriendManager
}
type FriendLists struct {
	FriendLists []struct {
		Head      string `json:"Head"`
		Signature string `json:"Signature"`
		TagId     int    `json:"TagId"`
		Uid       string `json:"Uid"`
		Uin       int    `json:"Uin"`
		Nick      string `json:"Nick"`
		Sex       int    `json:"Sex"`
	} `json:"FriendLists"`
	LastBuffer string `json:"LastBuffer"`
	TagLists   []struct {
		TagId   int    `json:"TagId"`
		TagName string `json:"TagName"`
	} `json:"TagLists"`
}

func (b *Builder) FriendManager() IFriendManager {
	if b.CgiRequest == nil {
		b.CgiRequest = &CgiRequest{}
	}
	return b
}
func (b *Builder) GetFriendLists(LastUin int64) IFriendManager {
	b.CgiRequest.LastUin = &LastUin
	return b
}
