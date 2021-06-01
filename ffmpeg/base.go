package ffmpeg

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
)

type Config struct {
	Path     string `yaml:"path"`
	Params   string `yaml:"params"`
	SavePath string `yaml:"save_path"`
}

type Ffmpeg struct {
	c *Config
}

func New(c *Config) *Ffmpeg {
	return &Ffmpeg{c: c}
}

func (f *Ffmpeg) Do(vodUrl, name string) error {
	//检查文件是否存在
	savePath := f.c.SavePath + "/" + name
	err := os.MkdirAll(f.c.SavePath, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "os.MkdirAll path:%s", f.c.SavePath)
	}
	_, err = os.Stat(savePath)
	if err == nil {
		//文件存在，检查本地与目标视频时间差值是否小于1s
		//获取本地文件时常
		localDuration, err := exec.Command("sh",
			"-c",
			fmt.Sprintf(`ffmpeg -i "%s" 2>&1 | grep 'Duration' | cut -d ' ' -f 4 | sed s/,//`, savePath)).Output()
		if err != nil {
			return errors.Wrapf(err, "exec.Command err: %s, vodUrl:%s", err, vodUrl)
		}
		//没有查询到视频信息，说明不是完整视频，删除文件，开始下载
		if len(localDuration) == 0 {
			_ = os.Remove(savePath)
		}
	} else if !os.IsNotExist(err) {
		return errors.Wrap(err, "os.Stat")
	}
	cmd := exec.Command("sh", "-c", fmt.Sprintf(`ffmpeg -i "%s" -c:v h264_videotoolbox "%s"`, vodUrl, savePath))
	_, err = cmd.Output()
	if err != nil {
		return errors.Wrap(err, "exec.Command save")
	}
	return nil
}
