package main

import (
	"crawler/tencentKeTang/config"
	"crawler/tencentKeTang/project"
	"flag"
	"os"
)

func main() {
	//加载入参
	taskUrl := ""
	flag.StringVar(&taskUrl, "u", "", "初始任务链接")
	flag.Parse()
	if taskUrl == "" {
		panic("url is empty")
	}
	//加载配置文件
	configPath := ""
	if os.Getenv("tencentKeTang") == "dev" {
		configPath = "./config_dev.yaml"
	} else {
		configPath = "./config.yaml"
	}

	c := config.Load(configPath)
	if c.Http.Cookie == "" {
		panic("cookie is empty")
	}
	//执行任务
	if err := project.New(c).Do(taskUrl); err != nil {
		panic(err)
	}
	return
}
