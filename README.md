# OPQBot 🎉
![GitHub go.mod Go version (branch & subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/mcoo/OPQBot/master?filename=go.mod&style=for-the-badge&logo=go) ![GitHub tag (latest SemVer pre-release)](https://img.shields.io/github/v/tag/mcoo/OPQBot?include_prereleases&style=for-the-badge&logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAA1klEQVQ4T6XTvQ3CMBCG4fdKRgCxBQvQMgAZggUQFdDBCAxAAz1INGQAKCnYgR0O2VjBSc5JBCn98+Tusy2qegXG1L8HcAF2IvI05v2QqKqmJsP4C5iJyNFaFwM5sAeGwNJYnFlIDGxEZOE2quoJmHRBYmAtIqsA3IBRorVSJdUMzkAvEWrsFYgDmv7WlK9HHODKtkJrORw/nTmgD7gqBl12VNbcJYQ2BQ4/ALkH/kCyAvgB+YRYLVtVu7TzPUar7xakfJFSwSWQ2nuotRCDAZmHsa31mN4A6l46o4qtxAAAAABJRU5ErkJggg==) ![GitHub](https://img.shields.io/github/license/mcoo/OPQBot?style=for-the-badge&logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAABMklEQVQ4T5WTsS5EQRSGv18EhcgmOkqJVlQ6sZWGgkQhUW9hS4XS7jOICIUKkUi20CjxBrKdBxCFgifwy6y7nL3mbjjNnZz7zzfn/GdGhLC9B6yF1IOkVtSU1+onbO8CR8B9EK0AC5K6VZAISCcdANeFeApYBeqSInSAFQEN4CRz0iPwnsl3JB1+A5LA9l36Sqrbfis2JUA5RoBlYP8vgFRZP6bD+hhoxxa2gNMgmARGh0zgQtJOzsTYb+0/JlaNcVFSzodecWUPzoHZUPYr0JCUm8IPwHbatARslgDPQOr1duhFsn0JbANPwEsQzwDzQEtSOweR7TNgY9h9Bz6AK0nNX2/BtituWllbkzTgWc/EApDK60rq2J4r3kD6PwaMAxNFG5WAL0elBLwB1rP9Zir4BJmUbAFx6PbeAAAAAElFTkSuQmCC) ![GitHub Workflow Status](https://img.shields.io/github/workflow/status/opq-osc/OPQBot/CodeQL?style=for-the-badge)
### 功能 😄
|功能|是否实现|
|-|-|
|群消息处理事件|是|
|好友消息处理事件|是|
|机器人事件处理|是|
|所有支持的消息发送|是|
|At|是|
|表情|是|
|撤回|是|
|下载文件|是|
|禁言|是|
|群公告|是|

### 安装 💡
` github.com/mcoo/OPQBot@latest`

请看 [wiki](https://mcenjoy.cn/opqbotwiki/)

以example文件为准 [example](https://github.com/opq-osc/OPQBot/blob/main/example/main.go)

### 使用本项目的一些开源项目 ❤️
|项目|介绍|
|-|-|
|[opq-osc/OPQBot-GroupManager](https://github.com/opq-osc/OPQBot-GroupManager)|群管理机器人示例 - 带web面板和模块化示例|
|[icepie/oh-my-lit](https://github.com/icepie/oh-my-lit)|Oh My LIT (吾之洛理)|
|[OPQBOT-jikipedia](https://github.com/opq-osc/OPQBOT-jikipedia)|查梗插件|
|[Liumik233/Aria2M-OPQBot](https://github.com/Liumik233/Aria2M-OPQBot)|本bot通过文字指令的形式将下载链接或bt种子下推送到Aria2服务器进行下载，以后可能会增加下载度盘链接功能|
|[Liumik233/OPQ-Plugins](https://github.com/Liumik233/OPQ-Plugins)|一些琐碎的OPQBOT插件|


以上排名不分先后，感谢！

### 没人看的更新历史 ✏️
```
20210318    简化发送代码
20210319    将宏移出BotManager,添加对发送队列每次发送时间的控制
20210322    添加发送函数的中间件
20210403    增加发送回调和优化中间件,基础功能完善
20210405    添加对撤回功能的完善和支持 注意看一下 example
20210406    戳一戳功能，example 即是文档
20210407    删除多余log，完善戳戳
20210420    添加Mp3转Silk功能和一些其他的功能
20210424    添加事件的中间件，向下兼容以前的代码，使用看example，完善silk功能
20210427    修复SocketIO数据畸形的问题，添加群上传功能
20210428    添加内置session 相关内容看Wiki
20210512    packet现在修改为传递指针，请注意
20210523    添加快捷发信息的函数
20210529    修复群文件上传BUG
20210531    修复Wait函数的BUG
```
