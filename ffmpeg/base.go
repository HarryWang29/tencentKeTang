package ffmpeg

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
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
	c            *Config
	ffmpegExec   string
	ffmpegParams []string
	ffprobeExec  string
}

func New(c *Config) *Ffmpeg {
	f := &Ffmpeg{
		c: c,
	}
	if f.ffmpegExec != "" {
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
	remoteDuration := gjson.Get(ret, "format.duration").Float()
	//检查文件是否存在
	savePath := f.c.SavePath + "/" + name + ".mp4"
	err = os.MkdirAll(f.c.SavePath, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "os.MkdirAll path:%s", f.c.SavePath)
	}
	_, err = os.Stat(savePath)
	if err == nil {
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

	err = f.mergeAndDownload(vodUrl, savePath, f.TempSock(remoteDuration))
	if err != nil {
		return errors.Wrap(err, "mergeAndDownload")
	}
	return nil
}

func (f *Ffmpeg) TempSock(totalDuration float64) string {
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
				log.Printf("progress: %s%%, fps: %s", progress, lastFps[len(lastFps)-1])
				fps = [][]string{}
			}
		}
	}()

	return sockFileName
}
