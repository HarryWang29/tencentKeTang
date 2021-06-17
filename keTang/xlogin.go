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
	DaID           string `url:"daid"`
	Style          string `url:"style"`
	Theme          string `url:"theme"`
	LoginText      string `url:"login_text"`
	HideTitleBar   string `url:"hide_title_bar"`
	HideBorder     string `url:"hide_border"`
	Target         string `url:"target"`
	SUrl           string `url:"s_url"`
	Pt3rdAID       string `url:"pt_3rd_aid"`
	PtFeedbackLink string `url:"pt_feedback_link"`
}

func (a *api) XLogin() (cookie []*http.Cookie, err error) {
	req := XLoginReq{
		AppID:          "716027609",
		DaID:           "383",
		Style:          "33",
		Theme:          "2",
		LoginText:      "授权并登录",
		HideTitleBar:   "1",
		HideBorder:     "1",
		Target:         "self",
		SUrl:           "https://graph.qq.com/oauth2.0/login_jump",
		Pt3rdAID:       "101487368",
		PtFeedbackLink: "https://support.qq.com/products/77942?customInfo=www.qq.com.appid101487368",
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
