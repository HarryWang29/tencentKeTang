package project

import (
	"crawler/tencentKeTang/ffmpeg"
	"crawler/tencentKeTang/keTang"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
)

type Api interface {
	GetCatalogue(cid string, tid int64) (list []*Catalogue, err error)
	DownLoadByIndex(i int64) (err error)
	DownLoadByCID(cid string) (err error)
}

type api struct {
	keTang     keTang.Api
	ffmpeg     *ffmpeg.Ffmpeg
	catalogues []*Catalogue
	cookie     map[string]string
}

func New(kt keTang.Api, f *ffmpeg.Ffmpeg, cookie string) Api {
	a := &api{
		keTang:     kt,
		ffmpeg:     f,
		catalogues: make([]*Catalogue, 0),
	}
	a.cookie2Map(cookie)
	return a
}

func (a *api) cookie2Map(cookie string) {
	list := strings.Split(cookie, "; ")
	a.cookie = make(map[string]string)
	for _, s := range list {
		i := strings.Index(s, "=")
		a.cookie[s[:i]] = s[i+1:]
	}
}

func (a *api) string2list(s string) []int64 {
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
