package ffmpeg

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
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

func (f *Ffmpeg) Do(vodUrl, name string, done int64) error {
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

	a, err := ffmpeg.Probe(vodUrl)
	if err != nil {
		panic(err)
	}
	totalDuration := gjson.Get(a, "format.duration").Float()

	err = ffmpeg.Input(vodUrl).
		Output(savePath, ffmpeg.KwArgs{"c:v": "h264_videotoolbox"}).
		GlobalArgs("-progress", "unix://"+f.TempSock(totalDuration, done)).
		OverWriteOutput().
		Run()
	if err != nil {
		panic(err)
	}
	return nil
}

func (f *Ffmpeg) TempSock(totalDuration float64, done int64) string {
	rand.Seed(time.Now().Unix())
	sockFileName := path.Join(os.TempDir(), fmt.Sprintf("%d_sock", rand.Int()))
	l, err := net.Listen("unix", sockFileName)
	if err != nil {
		panic(err)
	}

	go func() {
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
				cp = fmt.Sprintf("%.0f", float64(c)/totalDuration/1000000*100)
			}
			if strings.Contains(data, "progress=end") {
				cp = "done"
			}
			if cp != progress {
				progress = cp
				lastFps := fps[len(fps)-1]
				log.Printf("done:%d, progress: %s%%, fps: %s", done, progress, lastFps[len(lastFps)-1])
				fps = [][]string{}
			}
		}
	}()

	return sockFileName
}
