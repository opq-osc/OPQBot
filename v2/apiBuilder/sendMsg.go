package apiBuilder

type ISendMsg interface {
	FriendMsg() IMsg
	GroupMsg() IMsg
}

type IMsg interface {
	ToUin(uin int64) IMsg
	TextMsg(text string) IMsg
	PicMsg(...*File) IMsg
	XmlMsg(xml string) IMsg
	JsonMsg(json string) IMsg
	At(uint ...int64) IMsg
	DoApi
}

func (b *Builder) SendMsg() ISendMsg {
	cmd := "MessageSvc.PbSendMsg"
	b.CgiCmd = &cmd
	if b.CgiRequest == nil {
		b.CgiRequest = &CgiRequest{}
	}
	return b
}

func (b *Builder) FriendMsg() IMsg {
	if b.CgiRequest == nil {
		b.CgiRequest = &CgiRequest{}
	}
	toType := 1
	b.CgiRequest.ToType = &toType
	return b
}

func (b *Builder) GroupMsg() IMsg {
	if b.CgiRequest == nil {
		b.CgiRequest = &CgiRequest{}
	}
	toType := 2
	b.CgiRequest.ToType = &toType
	return b
}

func (b *Builder) ToUin(uin int64) IMsg {
	if b.CgiRequest == nil {
		b.CgiRequest = &CgiRequest{}
	}
	b.CgiRequest.ToUin = &uin
	return b
}

func (b *Builder) TextMsg(text string) IMsg {
	if b.CgiRequest == nil {
		b.CgiRequest = &CgiRequest{}
	}
	b.CgiRequest.Content = &text
	return b
}
func (b *Builder) PicMsg(pics ...*File) IMsg {
	if b.CgiRequest == nil {
		b.CgiRequest = &CgiRequest{}
	}
	b.CgiRequest.Images = append(b.CgiRequest.Images, pics...)
	return b
}
func (b *Builder) XmlMsg(xml string) IMsg {
	if b.CgiRequest == nil {
		b.CgiRequest = &CgiRequest{}
	}
	subMsgType := 12
	b.CgiRequest.SubMsgType = &subMsgType
	b.CgiRequest.Content = &xml
	return b
}
func (b *Builder) JsonMsg(json string) IMsg {
	if b.CgiRequest == nil {
		b.CgiRequest = &CgiRequest{}
	}
	subMsgType := 51
	b.CgiRequest.SubMsgType = &subMsgType
	b.CgiRequest.Content = &json
	return b
}
func (b *Builder) At(uin ...int64) IMsg {
	if b.CgiRequest == nil {
		b.CgiRequest = &CgiRequest{}
	}
	for _, v := range uin {
		qq := struct {
			Uin *int64 `json:"Uin,omitempty"`
		}{Uin: &v}
		b.CgiRequest.AtUinLists = append(b.CgiRequest.AtUinLists, qq)
	}
	return b
}
