package ffmpeg

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type task struct {
	vodUrl   string
	tsUrl    string
	m3u8Dir  string
	savePath string
	fileName string
	m3u8Path string
	bitrate  int
}

func (f *Ffmpeg) taskProcess() {
	for {
		task := <-f.workerChannel
		_ = f.doDownloadTs(task.tsUrl, task.savePath)
		finishCount, ok := f.finishMap.Load(task.vodUrl)
		if !ok {
			finishCount = new(int)
			finishCount = 1
			f.finishMap.Store(task.vodUrl, finishCount)
		} else {
			finishCount = finishCount.(int) + 1
			f.finishMap.Store(task.vodUrl, finishCount)
		}
		taskCount, ok := f.taskMap.Load(task.vodUrl)
		if !ok {
			continue
		}
		if taskCount.(int) == finishCount.(int) {
			f.finishChannel <- task
		}
	}
}

func (f *Ffmpeg) finish() {
	for {
		task := <-f.finishChannel
		task.fileName = strings.ReplaceAll(task.fileName, filepath.Ext(task.fileName), ".mp4")
		err := f.merge(task.m3u8Path, filepath.Join(task.m3u8Dir, "..", task.fileName), task.bitrate)
		if err != nil {
			log.Printf("merge error: %v", err)
			continue
		}
		err = os.RemoveAll(task.m3u8Dir)
		if err != nil {
			log.Printf("remove error: %v", err)
			continue
		}
	}
}
