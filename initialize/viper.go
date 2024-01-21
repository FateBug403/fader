package initialize

import (
	"fmt"
	"github.com/FateBug403/fader/global"
	"github.com/FateBug403/util"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)

// InitViper 初始化配置信息
func InitViper() error {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		// 不存在配置文件则重新生成,重新读取
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			cs := util.StructToMap(global.CONFIG)
			for key,value  := range cs {
				v.Set(key, value)
			}
			if err := v.WriteConfigAs("config.yaml"); err != nil {
				log.Println(err)
				return err
			}
			err = v.ReadInConfig()
			if err != nil {
				log.Println(err)
				return err
			}
			log.Println("没有找到配置文件，已成功生成，请重新运行程序")
			os.Exit(0)
		} else {
			fmt.Println("Error reading config file:", err)
			return err
		}
	}

	// 监视配置文件的变化并自动重新加载配置
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&global.CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&global.CONFIG); err != nil {
		fmt.Println(err)
		return err
	}

	// 初始化到全局变量
	global.VP =v

	return err
}