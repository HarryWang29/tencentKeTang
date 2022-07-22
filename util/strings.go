package util

import (
	"github.com/tidwall/gjson"
	"path/filepath"
	"regexp"
	"strconv"
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

func String2list(s string) []int64 {
	s = strings.ReplaceAll(s, "&quot;", "\"")
	list := make([]int64, 0)
	if id, err := strconv.ParseInt(s, 10, 64); err == nil {
		list = append(list, id)
	} else {
		arr := gjson.Get(s, "#()#").Array()
		for _, result := range arr {
			list = append(list, result.Int())
		}
	}
	return list
}
