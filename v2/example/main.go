package main

import (
	"context"
	"github.com/charmbracelet/log"
	_ "github.com/mcoo/OPQBot/session/provider"
	"github.com/opq-osc/OPQBot/v2"
	"github.com/opq-osc/OPQBot/v2/apiBuilder"
	"github.com/opq-osc/OPQBot/v2/events"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	core, err := OPQBot.NewCore("http://localhost:8086", 10, func(builder *apiBuilder.Builder) bool {
		if *builder.CgiRequest.Content == "test" {
			content := "replaced"
			builder.CgiRequest.Content = &content
		}
		if *builder.CgiRequest.Content == "拦截" {
			log.Error("已拦截")
			return false
		}
		return true
	})
	if err != nil {
		log.Fatal(err)
	}
	core.On(events.EventNameGroupMsg, func(ctx context.Context, event events.IEvent) {
		if event.GetMsgType() == events.MsgTypeGroupMsg {
			groupMsg := event.ParseGroupMsg()
			if groupMsg.AtBot() {
				text := groupMsg.ExcludeAtInfo().ParseTextMsg().GetTextContent()
				if text == "test" {
					event.GetApiBuilder().SendMsg().GroupMsg().ToUin(groupMsg.GetGroupUin()).TextMsg("test").Do(ctx)
				}
				if text == "拦截" {
					event.GetApiBuilder().SendMsg().GroupMsg().ToUin(groupMsg.GetGroupUin()).TextMsg("拦截").Do(ctx)
				}
			}
		}

	})
	err = core.ListenAndWait(context.Background())
	if err != nil {
		panic(err)
	}
}
