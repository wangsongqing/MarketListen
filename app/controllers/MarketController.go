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

	for _, v := range data.Data {
		//fmt.Printf("ts:%v,o:%v,h:%v,l:%v,c:%v,confirm:%v \n", v[0], v[1], v[2], v[3], v[4], v[5])
		// (开盘价格 - 收盘价格) / 开盘价格
		open, _ := strconv.ParseFloat(v[1], 0)
		ceil, _ := strconv.ParseFloat(v[4], 0)

		rate := (open - ceil) / open
		resRate := helpers.Decimal(rate*100, "2")

		sendRate := fmt.Sprintf("%.2f", resRate)

		// 跌幅
		if resRate < 0 {
			if resRate*-1 > m.Rise {
				email.SentEmail(receive, sendRate, instId)
				fmt.Printf("币种:%v,跌幅：%v \n", instId, resRate)
				break
			}
		} else {
			// 涨幅
			if resRate >= m.Rise {
				email.SentEmail(receive, sendRate, instId)
				fmt.Printf("币种:%v,涨幅：%v \n", instId, resRate)
				break
			}
		}
	}

	fmt.Println("Done")
}
