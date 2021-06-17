package keTang

import (
	"fmt"
	"github.com/iris-contrib/schema"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

type PtQrShowReq struct {
	AppID    string `url:"appid"`
	E        string `url:"e"`
	L        string `url:"l"`
	S        string `url:"s"`
	D        string `url:"d"`
	V        string `url:"v"`
	T        string `url:"t"`
	DAID     string `url:"daid"`
	Pt3rdAID string `url:"pt_3rd_aid"`
}

func (a *api) PtQrShow() (cookie []*http.Cookie, img []byte, err error) {
	req := &PtQrShowReq{
		AppID:    "716027609",
		E:        "2",
		L:        "M",
		S:        "3",
		D:        "72",
		V:        "4",
		T:        "0.560473424898503",
		DAID:     "383",
		Pt3rdAID: "101487368",
	}
	v := url.Values{}
	err = schema.NewEncoder().Encode(req, v)
	if err != nil {
		return nil, nil, errors.Wrap(err, "schema.NewEncoder().Encode")
	}
	cookie, err = a.get(fmt.Sprintf("%s%s", PtQrShow, v.Encode()), &img,
		"referer", "https://xui.ptlogin2.qq.com/")
	if err != nil {
		return nil, nil, errors.Wrap(err, "get")
	}
	return cookie, img, nil
}
