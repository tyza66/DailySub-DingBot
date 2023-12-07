package main

import (
	"github.com/blinkbean/dingtalk"
	"github.com/jasonlvhit/gocron"
	"net/http"
)

func main() {
	bot := dingtalk.InitDingTalkWithSecret("ea098c63f230ea3ad7fa7d0dbc09f262777ae9325332da668f89d95f1bf30da6", "SECb2847ebb5940e68a2dabc9e25cf4a280d6cbe08e32a38a11c5059b8817c33123")
	bot.SendTextMessage("订阅服务已经启动")
	gocron.Every(1).Day().At("21:30").Do(tellTheFlow)
	<-gocron.Start()
}

func tellTheFlow() {
	bot := dingtalk.InitDingTalkWithSecret("ea098c63f230ea3ad7fa7d0dbc09f262777ae9325332da668f89d95f1bf30da6", "SECb2847ebb5940e68a2dabc9e25cf4a280d6cbe08e32a38a11c5059b8817c33123")
	resp, err := http.Get("http://mxkj089.cn/app/simCard/getCardFlow?card=03224303")
	if err != nil {
		bot.SendTextMessage("获取CPE当前流量失败")
		return
	}
	defer resp.Body.Close()

	bot.SendTextMessage("订阅服务已经启动")

}
