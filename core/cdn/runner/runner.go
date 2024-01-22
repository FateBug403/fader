package runner

import (
	"bufio"
	"fmt"
	"github.com/FateBug403/cdn"
	"github.com/FateBug403/fader/core/cdn/config"
	"github.com/FateBug403/fader/global"
	"log"
	"os"
	"time"
)

type Runner struct {
	options *config.Options
}

func NewRunner(Options *config.Options)(*Runner,error){
	var err error
	runner := &Runner{options: Options}
	return runner,err
}

func (receiver *Runner)RunCDNChecks() error {
	var err error

	currentTime := time.Now()
	timestamp := currentTime.Format("20060102_150405")
	// 打开文件，如果文件不存在则创建
	file, err := os.Create(fmt.Sprintf(receiver.options.OutputPath+"\\ips%s.txt",timestamp))
	if err != nil {
		return err
	}
	defer file.Close()
	// 创建一个带缓冲的写入器
	writer := bufio.NewWriter(file)
	ipsMap:=make(map[string]bool)
	options := &cdn.Options{
		DnsOerverFile: global.CONFIG.CDN.DnsServer,
		OnResult: func(s string) {
			if _,ok:=ipsMap[s];!ok{
				// 写入内容到缓冲区
				_, err := writer.WriteString(s+"\n")
				if err != nil {
					return
				}
				// 将缓冲区的内容刷入文件
				err = writer.Flush()
				if err != nil {
					return
				}
				ipsMap[s]=true
				log.Println(s)
			}
		},
	}
	CDNClient := cdn.NewCDNClient(options)
	_,err=CDNClient.CDNChecks(receiver.options.Domains)
	if err != nil {
		return err
	}
	return nil
}
