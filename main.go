package main

import (
	"encoding/json"
	"fmt"
	"github.com/blinkbean/dingtalk"
	"github.com/jasonlvhit/gocron"
	"io/ioutil"
	"net/http"
)

// 用于存储昨天的流量情况
var lastDayFlow float64 = 0.0

func main() {
	bot := dingtalk.InitDingTalkWithSecret("ea098c62f330ea3ad7fa7d0dbc09f262777ae9325332da668f89d95f1bf30da6", "SECb2847ebb5940e68a2dabc9e25cf4a280d6cbe08e32a38a11c5059b8817c33123")
	bot.SendTextMessage("订阅服务已经启动")
	//每天晚上九点半汇报当时的流量情况
	tellTheFlow()
	gocron.Every(1).Day().At("21:30").Do(tellTheFlow)
	<-gocron.Start()
}

func tellTheFlow() {
	bot := dingtalk.InitDingTalkWithSecret("ea098c62f330ea3ad7fa7d0dbc09f262777ae9325332da668f89d95f1bf30da6", "SECb2847ebb5940e68a2dabc9e25cf4a280d6cbe08e32a38a11c5059b8817c33123")
	resp, err := http.Get("http://mxkj089.cn/app/simCard/phoneSimCard?card=03224303")
	if err != nil {
		bot.SendTextMessage("获取CPE当前流量失败")
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var jsonMap map[string]any
	json.Unmarshal(body, &jsonMap)
	if jsonMap["code"].(float64) == 0 {
		var msg = "获取当前CPE流量" + jsonMap["msg"].(string)
		msg += "，流量总量" + fmt.Sprintf("%.2f", jsonMap["data"].(map[string]any)["sumFlow"].(float64)) + "MB"
		msg += "，已用流量" + fmt.Sprintf("%.2f", jsonMap["data"].(map[string]any)["consumeFlow"].(float64)) + "MB"
		msg += "，剩余流量" + fmt.Sprintf("%.2f", jsonMap["data"].(map[string]any)["surplusFlow"].(float64)) + "MB"
		msg += "，流量过期时间为" + jsonMap["data"].(map[string]any)["mealEndTime"].(string)
		msg += "，校比昨日使用" + fmt.Sprintf("%.2f", jsonMap["data"].(map[string]any)["consumeFlow"].(float64)-lastDayFlow) + "MB"
		msg += "，卡上余额" + fmt.Sprintf("%.2f", jsonMap["data"].(map[string]any)["balance"].(float64)) + "元"
		msg += "，卡号为" + jsonMap["data"].(map[string]any)["meCardNum"].(string)
		msg += "，ICCID为" + jsonMap["data"].(map[string]any)["iccId"].(string)
		lastDayFlow = jsonMap["data"].(map[string]any)["consumeFlow"].(float64)
		bot.SendTextMessage(msg)
		return
	} else {
		get, err := http.Get("http://mxkj089.cn/app/simCard/getCardFlow?card=03224303")
		if err != nil {
			bot.SendTextMessage("备线获取CPE当前流量失败")
			return
		}
		defer get.Body.Close()
		body, err := ioutil.ReadAll(get.Body)
		var jsonMap map[string]any
		json.Unmarshal(body, &jsonMap)
		if jsonMap["code"].(float64) == 0 {
			var msg = "获取当前CPE流量" + jsonMap["msg"].(string)
			msg += "，流量总量" + fmt.Sprintf("%.2f", jsonMap["data"].(map[string]any)["sumFlow"].(float64)) + "MB"
			msg += "，已用流量" + fmt.Sprintf("%.2f", jsonMap["data"].(map[string]any)["consumeFlow"].(float64)) + "MB"
			msg += "，剩余流量" + fmt.Sprintf("%.2f", jsonMap["data"].(map[string]any)["surplusFlow"].(float64)) + "MB"
			msg += "，流量过期时间为" + jsonMap["data"].(map[string]any)["mealEndTime"].(string)
			msg += "，校比昨日使用" + fmt.Sprintf("%.2f", jsonMap["data"].(map[string]any)["consumeFlow"].(float64)-lastDayFlow) + "MB"
			lastDayFlow = jsonMap["data"].(map[string]any)["consumeFlow"].(float64)
			bot.SendTextMessage(msg)
			return
		} else {
			bot.SendTextMessage("所有获取CPE当前流量接口失效")
			return
		}
	}
}
