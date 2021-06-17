package keTang

import (
	"crawler/tencentKeTang/internal/httplib"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

const (
	ItemsUri      = "https://ke.qq.com/cgi-bin/course/get_terms_detail?"
	TokenUri      = "https://ke.qq.com/cgi-bin/qcloud/get_token?"
	MediaUri      = "https://playvideo.qcloud.com/getplayinfo/v2/1258712167/"
	InfoUri       = "https://ke.qq.com/cgi-bin/identity/info?"
	BasicInfoUri  = "https://ke.qq.com/cgi-bin/course/basic_info?"
	MiniAppQrcode = "https://ke.qq.com/cgi-proxy/get_miniapp_qrcode?"
	LoginState    = "https://ke.qq.com/cgi-proxy/get_login_state?"
	A2Login       = "https://ke.qq.com/cgi-proxy/account_login/a2_login?"
	XLogin        = "https://xui.ptlogin2.qq.com/cgi-bin/xlogin?"
	PtQrShow      = "https://ssl.ptlogin2.qq.com/ptqrshow?"
	PtQrLogin     = "https://ssl.ptlogin2.qq.com/ptqrlogin?"
	Check         = "https://ssl.ptlogin2.qq.com/check?"
)

type Config struct {
	Cookie string  `yaml:"cookie"`
	BKN    int64   `yaml:"bkn"`
	T      float32 `yaml:"t"`
	R      float32 `yaml:"r"`
}

type Api interface {
	SetCookie(cookies []*http.Cookie)
	Get(i *Items) (resp *ItemsResp, err error)
	Info() (result *InfoResult, err error)
	BasicInfo(cid string) (resp *BasicInfoResp, err error)
	Token(t *Token) (ret *TokenResult, err error)
	MediaInfo(m *MediaInfo) (info *MediaInfoResp, err error)
	MiniAppQrcode(rd, r float64) (resp *MiniAppQrcodeResp, err error)
	LoginState(rd, r float64) (resp *LoginStateResp, err error)
	A2Login(bkn string, r float64) (resp *A2LoginResp, err error)
	XLogin() (cookie []*http.Cookie, err error)
	PtQrShow() (cookie []*http.Cookie, img []byte, err error)
	PtQrLogin(ptQrToken int64, loginSig, ptDrvs, sID string) (*PtQrLoginResp, error)
}

type api struct {
	c *Config
}

func New(c *Config) Api {
	return &api{c: c}
}

func (a *api) post(url string, resp interface{}, headers ...string) (cookies []*http.Cookie, err error) {
	req := httplib.Post(url)
	if len(headers)%2 != 0 {
		return nil, errors.New("headers error")
	}
	for i := 0; i < len(headers)-1; i += 2 {
		req.Header(headers[i], headers[i+1])
	}
	body, err := req.Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "httplib.Get.Bytes")
	}

	httpResp, err := req.Response()
	if err != nil {
		return nil, errors.Wrap(err, "httplib.Response")
	}
	cookies = httpResp.Cookies()
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "json.Unmarshal,resp:%s", string(body))
	}
	return cookies, nil
}

func (a *api) get(url string, resp interface{}, headers ...string) (cookies []*http.Cookie, err error) {
	req := httplib.Get(url)
	if len(headers)%2 != 0 {
		return nil, errors.New("headers error")
	}
	for i := 0; i < len(headers)-1; i += 2 {
		req.Header(headers[i], headers[i+1])
	}
	body, err := req.Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "httplib.Get.Bytes")
	}

	httpResp, err := req.Response()
	if err != nil {
		return nil, errors.Wrap(err, "httplib.Response")
	}
	cookies = httpResp.Cookies()
	if resp == nil {
		return cookies, nil
	}
	switch t := resp.(type) {
	case *[]byte:
		*t = make([]byte, 0)
		*t = body
	case *string:
		*t = string(body)
	default:
		err = json.Unmarshal(body, &resp)
		if err != nil {
			return nil, errors.Wrapf(err, "json.Unmarshal,resp:%s", string(body))
		}
	}
	return cookies, nil
}

func (a *api) SetCookie(cookies []*http.Cookie) {
	a.c.Cookie = ""
	for _, cookie := range cookies {
		a.c.Cookie += fmt.Sprintf("%s=%s; ", cookie.Name, cookie.Value)
	}
}
