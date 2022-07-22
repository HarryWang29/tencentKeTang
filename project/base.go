package project

import (
	"crawler/tencentKeTang/ffmpeg"
	"crawler/tencentKeTang/keTang"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	ErrExpire = errors.New("expire")
)

type Api interface {
	GetCatalogue(cid string, tid int64) (list []*Catalogue, err error)
	DownLoadByIndex(i int64) (err error)
	DownLoadByCID(cid string) (err error)
	SetCookie(cookie string)
	WxQRLogin() (err error)
	QQQRLogin() (nickName string, err error)
}

type api struct {
	keTang        keTang.Api
	ffmpeg        *ffmpeg.Ffmpeg
	catalogues    []*Catalogue
	cookie        sync.Map
	catalogueName string
	vodUrlMap     sync.Map
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

func (a *api) SetCookie(cookie string) {
	a.cookie2Map(cookie)
}

func (a *api) AddCookie(cookies []*http.Cookie) {
	for _, cookie := range cookies {
		a.cookie.Store(cookie.Name, cookie.Value)
	}
	a.keTang.SetCookie(a.GetCookies())
}

func (a *api) GetCookies() []*http.Cookie {
	cookies := make([]*http.Cookie, 0)
	a.cookie.Range(func(key, value interface{}) bool {
		c := http.Cookie{
			Name:  key.(string),
			Value: value.(string),
		}
		cookies = append(cookies, &c)
		return true
	})
	return cookies
}

func (a *api) cookie2Map(cookie string) {
	if cookie == "" {
		return
	}
	list := strings.Split(cookie, "; ")
	for _, s := range list {
		i := strings.Index(s, "=")
		a.cookie.Store(s[:i], s[i+1:])
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

func (a *api) getCookieByKey(s string) string {
	v, ok := a.cookie.Load(s)
	if !ok {
		return ""
	}
	return v.(string)
}
