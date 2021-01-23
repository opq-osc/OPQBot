# OPQBot
OPQBot

### 功能
|功能|是否实现|
|-|-|
|群消息处理事件|是|
|好友消息处理事件|是|
|机器人事件处理|否|
### Install
`go get github.com/mcoo/OPQBot`
### Example
```golang
opqBot := OPQBot.NewBotManager("http://192.168.2.2:8899",config.Sysconfig.OPQBotUrl)
err := opqBot.Start()
if err != nil {
    log.Println(err_opq.Error())
}
err = opqBot.AddEvent(OPQBot.EventNameOnGroupMessage, func(botQQ int64, packet OPQBot.CurrentPacket) {
    log.Println(botQQ,packet)
})
if err != nil {
    log.Println(err.Error())
}
err = opqBot.AddEvent(OPQBot.EventNameOnFriendMessage, func(botQQ int64, packet OPQBot.SendTypeRelayContent) {
    log.Println(botQQ,packet)
})
if err != nil {
    log.Println(err.Error())
}
opqBot.Send(OPQBot.SendMsgPack{
    SendType:   OPQBot.SendTypePicMsgByUrl,
    SendToType: OPQBot.SendToTypeFriend,
    ToUserUid:  2435932516,
    Content:    OPQBot.SendTypePicMsgByUrlContent{Content: "你好", PicUrl: "https://img-home.csdnimg.cn/images/20201124032511.png"},
})

```

更多请看wiki