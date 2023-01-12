package email

import (
	"MarketListen/pkg/config"
	"MarketListen/pkg/helpers"
	"fmt"
	"gopkg.in/gomail.v2"
	"strconv"
)

// SentEmail wave 涨跌幅
// toMan 接收人
// coin 币中
func SentEmail(receive string, wave string, coin string) {
	m := gomail.NewMessage()

	sendEmail := helpers.FmtStrFromInterface(config.Env("EMAIL_QQ_ACCOUNT", ""))
	password := helpers.FmtStrFromInterface(config.Env("EMAIL_QQ_PASSWORD", ""))
	port, _ := strconv.Atoi(config.Env("EMAIL_QQ_PORT", "").(string))

	body := "<h4>币种:" + coin + "</h4>"
	body += "<h4>十五分钟线波动:" + wave + "%</h4>"

	//发送人
	m.SetHeader("From", sendEmail)
	//接收人
	m.SetHeader("To", receive)
	//抄送人
	//m.SetAddressHeader("Cc", "xxx@qq.com", "xiaozhujiao")
	//主题
	m.SetHeader("Subject", "K线波动剧烈")
	//内容
	m.SetBody("text/html", body)
	//附件
	//m.Attach("./myIpPic.png")

	//拿到token，并进行连接,第4个参数是填授权码
	d := gomail.NewDialer("smtp.qq.com", port, sendEmail, password)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("DialAndSend err %v:", err)
		panic(err)
	}
	fmt.Printf("send mail success\n")
}
