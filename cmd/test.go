/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/FateBug403/FoFa/pkg/fofa"
	"github.com/FateBug403/fader/global"
	httpx "github.com/FateBug403/httpx/runner"
	"github.com/spf13/cobra"
	"log"
	"os"
	"regexp"
)

// TestCmd represents the Test command
var TestCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `用于对一些配置做测试`,
	Run: func(cmd *cobra.Command, args []string) {
		if testParam.Proxy { // 测试配置文件代理配置
			regInfo := regexp.MustCompile("<pre>[\\s\\S]*</pre>")

			option := httpx.DefaultOptions
			option.HTTPProxy=global.CONFIG.Proxy
			hp,err:=httpx.New(option)
			if err != nil {
				log.Println(err)
				return
			}
			hp.RunAlone("cip.cc", func(r httpx.Result) {
				if r.StatusCode == 200{
					fmt.Println(regInfo.FindString(r.ResponseDateStr))
					os.Exit(0)
				}else {
					log.Println("获取出错，请检查网络配置")
				}
			})
			os.Exit(0)
		}

		if testParam.FoFa{
			_,err := fofa.NewFoFa(&fofa.Options{
				Baseurl: global.CONFIG.FoFa.Api,
				Email:   global.CONFIG.FoFa.Mail,
				Key:     global.CONFIG.FoFa.Key,
				Size:    10000,
			})
			if err != nil {
				log.Println(err)
			}
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(TestCmd)
	TestCmd.Flags().BoolVarP(&testParam.Proxy, "proxy", "p",false, "-p 或者 --proxy 测试配置文件中的代理配置是否生效")
	TestCmd.Flags().BoolVarP(&testParam.FoFa, "fofa", "f",false, "-f 或者 --fofa 测试配置文件中的FoFa账号信息是否正确")
}
