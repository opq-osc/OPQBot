package main

import (
	"context"
	"github.com/opq-osc/OPQBot/v2"
	"github.com/opq-osc/OPQBot/v2/events"
	"log"
)

func main() {
	core, err := OPQBot.NewCore("http://localhost:8086", 10)
	if err != nil {
		panic(err)
	}

	core.On(events.EventNameNewMsg, func(ctx context.Context, event events.IEvent) {
		groupMsg := event.ParseGroupMsg()
		if groupMsg.ParseTextMsg().GetTextContent() == "test" {
			//qr := event.GetApiBuilder().Qrcode()
			//qr.Get()
			//var buf bytes.Buffer
			//qr.PrintTerminal(&buf)
			pic, err := event.GetApiBuilder().Upload().SetFilePath(`F:\project\OPQBot\v2\example\img\opq.logo-only.min.png`).DoUpload(ctx)
			if err != nil {
				log.Println(pic)
				return
			}
			event.GetApiBuilder().SendMsg().GroupMsg().PicMsg("test", pic).ToUin(groupMsg.GetGroupUin()).Do(ctx)
		}
	})
	err = core.ListenAndWait(context.Background())
	if err != nil {
		panic(err)
	}
}
