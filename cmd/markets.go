package cmd

import (
	"MarketListen/app/controllers"
	"MarketListen/pkg/time"
	"github.com/spf13/cobra"
	"strconv"
)

// 可以根据参数名称--传参
func init() {
	rootCmd.AddCommand(marketsCmd)
	marketsCmd.Flags().StringP("bar", "b", "15m", "") // K线时间段
	marketsCmd.Flags().Float64P("wave", "w", 2, "")   // 涨跌幅通知阀值
}

// 运行项目命令 go run main.go markets -b 15m -w 2
// 监听所有的币种，如果15分钟K线涨跌幅大于2%就发送通知
var marketsCmd = &cobra.Command{
	Use:     "markets",
	Short:   "",
	Long:    ``,
	Example: "go run main.go markets -b 15m -w 2",
	Run: func(cmd *cobra.Command, args []string) {
		markets := controllers.Markets{}
		wave, _ := cmd.Flags().GetFloat64("wave")
		markets.Rise = wave // 涨跌幅度通知阀值
		markets.BaseUrl = "https://www.okx.com/api/v5/market/index-candles"
		markets.TickersUrl = "https://www.okx.com/api/v5/market/tickers?instType=SPOT"
		bar, _ := cmd.Flags().GetString("bar")
		afterTime := (time.GetNowTime() - 1800) * 1000

		time := strconv.FormatInt(afterTime, 10)
		markets.SendPush(bar, time)
	},
}
