package ffmpeg

import (
	"fmt"
	"github.com/pkg/errors"
	"os/exec"
)

func (f *Ffmpeg) mergeAndDownload(vodUrl, name, sockFileName string) error {
	args := []string{"-i", vodUrl}
	if len(f.ffmpegParams) != 0 {
		args = append(args, f.ffmpegParams...)
	}
	args = append(args, "-progress", fmt.Sprintf("unix://%s", sockFileName))
	args = append(args, name, "-y")
	cmd := exec.Command(f.ffmpegExec, args...)
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "exec.Run")
	}

	return nil
}
