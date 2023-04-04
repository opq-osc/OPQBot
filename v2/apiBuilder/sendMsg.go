package apiBuilder

import "context"

type ISendMsg interface {
	FriendMsg() IMsg
	GroupMsg() IMsg
}

type IMsg interface {
	ToUin(uin int64) IMsg
	TextMsg(text string) IMsg
	PicMsgWithMd5(...Md5Pic) IMsg
	XmlMsg(xml string) IMsg
	JsonMsg(json string) IMsg
	At(uint ...int64) IMsg
	Do(ctx context.Context) error
	DoAndResponse(ctx context.Context) (*Response, error)
}

type Md5Pic struct {
	FileMd5 string `json:"FileMd5,omitempty"`
	Size    int64  `json:"Size,omitempty"`
}
