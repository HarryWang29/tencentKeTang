package keTang

import (
	"fmt"
	"github.com/iris-contrib/schema"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

type QQLoginReq struct {
	BKN int64   `url:"bkn"`
	R   float32 `url:"r"`
}

func (a *api) QQLogin() (cookie []*http.Cookie, err error) {
	req := QQLoginReq{
		BKN: a.c.BKN,
		R:   a.c.R,
	}
	v := url.Values{}
	err = schema.NewEncoder().Encode(req, v)
	if err != nil {
		return nil, errors.Wrap(err, "schema.NewEncoder().Encode")
	}
	cookies, err := a.get(fmt.Sprintf("%s%s", QQLoginUri, v.Encode()), nil,
		"origin", "https://ke.qq.com",
		"authority", "ke.qq.com",
		"user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36",
		"referer", "https://ke.qq.com/user/login-proxy/login-proxy.html?_wv=2147487745&type=1&redirect_url=https%3A%2F%2Fke.qq.com%2Fcourse%2F338702%23term_id%3D100402399",
		"sec-ch-ua", "\".Not/A)Brand\";v=\"99\", \"Google Chrome\";v=\"103\", \"Chromium\";v=\"103\"",
		"sec-ch-ua-mobile", "?0",
		"sec-ch-ua-platform", "\"macOS\"",
		"sec-fetch-dest", "empty",
		"sec-fetch-mode", "cors",
		"sec-fetch-site", "same-origin",
		"cookie", a.c.Cookie,
	)
	if err != nil {
		return nil, errors.Wrap(err, "get")
	}
	return cookies, nil
}
