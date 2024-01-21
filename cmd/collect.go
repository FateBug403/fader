/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/FateBug403/fader/core/collect/config"
	"github.com/FateBug403/fader/core/collect/runner"
	"github.com/FateBug403/fader/global"
	"github.com/FateBug403/util"
	"github.com/spf13/cobra"
	"os"
)

// collectCmd represents the collect command
var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "对目标的信息进行收集",
	Long: `对目标的信息进行收集`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		_,err=os.Stat(cliParam.InputFile)
		if os.IsNotExist(err){
			cobra.CheckErr(fmt.Errorf("输入文件未找到"))
		}

		_,err = os.Stat(global.CONFIG.OneForAll)
		if os.IsNotExist(err){
			cobra.CheckErr(fmt.Errorf("没有找到系统指定的OneForAll文件夹,请在config.yaml文件中指定或者重新配置"))
		}

		if cliParam.CollectOutput==""{
			cliParam.CollectOutput = "collectOut"
		}
		err = util.CreatePath(cliParam.CollectOutput)
		if err != nil {
			cobra.CheckErr(err)
		}

		// 开始运行Collect
		options := &config.Options{
			AliveVerify:  cliParam.AliveVerify,
			OutputPath:   cliParam.CollectOutput,
			Targets:      util.ReadFile(cliParam.InputFile),
			OnforAllPath: global.CONFIG.OneForAll,
		}
		r,err := runner.NewRunner(options)
		if err != nil {
			cobra.CheckErr(err)
		}
		err=r.RunCollect()
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(collectCmd)
	collectCmd.Flags().StringVarP(&cliParam.InputFile, "inputFile", "f","", "-f 或者 --inputFile 来传递输入文件")
	collectCmd.Flags().StringVarP(&cliParam.CollectOutput, "output", "o","", "-o 或者 --output 来传递输出文件")
	collectCmd.Flags().BoolVarP(&cliParam.AliveVerify,"aliveVerify","a",true,"-a 或者 --aliveVerify 来对收集的网站进行存活探测")
	err := collectCmd.MarkFlagRequired("inputFile")
	if err != nil {
		return
	}
}
