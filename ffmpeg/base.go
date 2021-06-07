package ffmpeg

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"log"
	"math"
	"net"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type Config struct {
	Path     string `yaml:"path"`
	Params   string `yaml:"params"`
	SavePath string `yaml:"save_path"`
}

type Ffmpeg struct {
	c              *Config
	ffmpegExec     string
	ffmpegParams   []string
	ffprobeExec    string
	address        string
	remoteDuration float64
}

func New(c *Config) *Ffmpeg {
	f := &Ffmpeg{
		c: c,
	}
	if c.Path != "" {
		f.ffmpegExec = c.Path + "/ffmpeg"
		f.ffprobeExec = c.Path + "/ffprobe"
	} else {
		f.ffmpegExec = "ffmpeg"
		f.ffprobeExec = "ffprobe"
	}
	if runtime.GOOS == "windows" {
		f.ffmpegExec += ".exe"
		f.ffprobeExec += ".exe"
	}
	if c.Params != "" {
		f.ffmpegParams = strings.Split(c.Params, " ")
	}
	return f
}

func (f *Ffmpeg) Do(vodUrl, name string) error {
	//获取目标视频帧数
	ret, err := f.probe(vodUrl)
	if err != nil {
		return errors.Wrap(err, "probe")
	}
	f.remoteDuration = gjson.Get(ret, "format.duration").Float()
	//检查文件是否存在
	savePath := f.c.SavePath + "/" + name + ".mp4"
	err = os.MkdirAll(f.c.SavePath, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "os.MkdirAll path:%s", f.c.SavePath)
	}
	_, err = os.Stat(savePath)
	if err == nil {
		//获取本地文件时常
		ret, err := f.probe(savePath)
		//没有查询到视频信息，说明不是完整视频，删除文件，开始下载
		if err != nil {
			_ = os.Remove(savePath)
		}
		localDuration := gjson.Get(ret, "format.duration").Float()
		if math.Abs(localDuration-f.remoteDuration) > 10 {
			//相差10帧以上则删除重新下载
			_ = os.Remove(savePath)
		} else {
			//小于10帧则成功处理
			return nil
		}
	} else if !os.IsNotExist(err) {
		return errors.Wrap(err, "os.Stat")
	}

	l, err := net.Listen("tcp", ":8829")
	if err != nil {
		panic(err)
	}
	go f.progress(l)
	f.address = "127.0.0.1:8829"

	err = f.mergeAndDownload(vodUrl, savePath, f.address)
	if err != nil {
		return errors.Wrap(err, "mergeAndDownload")
	}
	return nil
}

func (f *Ffmpeg) progress(l net.Listener) {
	re := regexp.MustCompile(`out_time_ms=(\d+)`)
	reFps := regexp.MustCompile(`fps=(\d+)`)
	fd, err := l.Accept()
	if err != nil {
		log.Fatal("accept error:", err)
	}
	buf := make([]byte, 16)
	data := ""
	progress := ""
	for {
		_, err := fd.Read(buf)
		if err != nil {
			return
		}
		data += string(buf)
		a := re.FindAllStringSubmatch(data, -1)
		cp := ""
		fps := reFps.FindAllStringSubmatch(data, -1)
		if len(a) > 0 && len(a[len(a)-1]) > 0 {
			c, _ := strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
			cp = fmt.Sprintf("%.0f", float64(c)/f.remoteDuration/1000000*100)
		}
		if strings.Contains(data, "progress=end") {
			cp = "done"
		}
		if cp != progress {
			progress = cp
			lastFps := fps[len(fps)-1]
			log.Printf("progress: %s%%, fps: %s", progress, lastFps[len(lastFps)-1])
			fps = [][]string{}
		}
	}
}
