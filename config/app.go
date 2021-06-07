package config

import (
	"crawler/tencentKeTang/ffmpeg"
	"crawler/tencentKeTang/keTang"
	"crawler/tencentKeTang/project"
	"os"
)

type App struct {
	Config  *Config
	KeTang  keTang.Api
	Project project.Api
	FFmpeg  *ffmpeg.Ffmpeg
}

func NewApp() *App {
	//加载配置文件
	configPath := ""
	if os.Getenv("tencentKeTang") == "dev" {
		configPath = "./config_dev.yaml"
	} else {
		configPath = "./config.yaml"
	}

	app := &App{
		Config: Load(configPath),
	}
	var err error
	app.KeTang = keTang.New(&app.Config.Http)
	app.FFmpeg, err = ffmpeg.New(&app.Config.Ffmpeg)
	if err != nil {
		panic(err)
	}
	app.Project = project.New(app.KeTang, app.FFmpeg, app.Config.Http.Cookie)
	return app
}
