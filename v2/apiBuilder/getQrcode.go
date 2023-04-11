package apiBuilder

import (
	"bytes"
	"encoding/base64"
	"github.com/imroc/req/v3"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/mdp/qrterminal/v3"
	"github.com/opq-osc/OPQBot/v2/errors"
	"image"
	_ "image/png"
	"io"
	"regexp"
)

func (b *Builder) Qrcode() IQrcode {
	return &QrcodeStruct{url: b.url}
}

type IQrcode interface {
	Get() error
	GetUrl() *string
	PrintTerminal(io.Writer) error
	GetImageBytes() []byte
}

type QrcodeStruct struct {
	url      string
	loginUrl *string
	imgByte  []byte
}

func (q *QrcodeStruct) GetImageBytes() []byte {
	return q.imgByte
}
func (q *QrcodeStruct) GetUrl() *string {
	return q.loginUrl
}

func (q *QrcodeStruct) PrintTerminal(writer io.Writer) error {
	if q.loginUrl == nil {
		return errors.ErrorData
	}
	qrterminal.GenerateWithConfig(*q.loginUrl, qrterminal.Config{
		Level:     qrterminal.L,
		Writer:    writer,
		BlackChar: qrterminal.BLACK,
		WhiteChar: qrterminal.WHITE,
		QuietZone: 0,
	})
	return nil
}

func (q *QrcodeStruct) Get() error {
	resp, err := req.R().Get(q.url + "/v1/login/getqrcode")
	if err != nil {
		return err
	}
	rgx := regexp.MustCompile(`data:image/png;base64,(.*?)"`)
	data := rgx.FindSubmatch(resp.Bytes())
	if len(data) != 2 {
		return errors.ErrorData
	}
	img := string(data[1])
	imgByte, err := base64.StdEncoding.DecodeString(img)
	if err != nil {
		return err
	}
	q.imgByte = imgByte
	buf := bytes.NewBuffer(imgByte)
	png, _, err := image.Decode(buf)
	if err != nil {
		return err
	}
	p, err := gozxing.NewBinaryBitmapFromImage(png)
	if err != nil {
		return err
	}
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(p, nil)
	if err != nil {
		return err
	}
	loginUrl := result.String()
	q.loginUrl = &loginUrl
	return nil
}
