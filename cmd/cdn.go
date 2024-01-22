/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/FateBug403/fader/core/cdn/config"
	"github.com/FateBug403/fader/global"
	"github.com/FateBug403/util"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/FateBug403/fader/core/cdn/runner"
)

// cdnCmd represents the cdn command
var cdnCmd = &cobra.Command{
	Use:   "cdn",
	Short: "A brief description of your command",
	Long: `用来探测域名是否采用了CDN技术，批量提取真实IP`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		_,err=os.Stat(cdnParam.InputFile)
		if os.IsNotExist(err){
			cobra.CheckErr(fmt.Errorf("输入文件未找到"))
		}

		_,err = os.Stat(global.CONFIG.CDN.DnsServer)
		if os.IsNotExist(err){
			cobra.CheckErr(fmt.Errorf("没有找到配置的DNS服务器地址的文件,请在config.yaml文件中指定或者重新配置"))
		}

		if cdnParam.OutPut==""{
			cdnParam.OutPut = "CdnOut"
		}
		err = util.CreatePath(cdnParam.OutPut)
		if err != nil {
			cobra.CheckErr(err)
		}

		domains := util.ReadFile(cdnParam.InputFile)
		if len(domains)<1{
			cobra.CheckErr(fmt.Errorf("文件为空，未找到目标"))
		}

		options := &config.Options{
			Domains:    domains,
			OutputPath: cdnParam.OutPut,
		}
		r,err := runner.NewRunner(options)
		err = r.RunCDNChecks()
		if err != nil {
			log.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(cdnCmd)

	cdnCmd.Flags().StringVarP(&cdnParam.InputFile, "inputFile", "f","", "-f 或者 --inputFile 来传递输入文件,可以包含域名、IP和Host")
	cdnCmd.Flags().StringVarP(&cdnParam.OutPut, "output", "o","", "-o 或者 --output 来传递输出的ip地址的文件")
	err := cdnCmd.MarkFlagRequired("inputFile")
	if err != nil {
		return
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cdnCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cdnCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
