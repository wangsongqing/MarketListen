package cmd

import (
	"MarketListen/app/controllers"
	"MarketListen/pkg/time"
	"github.com/spf13/cobra"
	"strconv"
)

// 可以根据参数名称--传参
func init() {
	rootCmd.AddCommand(marketCmd)
}

// 运行项目命令 go run main.go market
// 监听15分钟K线，如果有涨幅大于 Rise 的通知
var marketCmd = &cobra.Command{
	Use:     "market",
	Short:   "",
	Long:    ``,
	Example: "go run main.go market",
	Run: func(cmd *cobra.Command, args []string) {
		market := controllers.Market{}
		market.Rise = 0.8 // 涨跌幅度通知阀值
		market.BaseUrl = "https://www.okx.com/api/v5/market/index-candles"

		bar := "15m"
		instId := "LTC-USD"

		afterTime := (time.GetNowTime() - 1800) * 1000

		time := strconv.FormatInt(afterTime, 10)
		market.Market15m(instId, bar, time)
	},
}
