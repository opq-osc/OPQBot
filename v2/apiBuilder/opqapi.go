package apiBuilder

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"strconv"
)

type Builder struct {
	qqBot      int64
	url        string
	CgiCmd     *string     `json:"CgiCmd,omitempty"`
	CgiRequest *CgiRequest `json:"CgiRequest,omitempty"`
}
type CgiRequest struct {
	ToUin      *int64   `json:"ToUin,omitempty"`
	ToType     *int     `json:"ToType,omitempty"`
	Content    *string  `json:"Content,omitempty"`
	SubMsgType *int     `json:"SubMsgType,omitempty"`
	Images     []Md5Pic `json:"Images,omitempty"`
	AtUinLists []struct {
		QQUid *int64 `json:"QQUid,omitempty"`
	} `json:"AtUinLists,omitempty"`
}

func (b *Builder) SendMsg() ISendMsg {
	cmd := "MessageSvc.PbSendMsg"
	b.CgiCmd = &cmd
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
func (b *Builder) PicMsgWithMd5(pics ...Md5Pic) IMsg {
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
			QQUid *int64 `json:"QQUid,omitempty"`
		}{QQUid: &v}
		b.CgiRequest.AtUinLists = append(b.CgiRequest.AtUinLists, qq)
	}
	return b
}
func (b *Builder) BuildStringBody() (string, error) {
	body, err := json.Marshal(b)
	return string(body), err
}

func (b *Builder) Do(ctx context.Context) error {
	r, err := b.DoAndResponse(ctx)
	if err != nil {
		return err
	}
	if !r.Ok() {
		return fmt.Errorf(r.ErrorMsg())
	}
	return nil
}
func (b *Builder) DoAndResponse(ctx context.Context) (*Response, error) {
	body, err := b.BuildStringBody()
	if err != nil {
		return nil, err
	}
	resp, err := req.SetContext(ctx).SetQueryParam("funcname", "MagicCgiCmd").SetQueryParam("qq", strconv.FormatInt(b.qqBot, 10)).SetBodyJsonString(body).Post(b.url)
	if err != nil {
		return nil, err
	}
	r := NewResponse(resp.Bytes())
	if !r.Ok() {
		return nil, fmt.Errorf(r.ErrorMsg())
	}
	return r, nil
}

func NewApi(url string, botQQ int64) IMainFunc {
	return &Builder{qqBot: botQQ, url: url}
}
