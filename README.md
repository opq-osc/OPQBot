# OPQBot Golang SDK V2 🎉
正在开发，有想法欢迎 PR 想法 

![GitHub go.mod Go version (branch & subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/mcoo/OPQBot/v2?filename=v2/go.mod&style=for-the-badge&logo=go)
![GitHub tag (latest SemVer pre-release)](https://img.shields.io/github/v/tag/mcoo/OPQBot?include_prereleases&style=for-the-badge&logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAA1klEQVQ4T6XTvQ3CMBCG4fdKRgCxBQvQMgAZggUQFdDBCAxAAz1INGQAKCnYgR0O2VjBSc5JBCn98+Tusy2qegXG1L8HcAF2IvI05v2QqKqmJsP4C5iJyNFaFwM5sAeGwNJYnFlIDGxEZOE2quoJmHRBYmAtIqsA3IBRorVSJdUMzkAvEWrsFYgDmv7WlK9HHODKtkJrORw/nTmgD7gqBl12VNbcJYQ2BQ4/ALkH/kCyAvgB+YRYLVtVu7TzPUar7xakfJFSwSWQ2nuotRCDAZmHsa31mN4A6l46o4qtxAAAAABJRU5ErkJggg==)
![GitHub](https://img.shields.io/github/license/mcoo/OPQBot?style=for-the-badge&logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAABMklEQVQ4T5WTsS5EQRSGv18EhcgmOkqJVlQ6sZWGgkQhUW9hS4XS7jOICIUKkUi20CjxBrKdBxCFgifwy6y7nL3mbjjNnZz7zzfn/GdGhLC9B6yF1IOkVtSU1+onbO8CR8B9EK0AC5K6VZAISCcdANeFeApYBeqSInSAFQEN4CRz0iPwnsl3JB1+A5LA9l36Sqrbfis2JUA5RoBlYP8vgFRZP6bD+hhoxxa2gNMgmARGh0zgQtJOzsTYb+0/JlaNcVFSzodecWUPzoHZUPYr0JCUm8IPwHbatARslgDPQOr1duhFsn0JbANPwEsQzwDzQEtSOweR7TNgY9h9Bz6AK0nNX2/BtituWllbkzTgWc/EApDK60rq2J4r3kD6PwaMAxNFG5WAL0elBLwB1rP9Zir4BJmUbAFx6PbeAAAAAElFTkSuQmCC)

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
