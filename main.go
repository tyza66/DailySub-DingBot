package main

import (
	"encoding/json"
	"fmt"
	"github.com/blinkbean/dingtalk"
	"github.com/jasonlvhit/gocron"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// 用于存储昨天的流量情况
var lastDayFlow float64 = 0.0

func main() {
	bot := dingtalk.InitDingTalkWithSecret("ea098c62f330ea3ad7fa7d0dbc09f262777ae9325332da668f89d95f1bf30da6", "SECb2847ebb5940e68a2dabc9e25cf4a280d6cbe08e32a38a11c5059b8817c33123")
	bot.SendTextMessage("订阅服务已经启动")
	//每天晚上九点半汇报当时的流量情况
	tellTheFlow()
	testSpeed()
	gocron.Every(1).Day().At("21:30").Do(tellTheFlow)
	//每天早上七点半、中午十一点半、晚上八点半测速
	gocron.Every(2).Day().At("07:30").Do(testSpeed)
	gocron.Every(2).Day().At("11:30").Do(testSpeed)
	gocron.Every(2).Day().At("20:30").Do(testSpeed)
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
		var msg = "### 流量使用情况" //"获取当前CPE流量" + jsonMap["msg"].(string)
		msg += "  \n- 流量总量：" + fmt.Sprintf("%.2f", jsonMap["data"].(map[string]any)["sumFlow"].(float64)) + "MB"
		msg += "  \n- 已用流量：" + fmt.Sprintf("%.2f", jsonMap["data"].(map[string]any)["consumeFlow"].(float64)) + "MB"
		msg += "  \n- 剩余流量：" + fmt.Sprintf("%.2f", jsonMap["data"].(map[string]any)["surplusFlow"].(float64)) + "MB"
		msg += "  \n- 过期时间：" + jsonMap["data"].(map[string]any)["mealEndTime"].(string)
		msg += "  \n- 较比昨日：" + fmt.Sprintf("%.2f", jsonMap["data"].(map[string]any)["consumeFlow"].(float64)-lastDayFlow) + "MB"
		msg += "  \n- 卡上余额：" + fmt.Sprintf("%.2f", jsonMap["data"].(map[string]any)["balance"].(float64)) + "元"
		msg += "  \n- 使用卡号：" + jsonMap["data"].(map[string]any)["meCardNum"].(string)
		msg += "  \n- ICCID：" + jsonMap["data"].(map[string]any)["iccId"].(string)
		lastDayFlow = jsonMap["data"].(map[string]any)["consumeFlow"].(float64)
		bot.SendMarkDownMessage("流量使用情况", msg)
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
			var msg = "### 流量使用情况" //"获取当前CPE流量" + jsonMap["msg"].(string)
			msg += "  \n- 流量总量：" + fmt.Sprintf("%.2f", jsonMap["data"].(map[string]any)["sumFlow"].(float64)) + "MB"
			msg += "  \n- 已用流量：" + fmt.Sprintf("%.2f", jsonMap["data"].(map[string]any)["consumeFlow"].(float64)) + "MB"
			msg += "  \n- 剩余流量：" + fmt.Sprintf("%.2f", jsonMap["data"].(map[string]any)["surplusFlow"].(float64)) + "MB"
			msg += "  \n- 过期时间：" + jsonMap["data"].(map[string]any)["mealEndTime"].(string)
			msg += "  \n- 较比昨日：" + fmt.Sprintf("%.2f", jsonMap["data"].(map[string]any)["consumeFlow"].(float64)-lastDayFlow) + "MB"
			lastDayFlow = jsonMap["data"].(map[string]any)["consumeFlow"].(float64)
			bot.SendMarkDownMessage("流量使用情况", msg)
			return
		} else {
			bot.SendTextMessage("所有获取CPE当前流量接口失效")
			return
		}
	}
}

func testSpeed() {
	bot := dingtalk.InitDingTalkWithSecret("ea098c62f230ea3ad7fa7d0dbc09f262777ae9325332da668f89d95f1bf30da6", "SECb2847ebb5940e68a2dabc9e25cf4a280d6cbe08e32a38a11c5059b8817c33123")
	msg := "### 测速汇总"
	msg += "  \n- 百度服务器：" + getSpeed("www.baidu.com:80")
	msg += "  \n- 腾讯服务器：" + getSpeed("www.qq.com:80")
	msg += "  \n- 阿里服务器：" + getSpeed("www.aliyun.com:80")
	msg += "  \n- 微信服务器：" + getSpeed("www.weixin.com:80")
	msg += "  \n- 网易服务器：" + getSpeed("www.163.com:80")
	msg += "  \n- 新浪服务器：" + getSpeed("www.sina.com:80")
	msg += "  \n- 搜狐服务器：" + getSpeed("www.sohu.com:80")
	msg += "  \n- 谷歌服务器：" + getSpeed("www.google.com:80")
	bot.SendMarkDownMessage("测速汇总", msg)
}

func getSpeed(url string) string {
	timeout := time.Duration(5 * time.Second)
	start := time.Now()
	_, err := net.DialTimeout("tcp", url, timeout)
	if err != nil {
		fmt.Println("Error:", err)
		return "超时"
	}
	elapsed := time.Since(start)
	return elapsed.String()
}
