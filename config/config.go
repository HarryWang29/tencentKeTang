package config

import (
	"crawler/tencentKeTang/ffmpeg"
	"crawler/tencentKeTang/keTang"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type AppConfig struct {
	Debug bool `yaml:"debug"`
}

type Config struct {
	Ffmpeg ffmpeg.Config `yaml:"ffmpeg"`
	Http   keTang.Config `yaml:"http"`
	App    AppConfig     `yaml:"app"`

	KeTang keTang.Api
}

func Load(path string) *Config {
	c := &Config{}
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(bs, c)
	if err != nil {
		panic(err)
	}
	if c.Ffmpeg.SavePath == "" {
		c.Ffmpeg.SavePath = "./download"
	}
	return c
}
