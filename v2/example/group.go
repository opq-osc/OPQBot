package main

import (
	"context"
	_ "github.com/mcoo/OPQBot/session/provider"
	"github.com/opq-osc/OPQBot/v2"
	"github.com/opq-osc/OPQBot/v2/apiBuilder"
	"github.com/opq-osc/OPQBot/v2/events"
	"regexp"
	"sync"
)

/*
* @Author <bypanghu> (bypanghu@163.com)
* @Date 2023/7/31 14:48
**/

var (
	groupUids = []int64{123455}
	adminUids = []int64{123455}
)

func IsInGroupS(str int64, groupUids []int64) bool {
	for _, s := range groupUids {
		if s == str {
			return true
		}
	}
	return false
}

func IsAdmins(uid int64, adminUids []int64) bool {
	for _, s := range adminUids {
		if s == uid {
			return true
		}
	}
	return false
}

func ContainsURL(str string) bool {
	// 定义匹配网址的正则表达式模式
	pattern := `(?i)(https?:\/\/)?([\w-]+\.)?([\w-]+\.[\w-]+)`
	// 编译正则表达式模式
	regex := regexp.MustCompile(pattern)
	// 使用正则表达式进行匹配
	return regex.MatchString(str)
}

func HandleGroupMsg(core *OPQBot.Core) {
	core.On(events.EventNameGroupJoin, func(ctx context.Context, event events.IEvent) {
		groupMsg := event.PraseGroupJoinEvent()
		apiBuilder.New(apiUrl, event.GetCurrentQQ()).SendMsg().GroupMsg().ToUin(groupMsg.GetGroupUId()).TextMsg("欢迎新人~").Do(ctx)
	})

	core.On(events.EventNameGroupMsg, func(ctx context.Context, event events.IEvent) {
		if event.GetMsgType() == events.MsgTypeGroupMsg {
			groupMsg := event.ParseGroupMsg()
			// 监控的群聊里面禁止所有非管理员发送网址信息
			if IsInGroupS(groupMsg.GetGroupUin(), groupUids) {
				text := groupMsg.ExcludeAtInfo().ParseTextMsg().GetTextContent()
				if IsAdmins(groupMsg.GetSenderUin(), adminUids) {
					if ContainsURL(text) {
						wg := sync.WaitGroup{}
						wg.Add(2)
						go func() {
							apiBuilder.New(apiUrl, event.GetCurrentQQ()).SendMsg().GroupMsg().ToUin(groupMsg.GetGroupUin()).TextMsg("发送内容中存在违规网址，系统已撤回，警告一次！").Do(ctx)
							wg.Done()
						}()
						go func() {
							apiBuilder.New(apiUrl, event.GetCurrentQQ()).GroupManager().RevokeMsg().ToGUin(groupMsg.GetGroupUin()).MsgSeq(groupMsg.GetMsgSeq()).MsgRandom(groupMsg.GetMsgRandom()).Do(ctx)
							wg.Done()
						}()

						wg.Wait()
					}
				}
			}
		}

	})
}
