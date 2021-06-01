package config

import (
	"crawler/tencentKeTang/ffmpeg"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type HttpConfig struct {
	Cookie string  `yaml:"cookie"`
	BKN    int64   `yaml:"bkn"`
	T      float32 `yaml:"t"`
}

type AppConfig struct {
	Debug bool `yaml:"debug"`
}

type Config struct {
	Ffmpeg ffmpeg.Config `yaml:"ffmpeg"`
	Http   HttpConfig    `yaml:"http"`
	App    AppConfig     `yaml:"app"`
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
		c.Ffmpeg.SavePath = "./"
	}
	return c
}
