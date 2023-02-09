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
	marketCmd.Flags().StringP("instId", "i", "LTC-USD", "") // 币种
	marketCmd.Flags().StringP("bar", "b", "15m", "")        // K线时间段
	marketCmd.Flags().Float64P("wave", "w", 1, "")          // 涨跌幅通知阀值
}

// 运行项目命令 go run main.go market
// 监听15分钟K线，如果有涨幅大于 Rise 的通知
var marketCmd = &cobra.Command{
	Use:     "market",
	Short:   "",
	Long:    ``,
	Example: "go run main.go market -i OKB-USD -b 15m -w 2",
	Run: func(cmd *cobra.Command, args []string) {
		market := controllers.Market{}
		wave, _ := cmd.Flags().GetFloat64("wave")
		market.Rise = wave // 涨跌幅度通知阀值
		market.BaseUrl = "https://www.okx.com/api/v5/market/index-candles"

		bar, _ := cmd.Flags().GetString("bar")
		instId, _ := cmd.Flags().GetString("instId")

		afterTime := (time.GetNowTime() - 1800) * 1000

		time := strconv.FormatInt(afterTime, 10)
		market.Market15m(instId, bar, time)
	},
}
