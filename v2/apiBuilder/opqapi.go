package apiBuilder

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/imroc/req/v3"
	"net/url"
	"strconv"
)

type DoApi interface {
	Do(ctx context.Context) error
	DoAndResponse(ctx context.Context) (*Response, error)
}
type Builder struct {
	qqBot      int64
	url        string
	path       *string
	method     *string
	CgiCmd     *string     `json:"CgiCmd,omitempty"`
	CgiRequest *CgiRequest `json:"CgiRequest,omitempty"`
}
type CgiRequest struct {
	LastUin    *int64  `json:"LastUin,omitempty"`
	OpCode     *int    `json:"OpCode,omitempty"`
	MsgSeq     *int64  `json:"MsgSeq,omitempty"`
	MsgType    *int    `json:"MsgType,omitempty"`
	GroupCode  *int64  `json:"GroupCode,omitempty"`
	Uin        *int64  `json:"Uin,omitempty"`
	LastBuffer *string `json:"LastBuffer,omitempty"`
	CommandId  *int    `json:"CommandId,omitempty"`
	FilePath   *string `json:"FilePath,omitempty"`
	Base64Buf  *string `json:"Base64Buf,omitempty"`
	FileUrl    *string `json:"FileUrl,omitempty"`
	ToUin      *int64  `json:"ToUin,omitempty"`
	ToType     *int    `json:"ToType,omitempty"`
	Content    *string `json:"Content,omitempty"`
	SubMsgType *int    `json:"SubMsgType,omitempty"`
	Images     []*File `json:"Images,omitempty"`
	Uid        *string `json:"Uid,omitempty"`
	AtUinLists []struct {
		Uin *int64 `json:"Uin,omitempty"`
	} `json:"AtUinLists,omitempty"`
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
	log.Debug("request", "body", body)
	client := req.SetContext(ctx)
	if b.path != nil {
		u, _ := url.JoinPath(b.url, *b.path)
		client.SetURL(u)
	} else {
		u, _ := url.JoinPath(b.url, "/v1/LuaApiCaller")
		client.SetURL(u)
	}
	if b.method != nil {
		client.Method = *b.method
	} else {
		client.Method = "POST"
	}

	resp := client.SetQueryParam("funcname", "MagicCgiCmd").SetQueryParam("qq", strconv.FormatInt(b.qqBot, 10)).SetBodyJsonString(body).Do()
	if resp.Err != nil {
		return nil, resp.Err
	}
	r := NewResponse(resp.Bytes())
	if !r.Ok() {
		return nil, fmt.Errorf(r.ErrorMsg())
	}
	return r, nil
}

func New(url string, botQQ int64) IMainFunc {
	return &Builder{qqBot: botQQ, url: url}
}
