package controllers

import (
	"MarketListen/pkg/logger"
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"strings"
	"time"
)

type Markets struct {
	BaseUrl    string
	Rise       float64
	TickersUrl string
}

type ResponseDatas struct {
	Code string `json:"code"`
	Data []struct {
		AskPx     string `json:"askPx"`
		AskSz     string `json:"askSz"`
		BidPx     string `json:"bidPx"`
		BidSz     string `json:"bidSz"`
		High24h   string `json:"high24h"`
		InstID    string `json:"instId"`
		InstType  string `json:"instType"`
		Last      string `json:"last"`
		LastSz    string `json:"lastSz"`
		Low24h    string `json:"low24h"`
		Open24h   string `json:"open24h"`
		SodUtc0   string `json:"sodUtc0"`
		SodUtc8   string `json:"sodUtc8"`
		Ts        string `json:"ts"`
		Vol24h    string `json:"vol24h"`
		VolCcy24h string `json:"volCcy24h"`
	} `json:"data"`
	Msg string `json:"msg"`
}

func (m *Markets) SendPush(bar string, times string) {
	market := Market{}
	market.Rise = m.Rise
	market.BaseUrl = m.BaseUrl

	request := gorequest.New()
	resp, body, errs := request.Get(m.TickersUrl).EndBytes()

	if resp.StatusCode != 200 {
		if jsonErr, ok := json.Marshal(&errs); ok != nil {
			logger.Info("请求失败：" + string(jsonErr))
		}
		return
	}

	data := ResponseDatas{}
	json.Unmarshal(body, &data)

	for _, v := range data.Data {
		if ok := strings.Contains(v.InstID, "USDT"); !ok {
			continue
		}
		//fmt.Printf("k:%v,v:%v \n", k, v.InstID)
		market.Market15m(v.InstID, bar, times)

		time.Sleep(time.Second)
	}

}
