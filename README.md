![OPQBot](https://socialify.git.ci/opq-osc/OPQBot/image?description=1&font=Jost&forks=1&issues=1&language=1&name=1&owner=1&pattern=Plus&pulls=1&stargazers=1&theme=Light)
# OPQBot Golang SDK V2 🎉
欢迎 Star 👍 

## 安装 💡

```shell
go get -u github.com/opq-osc/OPQBot/v2@latest
```

## 使用方法

```go
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
	core.On(events.EventNameGroupMsg, func(ctx context.Context, event events.IEvent) {
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
```
