package ffmpeg

import (
	"fmt"
	"github.com/HarryWang29/progressbar/v4"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type Config struct {
	Path     string `yaml:"path"`
	Params   string `yaml:"params"`
	SavePath string `yaml:"save_path"`
	Worker   int    `yaml:"worker"`
	Port     int    `yaml:"port"`
}

type Ffmpeg struct {
	c              *Config
	ffmpegExec     string
	ffmpegParams   []string
	ffprobeExec    string
	address        string
	remoteDuration float64
	listener       net.Listener
	workerChannel  chan *task
	finishChannel  chan *task
	taskMap        sync.Map
	finishMap      sync.Map
	httpclient     *http.Client
	bars           *progressbar.MultiProgress
	mergeFileName  string
}

func New(c *Config) (*Ffmpeg, error) {
	f := &Ffmpeg{
		c:    c,
		bars: progressbar.NewMultiProgress(),
	}
	f.bars.Start()
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
	//检查ffmpeg
	err := f.checkFfmpeg()
	if err != nil {
		return nil, errors.Wrap(err, "调用ffmpeg出错，请检查地址")
	}
	//检查ffprobe
	err = f.checkProbe()
	if err != nil {
		return nil, errors.Wrap(err, "调用ffprobe出错，请检查地址")
	}

	f.workerChannel = make(chan *task, f.c.Worker)
	for i := 0; i < f.c.Worker; i++ {
		go f.taskProcess()
	}

	//启动一个协程进行合并视频
	f.finishChannel = make(chan *task, 1)
	go f.finish()

	//打开ffmpeg监听端口
	f.address = "localhost:" + strconv.Itoa(f.c.Port)
	f.listener, err = net.Listen("tcp", f.address)
	if err != nil {
		return nil, errors.Wrap(err, "net.Listen")
	}
	go f.progress()

	f.httpclient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 100,
		},
	}
	return f, nil
}

func (f *Ffmpeg) Do(vodUrl, dk string, bitrate int, path []string) error {
	//获取目标视频帧数
	ret, err := f.probe(vodUrl)
	if err != nil {
		return errors.Wrap(err, "probe")
	}
	f.remoteDuration = gjson.Get(ret, "format.duration").Float()
	//检查文件是否存在
	path = append([]string{f.c.SavePath}, path...)
	//savePath := util.PathJoin(path...) + ".mp4"
	savePath := filepath.Join(path...) + ".mp4"
	saveDir, _ := filepath.Split(savePath)
	err = os.MkdirAll(saveDir, os.ModePerm)
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

	err = f.downloadTs(vodUrl, dk, bitrate, savePath)
	if err != nil {
		return errors.Wrap(err, "mergeAndDownload")
	}
	return nil
}

func (f *Ffmpeg) addMergeBar(name string, max int) error {
	return f.addBar("merge|"+name, fmt.Sprintf("[fps:0000] merging %s ...", name), max)
}

func (f *Ffmpeg) addDownloadBar(name string, max int) error {
	return f.addBar("down|"+name, fmt.Sprintf("downloading %s ...", name), max)
}

func (f *Ffmpeg) addBar(name, desc string, max int) error {
	return f.bars.AddBar(name, progressbar.NewOptions(max,
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription(desc),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		})))
}

func (f *Ffmpeg) progress() {
	for {
		fd, err := f.listener.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}
		f.progressRead(fd)
	}
}

func (f *Ffmpeg) progressRead(fd net.Conn) {
	re := regexp.MustCompile(`out_time_ms=(\d+)`)
	reFps := regexp.MustCompile(`fps=(\d+)`)
	buf := make([]byte, 16)
	data := ""
	fpsShow := ""
	for {
		_, err := fd.Read(buf)
		if err != nil {
			return
		}
		data += string(buf)
		datas := strings.Split(data, "\n")
		for i := 0; i < len(datas)-1; i++ {
			fps := reFps.FindStringSubmatch(datas[i])
			if len(fps) > 0 {
				fpsShow = fps[len(fps)-1]
			}
			a := re.FindStringSubmatch(datas[i])
			if len(a) > 0 {
				c, err := strconv.Atoi(a[len(a)-1])
				if err != nil {
					continue
				}
				f.bars.BarSet("merge|"+f.mergeFileName, c)
				iFps, _ := strconv.Atoi(fpsShow)
				f.bars.BarDescribe("merge|"+f.mergeFileName, fmt.Sprintf("[fps:%4d] merging %s ...", iFps, f.mergeFileName))
			}
			if strings.Contains(datas[i], "progress=end") {
				f.bars.BarFinish("merge|" + f.mergeFileName)
				f.mergeFileName = ""
				continue
			}
		}
		data = datas[len(datas)-1]
	}

}
