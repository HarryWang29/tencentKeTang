package ffmpeg

import (
	"bytes"
	"crawler/tencentKeTang/internal/httplib"
	"crawler/tencentKeTang/util"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func (f *Ffmpeg) downloadTs(vodUrl, dk string, bitrate int, mp4Path string) error {
	dir := filepath.Dir(mp4Path)
	fileName := filepath.Base(mp4Path)
	m3u8Dir := filepath.Join(dir, fileName+"m3u8")

	err := os.MkdirAll(m3u8Dir, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "os.MkdirAll(%s)", m3u8Dir)
	}
	key, err := base64.StdEncoding.DecodeString(dk)
	if err != nil {
		return errors.Wrap(err, "base64.StdEncoding.DecodeString")
	}
	keyFile, err := os.OpenFile(filepath.Join(m3u8Dir, "key"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return errors.Wrapf(err, "os.OpenFile(%s)", filepath.Join(m3u8Dir, "index.m3u8"))
	}
	_, err = keyFile.Write(key)
	if err != nil {
		return errors.Wrap(err, "keyFile.Write")
	}
	_ = keyFile.Close()

	m3u8, err := httplib.Get(vodUrl).String()
	if err != nil {
		return errors.Wrap(err, "httplib.Get")
	}
	loadUrl, err := url.Parse(vodUrl)
	if err != nil {
		return errors.Wrap(err, "url.Parse")
	}
	loadUrl.RawQuery = ""

	m3u8Path := filepath.Join(m3u8Dir, fileName)
	m3u8Path += ".m3u8"
	lines := strings.Split(m3u8, "\n")
	file, err := os.OpenFile(m3u8Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return errors.Wrap(err, "os.OpenFile")
	}
	//先用lines设置max
	f.addDownloadBar(fileName, len(lines))

	tsCount := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "#EXT-X-KEY") {
			reg1 := regexp.MustCompile(`URI="(.*)"`)
			if reg1 == nil {
				return errors.Wrap(err, "regexp.MustCompile")
			}
			line = reg1.ReplaceAllString(line, fmt.Sprintf(`URI="./%s"`, "key"))
		} else if strings.HasPrefix(line, "#") {
		} else {
			parm := strings.Split(line, "?")
			urlPath := strings.Split(loadUrl.Path, "/")
			loadUrl.Path = strings.Join(append(urlPath[:len(urlPath)-1], parm[0]), "/")
			loadUrl.RawQuery = parm[1]

			line = util.ReplaceName(line)
			downloadPath := filepath.Join(m3u8Dir, line)
			task := &task{
				vodUrl:   vodUrl,
				tsUrl:    loadUrl.String(),
				savePath: downloadPath,
				fileName: fileName,
				m3u8Dir:  m3u8Dir,
				m3u8Path: m3u8Path,
				bitrate:  bitrate,
			}
			f.asyncDownload(task)
			tsCount++
		}
		_, _ = file.WriteString(line + "\n")
		if line == "#EXT-X-ENDLIST" {
			break
		}
	}
	//max设置完成后更新
	f.bars.BarChangeMax("down|"+fileName, tsCount)
	_ = file.Close()

	return nil
}

func (f *Ffmpeg) asyncDownload(t *task) {
	f.workerChannel <- t
	count, ok := f.taskMap.Load(t.vodUrl)
	if !ok {
		f.taskMap.Store(t.vodUrl, 1)
	} else {
		f.taskMap.Store(t.vodUrl, count.(int)+1)
	}
}

func (f *Ffmpeg) doDownloadTs(tsUrl, savePath string) error {
	resp, err := f.httpclient.Get(tsUrl)
	saveFile, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer saveFile.Close()

	if resp.Body == nil {
		return nil
	}
	defer resp.Body.Close()
	_, err = io.Copy(saveFile, resp.Body)
	if err != nil {
		return errors.Wrap(err, "httplib.Get")
	}
	return nil
}

func (f *Ffmpeg) merge(src, dst string, bitrate int) error {
	args := []string{
		"-allowed_extensions", "ALL",
		"-hwaccel", "auto",
		"-i", src,
	}
	if len(f.ffmpegParams) != 0 {
		args = append(args, f.ffmpegParams...)
		args = append(args, "-b:v", fmt.Sprint(bitrate))
	}
	args = append(args, "-progress", fmt.Sprintf(`tcp://%s`, f.address))
	args = append(args, dst, "-y")
	cmd := exec.Command(f.ffmpegExec, args...)
	buf := bytes.NewBuffer(nil)
	cmd.Stderr = buf
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%s\n", buf.String())
		return errors.Wrapf(err, "exec.Run err: %s", buf.String())
	}

	return nil
}

func (f *Ffmpeg) mergeAndDownload(vodUrl, name, sockFileName string) error {
	args := []string{"-i", fmt.Sprintf(`%s`, vodUrl)}
	if len(f.ffmpegParams) != 0 {
		args = append(args, f.ffmpegParams...)
	}
	args = append(args, "-progress", fmt.Sprintf(`tcp://%s`, sockFileName))
	args = append(args, name, "-y")
	cmd := exec.Command(f.ffmpegExec, args...)
	buf := bytes.NewBuffer(nil)
	cmd.Stderr = buf
	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "exec.Run err: %s", buf.String())
	}

	return nil
}

func (f *Ffmpeg) checkFfmpeg() (err error) {
	cmd := exec.Command(f.ffmpegExec, "-version")
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "exec.Run")
	}
	return nil
}
