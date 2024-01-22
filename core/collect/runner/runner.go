package runner

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/FateBug403/FoFa/pkg/fofa"
	"github.com/FateBug403/OneForAll_go/pkg/oneforall"
	"github.com/FateBug403/fader/core/collect/config"
	"github.com/FateBug403/fader/global"
	httpx "github.com/FateBug403/httpx/runner"
	"github.com/FateBug403/util"
	"log"
	"net"
	"regexp"
	"strconv"
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

func (receiver *Runner) RunCollect() error {

	var hosts []string

	currentTime := time.Now()
	timestamp := currentTime.Format("20060102_150405")

	// 提取域名调用OneForAll进行查询
	var domains []string
	pattern := `^([a-zA-Z0-9.-]+\.[a-zA-Z]{2,})$`
	regex := regexp.MustCompile(pattern)
	for _,value:=range receiver.options.Targets{
		if regex.MatchString(value){
			domains = append(domains,value)
		}
	}
	ofa,err := oneforall.NewOneForAll(oneforall.Options{ExePath: receiver.options.OnforAllPath})
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("开始对"+strconv.Itoa(len(domains))+"个域名进行OneForAll域名收集")
	subdomains,err :=ofa.GetSubDomains(domains)
	if err != nil {
		log.Println(err)
		return err
	}
	hosts = append(hosts,subdomains...)
	log.Println("从OneForAll收集到"+strconv.Itoa(len(subdomains))+"条子域名")

	if len(subdomains)>0{
		err = util.WriteLineFile(fmt.Sprintf(receiver.options.OutputPath+"\\subdomains%s.txt",timestamp), subdomains)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	// 从Fofa采集数据并和子域名整理一起
	if global.CONFIG.FoFa.Mail!=""&&global.CONFIG.FoFa.Key!=""&&global.CONFIG.FoFa.Api!=""{
		fofaHost:=HostsProbing(receiver.options.Targets)
		log.Println("从Fofa收集到数据:"+ strconv.Itoa(len(fofaHost)))
		hosts = util.RemoveDuplicateElement(append(hosts,fofaHost...))
	}else {
		log.Println("没有配置Fofa邮箱、key或者Api，无法从Fofa获取数据")
	}

	// 保存收集到的所有host
	if len(hosts)>0{
		err = util.WriteLineFile(fmt.Sprintf(receiver.options.OutputPath+"\\all_host_%s.txt",timestamp), hosts)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	// 收集到的信息进行存活探测
	if receiver.options.AliveVerify{
		log.Println("开始对"+ strconv.Itoa(len(hosts)) +"条信息进行存活探测")
		aliveLinks,err := LinksAliveProbing(hosts)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("存活探测完毕，共计获取到"+ strconv.Itoa(len(aliveLinks)) +"条信息")

		if len(aliveLinks)>0{
			// 保存数据到表格中
			file := excelize.NewFile()
			sheetName := "Sheet1"
			//index := file.NewSheet(sheetName)
			// 设置表头
			headers := []string{"Url", "Host","Title", "StatusCode","Technologies","Words"}
			for colIndex, header := range headers {
				cell := excelize.ToAlphaString(colIndex+1) + "1"
				file.SetCellValue(sheetName, cell, header)
			}
			// 将数据写入表格
			for rowIndex, rowData := range aliveLinks {
				cell := excelize.ToAlphaString(1) + fmt.Sprintf("%d", rowIndex+2)
				file.SetCellValue(sheetName, cell, rowData.URL)
				cell = excelize.ToAlphaString(2) + fmt.Sprintf("%d", rowIndex+2)
				file.SetCellValue(sheetName, cell, rowData.Host)
				cell = excelize.ToAlphaString(3) + fmt.Sprintf("%d", rowIndex+2)
				file.SetCellValue(sheetName, cell, rowData.Title)
				cell = excelize.ToAlphaString(4) + fmt.Sprintf("%d", rowIndex+2)
				file.SetCellValue(sheetName, cell, rowData.StatusCode)
				cell = excelize.ToAlphaString(5) + fmt.Sprintf("%d", rowIndex+2)
				file.SetCellValue(sheetName, cell, rowData.Technologies)
				cell = excelize.ToAlphaString(6) + fmt.Sprintf("%d", rowIndex+2)
				file.SetCellValue(sheetName, cell, rowData.Words)
			}

			// 将结果保存到 Excel 文件
			err := file.SaveAs(fmt.Sprintf(receiver.options.OutputPath+"\\aliveLink_%s.xlsx",timestamp))
			if err != nil {
				fmt.Println("Error saving Excel file:", err)
				return err
			}
		}
	}

	return nil
}

// HostsProbing 通过域名或者IP从FOFA从收集HOST信息
func HostsProbing(rules []string) []string  {
	var hosts []string

	//提取规则里的域名还有IP地址
	var ips []string
	var domains []string
	// 定义一个正则表达式模式，用于匹配域名
	pattern := `^([a-zA-Z0-9.-]+\.[a-zA-Z]{2,})$`
	// 编译正则表达式
	regex := regexp.MustCompile(pattern)
	for _,value:=range rules{
		// 使用net.ParseIP()函数来尝试解析IP地址
		ip := net.ParseIP(value)
		if ip == nil {
			if regex.MatchString(value){
				domains = append(domains,value)
			}
			continue
		} else {
			ips = append(ips,value)
		}
	}

	// 1. 通过Fofa搜索用户输入的域名，获取链接
	log.Println("开始从Fofa进行链接收集")
	var domainRules []string
	FofaClient,err := fofa.NewFoFa(&fofa.Options{
		Baseurl: global.CONFIG.FoFa.Api,
		Email:   global.CONFIG.FoFa.Mail,
		Key:     global.CONFIG.FoFa.Key,
		Size:    10000,
	})
	if err != nil {
		log.Println(err)
		return hosts
	}

	for _,domainTmp:= range domains{
		domainRules = append(domainRules,"domain=\""+domainTmp+"\"")
	}

	FofaDomainResult:=FofaClient.SearchAllS(domainRules)
	hosts = append(hosts,FofaDomainResult.GetHosts()...)

	//2. 通过Fofa搜索用户输入的IP，获取链接
	var ipsRules []string
	for _,ipsTmp:= range ips{
		ipsRules = append(ipsRules,"ip=\""+ipsTmp+"\"")
	}
	FofaIpResult:=FofaClient.SearchAllS(ipsRules)
	hosts = append(hosts,FofaIpResult.GetHosts()...)

	return hosts
}

func LinksAliveProbing(targets []string) ([]httpx.Result,error) {
	option := httpx.DefaultOptions
	option.HTTPProxy=global.CONFIG.Proxy
	option.InputTargetHost = targets
	hp,err:=httpx.New(option)
	if err != nil {
		log.Println(err)
		return nil,err
	}
	var result []httpx.Result
	hp.RunEnumeration(func(r httpx.Result) {
		result = append(result, r)
	})
	return result,nil
}