package apiBuilder

import (
	"context"
	"errors"
)

type IUpload interface {
	SetFilePath(path string) IUpload
	SetFileUrlPath(url string) IUpload
	DoUpload(ctx context.Context) (*File, error)
}

func (b *Builder) Upload() IUpload {
	cmd := "PicUp.DataUp"
	path := "/v1/upload"
	b.CgiCmd = &cmd
	b.path = &path
	return b
}

func (b *Builder) SetFilePath(path string) IUpload {
	if b.CgiRequest == nil {
		b.CgiRequest = &CgiRequest{}
	}
	CommandId := 1
	b.CgiRequest.CommandId = &CommandId
	b.CgiRequest.FilePath = &path
	return b
}
func (b *Builder) SetFileUrlPath(url string) IUpload {
	if b.CgiRequest == nil {
		b.CgiRequest = &CgiRequest{}
	}
	CommandId := 1
	b.CgiRequest.CommandId = &CommandId
	b.CgiRequest.FileUrl = &url
	return b
}

type File struct {
	FileMd5   string `json:"FileMd5"`
	FileSize  int    `json:"FileSize"`
	FileToken string `json:"FileToken"`
}

func (b *Builder) DoUpload(ctx context.Context) (*File, error) {
	resp, err := b.DoAndResponse(ctx)
	if err != nil {
		return nil, err
	}
	if !resp.Ok() {
		return nil, errors.New(resp.ErrorMsg())
	}
	var pic = File{}
	err = resp.GetData(&pic)
	if err != nil {
		return nil, err
	}
	return &pic, nil
}
