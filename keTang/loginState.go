package keTang

import (
	"fmt"
	"github.com/iris-contrib/schema"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

type LoginStateReq struct {
	Rd  float64 `url:"rd"`
	BKN string  `url:"bkn"`
	R   float64 `url:"r"`
}

type LoginStateResp struct {
	Result struct {
		State int `json:"state"`
	} `json:"result"`
	Retcode int `json:"retcode"`
	Cookies []*http.Cookie
}

func (a *api) LoginState(rd, r float64) (resp *LoginStateResp, err error) {
	req := LoginStateReq{
		Rd:  rd,
		BKN: "",
		R:   r,
	}
	v := url.Values{}
	err = schema.NewEncoder().Encode(req, v)
	if err != nil {
		return nil, errors.Wrap(err, "schema.NewEncoder().Encode")
	}
	cookies, err := a.get(fmt.Sprintf("%s%s", LoginState, v.Encode()), &resp,
		"referer", "https://ke.qq.com/webcourse/index.html",
		"cookie", a.c.Cookie,
	)
	if err != nil {
		return nil, errors.Wrap(err, "get")
	}
	resp.Cookies = cookies
	return resp, nil
}
