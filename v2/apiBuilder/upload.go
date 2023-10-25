package apiBuilder

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/charmbracelet/log"
)

type IUpload interface {
	SetFilePath(path string) IUpload
	SetFileUrlPath(url string) IUpload
	SetBase64Buf(base64Buf string) IUpload
	FriendPic() IUpload
	FriendVoice() IUpload
	GroupVoice() IUpload
	GroupPic() IUpload
	DoUpload(ctx context.Context) (*File, error)
}

func (b *Builder) FriendVoice() IUpload {
	commandId := 26
	b.CgiRequest.CommandId = &commandId
	return b
}
func (b *Builder) GroupVoice() IUpload {
	commandId := 29
	b.CgiRequest.CommandId = &commandId
	return b
}
func (b *Builder) GroupPic() IUpload {
	commandId := 2
	b.CgiRequest.CommandId = &commandId
	return b
}
func (b *Builder) FriendPic() IUpload {
	commandId := 1
	b.CgiRequest.CommandId = &commandId
	return b
}
func (b *Builder) Upload() IUpload {
	cmd := "PicUp.DataUp"
	path := "/v1/upload"
	b.CgiCmd = &cmd
	b.path = &path
	if b.CgiRequest == nil {
		b.CgiRequest = &CgiRequest{}
	}
	return b
}

func (b *Builder) SetFilePath(path string) IUpload {
	if b.CgiRequest == nil {
		b.CgiRequest = &CgiRequest{}
	}
	b.CgiRequest.FilePath = &path
	return b
}
func (b *Builder) SetBase64Buf(base64Buf string) IUpload {
	if b.CgiRequest == nil {
		b.CgiRequest = &CgiRequest{}
	}
	b.CgiRequest.Base64Buf = &base64Buf
	return b
}
func (b *Builder) SetFileUrlPath(url string) IUpload {
	if b.CgiRequest == nil {
		b.CgiRequest = &CgiRequest{}
	}
	b.CgiRequest.FileUrl = &url
	return b
}

type File struct {
	FileMd5   string `json:"FileMd5"`
	FileSize  int    `json:"FileSize"`
	FileToken string `json:"FileToken"`
	FileId    int64  `json:"FileId"`
	Height    int    `json:"Height"`
	Width     int    `json:"Width"`
}

func (b *Builder) DoUpload(ctx context.Context) (*File, error) {
	resp, err := b.DoAndResponse(ctx)
	if err != nil {
		return nil, err
	}
	log.Debug(string(resp.GetOrigin()))
	if !resp.Ok() {
		return nil, errors.New(resp.ErrorMsg())
	}
	var pic = File{}
	err = resp.GetData(&pic)
	if err != nil {
		return nil, err
	}
	if b.CgiRequest.CommandId != nil && (*b.CgiRequest.CommandId == 1 || *b.CgiRequest.CommandId == 2) {
		var picBytes []byte = nil
		if b.CgiRequest.Base64Buf != nil {
			picBytes, err = base64.StdEncoding.DecodeString(*b.CgiRequest.Base64Buf)
			if err == nil {
				h, w := GetPicHW(picBytes)
				pic.Height = h
				pic.Width = w
			} else {
				log.Debug(err)
			}
		}

	}
	return &pic, nil
}
func GetPicHW(pic []byte) (height, width int) {
	img, _, _ := image.Decode(bytes.NewBuffer(pic))
	return img.Bounds().Dy(), img.Bounds().Dx()
}
