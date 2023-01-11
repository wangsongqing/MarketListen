package cmd

import (
	"MarketListen/app/controllers"
	"github.com/spf13/cobra"
)

// 可以根据参数名称--传参
func init() {
	rootCmd.AddCommand(indexCmd)
}

// 运行项目命令 go run main.go index

var indexCmd = &cobra.Command{
	Use:     "index",
	Short:   "",
	Long:    ``,
	Example: "go run main.go index",
	Run: func(cmd *cobra.Command, args []string) {

		twitter := controllers.Twitter{}
		twitter.Url = "https://twitter.com/elonmusk"
		twitter.ListenTwitter()
		//twitter.SentEmail()
	},
}
