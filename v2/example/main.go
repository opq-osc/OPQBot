package main

import (
	"context"
	"encoding/base64"
	"github.com/charmbracelet/log"
	_ "github.com/mcoo/OPQBot/session/provider"
	"github.com/opq-osc/OPQBot/v2"
	"github.com/opq-osc/OPQBot/v2/events"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	core, err := OPQBot.NewCore("http://localhost:8086", 10)
	if err != nil {
		log.Fatal(err)
	}
	core.On(events.EventNameGroupMsg, func(ctx context.Context, event events.IEvent) {
		groupMsg := event.ParseGroupMsg()
		if groupMsg.GetMsgType() == 82 && groupMsg.AtBot() {
			text := groupMsg.ExcludeAtInfo().ParseTextMsg().GetTextContent()
			if text == "login" {
				qr := event.GetApiBuilder().Qrcode()
				err := qr.Get()
				if err != nil {
					log.Error(err)
				}
				pic, err := event.GetApiBuilder().Upload().GroupPic().SetBase64Buf(base64.StdEncoding.EncodeToString(qr.GetImageBytes())).DoUpload(ctx)
				if err != nil {
					log.Error(err)
				}
				err = event.GetApiBuilder().SendMsg().GroupMsg().ToUin(groupMsg.GetGroupUin()).PicMsg(pic).Do(ctx)
				if err != nil {
					log.Error(err)
				}
			}

		}

	})
	err = core.ListenAndWait(context.Background())
	if err != nil {
		panic(err)
	}
}
