package controllers

import (
	"MarketListen/pkg/config"
	"MarketListen/pkg/email"
	"MarketListen/pkg/helpers"
	"MarketListen/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"strconv"
	"strings"
	times "time"
)

type Market struct {
	BaseUrl string
	Rise    float64
}

type ResponseData struct {
	Code int        `json:"code"`
	Msg  string     `json:"msg"`
	Data [][]string `json:"data"`
}

func (m *Market) Market15m(instId string, bar string, time string) {

	request := gorequest.New()
	url := m.BaseUrl + "?instId=" + instId + "&bar=" + bar + "&before=" + time
	resp, body, errs := request.Get(url).
		EndBytes()

	if resp.StatusCode != 200 {
		if jsonErr, ok := json.Marshal(&errs); ok != nil {
			logger.Info("请求失败：" + string(jsonErr))
		}
		return
	}

	data := ResponseData{}
	json.Unmarshal(body, &data)

	if data.Code != 0 {
		logger.NewGormLogger().Info(context.TODO(), "数据异常："+helpers.FmtStrFromInterface(data.Msg)) // 写入日志
	}

	receive := helpers.FmtStrFromInterface(config.Env("EMAIL_QQ_RECEIVE", ""))
	splitReceive := strings.Split(receive, ",")
	for _, v := range data.Data {
		// fmt.Printf("ts:%v,o:%v,h:%v,l:%v,c:%v,confirm:%v \n", v[0], v[1], v[2], v[3], v[4], v[5])
		// (收盘价格 - 开盘价格) / 开盘价格  // 涨跌幅
		open, _ := strconv.ParseFloat(v[1], 0)
		ceil, _ := strconv.ParseFloat(v[4], 0)

		// fmt.Printf("open:%v,ceil:%v \n", v[1], v[4])

		rate := (ceil - open) / open
		resRate := helpers.Decimal(rate*100, "2")

		sendRate := fmt.Sprintf("%.2f", resRate)

		// 跌幅
		if resRate < 0 {
			if resRate*-1 > m.Rise {
				sendEmails(sendRate, instId, splitReceive)
				fmt.Printf("币种:%v,跌幅：%v,开盘价:%v,收盘价:%v \n", instId, resRate, open, ceil)
				break
			}
		} else {
			// 涨幅
			if resRate >= m.Rise {
				sendEmails(sendRate, instId, splitReceive)
				fmt.Printf("币种:%v,涨幅：%v,开盘价:%v,收盘价:%v \n", instId, resRate, open, ceil)
				break
			}
		}
	}

	fmt.Println("Done")
}

// 批量发送邮件，多个邮箱都可以
func sendEmails(sendRate string, instId string, splitReceive []string) {
	for _, sendEmail := range splitReceive {
		if len(sendEmail) == 0 {
			continue
		}
		email.SentEmail(sendEmail, sendRate, instId)
		times.Sleep(times.Second)
	}
}
