package util

import (
	"path/filepath"
	"regexp"
	"strings"
)

func PathJoin(paths ...string) string {
	//替换名称中出现的特殊字符
	re, _ := regexp.Compile("[?、\\\\/*\"<>|]")
	names := make([]string, 0, len(paths))
	for _, path := range paths {
		path = filepath.Clean(path)
		path = re.ReplaceAllString(path, "")
		//去除两侧空格
		path = strings.TrimSpace(path)
		names = append(names, path)
	}
	return filepath.Join(names...)
}
