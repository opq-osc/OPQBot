package main

import (
	"context"
	"fmt"
	_ "github.com/mcoo/OPQBot/session/provider"
	"github.com/opq-osc/OPQBot/v2"
	"github.com/opq-osc/OPQBot/v2/apiBuilder"
	"github.com/opq-osc/OPQBot/v2/events"
	"regexp"
	"strings"
	"sync"
)

/*
* @Author <bypanghu> (bypanghu@163.com)
* @Date 2023/7/31 14:48
**/

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
		apiBuilder.New(apiUrl, event.GetCurrentQQ()).GroupManager().GetGroupMemberLists(groupMsg.GetGroupUId(), "").Do(ctx)
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

			// 指令： 禁言@用户 1h (禁言用户和时间中间需要有一个空格)
			if IsAdmins(groupMsg.GetSenderUin(), adminUids) {
				text := groupMsg.ExcludeAtInfo().ParseTextMsg().GetTextContent()
				// 禁言用户
				if strings.HasPrefix(text, "禁言") && groupMsg.ContainedAt() {
					// 获取被禁言的用户
					user := groupMsg.GetAtInfo()
					// 获取禁言时间
					time := strings.Split(text, " ")[2]
					intTime := 0
					// 时间转换为秒
					switch time {
					case "1m":
						intTime = 60
					case "1h":
						intTime = 3600
					case "1d":
						intTime = 86400
					case "1w":
						intTime = 604800
					}
					if len(user) > 0 {
						wg := sync.WaitGroup{}
						wg.Add(len(user))
						for _, u := range user {
							go func(u events.UserInfo) {
								// 由于现在 at user 中并没有 uid ， 等待后续更新
								apiBuilder.New(apiUrl, event.GetCurrentQQ()).GroupManager().ProhibitedUser().ToGUin(groupMsg.GetGroupUin()).ToUid("u_8Unz3AyEaG3wPEEdFxPsGA").ShutTime(intTime).Do(ctx)
								apiBuilder.New(apiUrl, event.GetCurrentQQ()).SendMsg().GroupMsg().ToUin(groupMsg.GetGroupUin()).TextMsg(fmt.Sprintf("用户%s已经被成功禁言 %s", u.Nick, time)).Do(ctx)
								wg.Done()
							}(u)
						}
						wg.Wait()
					} else {
						apiBuilder.New(apiUrl, event.GetCurrentQQ()).SendMsg().GroupMsg().ToUin(groupMsg.GetGroupUin()).TextMsg("禁言失败，未找到用户").Do(ctx)
					}

				}

				if (strings.HasPrefix(text, "踢") || strings.HasPrefix(text, "移除")) && groupMsg.ContainedAt() {
					// 获取被禁言的用户
					user := groupMsg.GetAtInfo()
					if len(user) > 0 {
						wg := sync.WaitGroup{}
						wg.Add(len(user))
						for _, u := range user {
							go func(u events.UserInfo) {
								// 由于现在 at user 中并没有 uid ， 等待后续更新
								apiBuilder.New(apiUrl, event.GetCurrentQQ()).GroupManager().RemoveUser().ToGUin(groupMsg.GetGroupUin()).ToUid("u_8Unz3AyEaG3wPEEdFxPsGA").Do(ctx)
								apiBuilder.New(apiUrl, event.GetCurrentQQ()).SendMsg().GroupMsg().ToUin(groupMsg.GetGroupUin()).TextMsg(fmt.Sprintf("用户%s已经被移除本群", u.Nick)).Do(ctx)
								wg.Done()
							}(u)
						}
						wg.Wait()
					} else {
						apiBuilder.New(apiUrl, event.GetCurrentQQ()).SendMsg().GroupMsg().ToUin(groupMsg.GetGroupUin()).TextMsg("移除失败，未找到用户").Do(ctx)
					}
				}

				if strings.HasPrefix(text, "修改昵称") && groupMsg.ContainedAt() {
					// 获取被禁言的用户
					user := groupMsg.GetAtInfo()
					if len(user) > 0 {
						wg := sync.WaitGroup{}
						wg.Add(len(user))
						for _, u := range user {
							go func(u events.UserInfo) {
								// 由于现在 at user 中并没有 uid ， 等待后续更新
								apiBuilder.New(apiUrl, event.GetCurrentQQ()).GroupManager().RenameUserNickName("测试").ToGUin(groupMsg.GetGroupUin()).ToUid("u_8Unz3AyEaG3wPEEdFxPsGA").Do(ctx)
								apiBuilder.New(apiUrl, event.GetCurrentQQ()).SendMsg().GroupMsg().ToUin(groupMsg.GetGroupUin()).TextMsg(fmt.Sprintf("用户%s昵称修改完成", u.Nick)).Do(ctx)
								wg.Done()
							}(u)
						}
						wg.Wait()
					} else {
						apiBuilder.New(apiUrl, event.GetCurrentQQ()).SendMsg().GroupMsg().ToUin(groupMsg.GetGroupUin()).TextMsg("修改失败，未找到用户").Do(ctx)
					}
				}

			}

		}

	})
}
