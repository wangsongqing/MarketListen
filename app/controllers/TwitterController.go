package controllers

import (
	"MarketListen/pkg/config"
	"MarketListen/pkg/helpers"
	"fmt"
	"github.com/gocolly/colly"
	"gopkg.in/gomail.v2"
	"log"
	"strconv"
	"time"
)

type Twitter struct {
	Url string
}

func (t *Twitter) SentEmail() {
	m := gomail.NewMessage()

	sendEmail := helpers.FmtStrFromInterface(config.Env("EMAIL_QQ_ACCOUNT", ""))
	password := helpers.FmtStrFromInterface(config.Env("EMAIL_QQ_PASSWORD", ""))
	port, _ := strconv.Atoi(config.Env("EMAIL_QQ_PORT", "").(string))

	//发送人
	m.SetHeader("From", sendEmail)
	//接收人
	m.SetHeader("To", "18201197923@163.com")
	//抄送人
	//m.SetAddressHeader("Cc", "xxx@qq.com", "xiaozhujiao")
	//主题
	m.SetHeader("Subject", "马斯克")
	//内容
	m.SetBody("text/html", "<h1>Twitt Success</h1>")
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

func (t *Twitter) ListenTwitter() {
	var Count int
	c := colly.NewCollector(
		//colly.Async(true),并发
		colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"),
	)
	c.SetRequestTimeout(time.Duration(35) * time.Second)
	c.Limit(&colly.LimitRule{DomainGlob: t.Url, Parallelism: 1}) //Parallelism代表最大并发数
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})
	//获取不同地区的总页数
	c.OnHTML("article", func(e *colly.HTMLElement) {
		list := e.ChildText("time")
		fmt.Println(list)
		Count += len(list)
	})
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
		c.Visit(t.Url)
	})
	c.Visit(t.Url)
	c.Wait()

	fmt.Printf("总数为%v", Count)
}
