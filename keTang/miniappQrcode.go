package keTang

import (
	"fmt"
	"github.com/iris-contrib/schema"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

type MiniAppQrcodeReq struct {
	AppID     string  `url:"appId"`
	AutoColor bool    `url:"autoColor"`
	Blue      int     `url:"blue"`
	Green     int     `url:"green"`
	IsHyaLine bool    `url:"isHyaLine"`
	Page      string  `url:"page"`
	Rd        float64 `url:"rd"`
	Red       int     `url:"red"`
	Width     int     `url:"width"`
	BKN       string  `url:"bkn"`
	R         float64 `url:"r"`
}

type MiniAppQrcodeResp struct {
	Result struct {
		QRCode string `json:"QRCode"`
		Expire int    `json:"expire"`
	} `json:"result"`
	Retcode int            `json:"retcode"`
	Cookie  []*http.Cookie `json:"-"`
}

func (a *api) MiniAppQrcode(rd, r float64) (resp *MiniAppQrcodeResp, err error) {
	req := MiniAppQrcodeReq{
		AppID:     "wxa2c453d902cdd452",
		AutoColor: false,
		Blue:      0,
		Green:     0,
		IsHyaLine: false,
		Page:      "pages%2Flogin%2Flogin",
		Rd:        rd,
		Red:       0,
		Width:     200,
		BKN:       "",
		R:         r,
	}
	v := url.Values{}
	err = schema.NewEncoder().Encode(req, v)
	if err != nil {
		return nil, errors.Wrap(err, "schema.NewEncoder().Encode")
	}
	cookies, err := a.get(fmt.Sprintf("%s%s", MiniAppQrcode, v.Encode()), &resp,
		"referer", "https://ke.qq.com/webcourse/index.html",
	)
	if err != nil {
		return nil, errors.Wrap(err, "a.get")
	}
	resp.Cookie = cookies
	return resp, nil
}
