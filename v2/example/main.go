package main

import (
	"context"
	"github.com/charmbracelet/log"
	_ "github.com/mcoo/OPQBot/session/provider"
	"github.com/opq-osc/OPQBot/v2"
	"github.com/opq-osc/OPQBot/v2/apiBuilder"
	"github.com/opq-osc/OPQBot/v2/events"
)

const apiUrl = "http://127.0.0.1:8086"

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	core, err := OPQBot.NewCore(apiUrl, OPQBot.WithMaxRetryCount(5))
	if err != nil {
		log.Fatal(err)
	}
	// 群组相关功能
	go HandleGroupMsg(core)
	core.On(events.EventNameGroupMsg, func(ctx context.Context, event events.IEvent) {
		if event.GetMsgType() == events.MsgTypeGroupMsg {
			groupMsg := event.ParseGroupMsg()
			if groupMsg.AtBot() {
				text := groupMsg.ExcludeAtInfo().ParseTextMsg().GetTextContent()
				if text == "test" {
					apiBuilder.New(apiUrl, event.GetCurrentQQ()).SendMsg().GroupMsg().ToUin(groupMsg.GetGroupUin()).TextMsg("test").Do(ctx)
				}
			}
		}

	})

	err = core.ListenAndWait(context.Background())
	if err != nil {
		panic(err)
	}
}
