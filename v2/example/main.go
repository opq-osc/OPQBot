package main

import (
	"context"
	"github.com/opq-osc/OPQBot/v2"
	"github.com/opq-osc/OPQBot/v2/events"
)

func main() {
	core, err := OPQBot.NewCore("http://localhost:8086")
	if err != nil {
		panic(err)
	}
	core.On(events.EventNewMsg, func(ctx context.Context, event events.IEvent) {
		apiBuilder := event.GetApiBuilder()
		groupMsg := event.ParseGroupMsg()
		if groupMsg.ParseTextMsg().GetTextContent() == "hello" {
			apiBuilder.SendMsg().GroupMsg().TextMsg("你好").ToUin(groupMsg.GetGroupUin()).Do(ctx)
		}
	})
	err = core.ListenAndWait(context.Background())
	if err != nil {
		panic(err)
	}
}
