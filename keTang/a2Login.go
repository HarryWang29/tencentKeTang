package keTang

import (
	"fmt"
	"github.com/iris-contrib/schema"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

type A2LoginReq struct {
	BKN string  `url:"bkn"`
	R   float64 `url:"r"`
}

type A2LoginResp struct {
	Cookies []*http.Cookie
}

func (a *api) A2Login(bkn string, r float64) (resp *A2LoginResp, err error) {
	req := A2LoginReq{
		BKN: bkn,
		R:   r,
	}
	v := url.Values{}
	err = schema.NewEncoder().Encode(req, v)
	if err != nil {
		return nil, errors.Wrap(err, "schema.NewEncoder().Encode")
	}
	cookies, err := a.post(fmt.Sprintf("%s%s", A2Login, v.Encode()), &resp,
		"referer", "https://ke.qq.com/user/login-proxy/login-proxy.html?_bid=167&_wv=2147487745&type=2&redirect_url=https%3A%2F%2Fke.qq.com%2Fwebcourse%2Findex.html%23cid%3D398381%26term_id%3D100475149%26taid%3D3385404892255277%26type%3D1024%26vid%3D5285890790625135834",
		"cookie", a.c.Cookie,
	)
	if err != nil {
		return nil, errors.Wrap(err, "post")
	}
	resp.Cookies = cookies
	return resp, nil
}
