package keTang

import (
	"fmt"
	"github.com/iris-contrib/schema"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

type XLoginReq struct {
	AppID          string `url:"appid"`
	HideCloseIcon  string `url:"hide_close_icon"`
	DaID           string `url:"daid"`
	Target         string `url:"target"`
	SUrl           string `url:"s_url"`
	ProxyUrl       string `url:"proxy_url"`
	LowLoginEnable string `url:"low_login_enable"`
}

func (a *api) XLogin() (cookie []*http.Cookie, err error) {
	req := XLoginReq{
		AppID:          "715030901",
		HideCloseIcon:  "1",
		DaID:           "233",
		Target:         "self",
		SUrl:           "https://ke.qq.com/login_proxy.html",
		ProxyUrl:       "https://ke.qq.com/login_proxy.html",
		LowLoginEnable: "1",
	}
	v := url.Values{}
	err = schema.NewEncoder().Encode(req, v)
	if err != nil {
		return nil, errors.Wrap(err, "schema.NewEncoder().Encode")
	}
	cookies, err := a.get(fmt.Sprintf("%s%s", XLogin, v.Encode()), nil,
		"referer", "https://graph.qq.com/")
	if err != nil {
		return nil, errors.Wrap(err, "get")
	}
	return cookies, nil
}
