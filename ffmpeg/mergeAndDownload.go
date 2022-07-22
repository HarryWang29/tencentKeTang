package ffmpeg

import (
	"bytes"
	"crawler/tencentKeTang/internal/httplib"
	"fmt"
	"github.com/pkg/errors"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func (f *Ffmpeg) downloadTs(vodUrl string, bitrate int, mp4Path string) error {
	dir := filepath.Dir(mp4Path)
	fileName := filepath.Base(mp4Path)
	m3u8Dir := filepath.Join(dir, fileName+"m3u8")

	err := os.MkdirAll(m3u8Dir, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "os.MkdirAll(%s)", m3u8Dir)
	}

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

	keys := make(map[string]string)
	tsCount := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "#EXT-X-KEY") {
			reg1 := regexp.MustCompile(`URI="(.*)"`)
			if reg1 == nil {
				return errors.Wrap(err, "regexp.MustCompile")
			}
			//根据规则提取关键信息
			result1 := reg1.FindAllStringSubmatch(line, -1)
			if len(result1) == 0 {
				return errors.Wrap(err, "regexp.FindAllStringSubmatch")
			}
			keyUrl := result1[0][1]
			keyFileName := ""
			if v, ok := keys[keyUrl]; ok {
				keyFileName = v
			} else {
				keyFileName = fmt.Sprintf("key%d", len(keys))
				err := httplib.Get(keyUrl).ToFile(filepath.Join(m3u8Dir, keyFileName))
				if err != nil {
					return errors.Wrap(err, "httplib.Get")
				}
				keys[keyUrl] = keyFileName
			}
			line = reg1.ReplaceAllString(line, fmt.Sprintf(`URI="./%s"`, keyFileName))
		} else if strings.HasPrefix(line, "#") {
		} else {
			parm := strings.Split(line, "?")
			loadUrl.Path = filepath.Join(filepath.Dir(loadUrl.Path), parm[0])
			loadUrl.RawQuery = parm[1]

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
	err := httplib.Get(tsUrl).ToFile(savePath)
	if err != nil {
		return errors.Wrap(err, "httplib.Get")
	}
	return nil
}

func (f *Ffmpeg) merge(src, dst string, bitrate int) error {
	args := []string{
		"-allowed_extensions", "ALL",
		"-i", src,
	}
	if len(f.ffmpegParams) != 0 {
		args = append(args, f.ffmpegParams...)
		args = append(args, "-b:v", fmt.Sprint(bitrate))
	}
	//args = append(args, "-progress", fmt.Sprintf(`tcp://%s`, sockFileName))
	args = append(args, dst, "-y")
	cmd := exec.Command(f.ffmpegExec, args...)
	buf := bytes.NewBuffer(nil)
	cmd.Stderr = buf
	err := cmd.Run()
	if err != nil {
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
