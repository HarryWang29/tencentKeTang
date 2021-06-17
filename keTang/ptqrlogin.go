package keTang

import (
	"fmt"
	"github.com/iris-contrib/schema"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

type PtQrLoginReq struct {
	U1             string `url:"u1"`
	PtQrToken      string `url:"ptqrtoken"`
	PtRedirect     string `url:"ptredirect"`
	H              string `url:"h"`
	T              string `url:"t"`
	G              string `url:"g"`
	FromUi         string `url:"from_ui"`
	PtLang         string `url:"ptlang"`
	Action         string `url:"action"`
	JsVer          string `url:"js_ver"`
	JsType         string `url:"js_type"`
	LoginSig       string `url:"login_sig"`
	PtUiStyle      string `url:"pt_uistyle"`
	LowLoginEnable string `url:"low_login_enable"`
	AID            string `url:"aid"`
	DAID           string `url:"daid"`
	//PtDrvs         string `url:"ptdrvs"`
	//SID            string `url:"sid"`
	Pt3rdAID string `url:"pt_3rd_aid"`
}

type PtQrLoginResp struct {
	Code        string
	RedirectUrl string
	Msg         string
	NickName    string
	Cookie      []*http.Cookie
}

func (a *api) PtQrLogin(ptQrToken int64, loginSig, ptDrvs, sID string) (*PtQrLoginResp, error) {
	req := &PtQrLoginReq{
		U1:             "https://graph.qq.com/oauth2.0/login_jump",
		PtQrToken:      fmt.Sprint(ptQrToken),
		PtRedirect:     "0",
		H:              "1",
		T:              "1",
		G:              "1",
		FromUi:         "1",
		PtLang:         "2052",
		Action:         fmt.Sprintf("0-0-%d000", time.Now().Unix()),
		JsVer:          "21061713",
		JsType:         "1",
		LoginSig:       loginSig,
		PtUiStyle:      "40",
		LowLoginEnable: "1",
		AID:            "716027609",
		DAID:           "383",
		Pt3rdAID:       "101487368",
	}
	v := url.Values{}
	err := schema.NewEncoder().Encode(req, v)
	if err != nil {
		return nil, errors.Wrap(err, "schema.NewEncoder().Encode")
	}
	respStr := ""
	resp := &PtQrLoginResp{}
	resp.Cookie, err = a.get(fmt.Sprintf("%s%s&", PtQrLogin, v.Encode()), &respStr,
		"referer", "https://xui.ptlogin2.qq.com/",
		"cookie", a.c.Cookie)
	if err != nil {
		return nil, errors.Wrap(err, "get")
	}
	re := regexp.MustCompile(`ptuiCB\('(.*)','0','(.*)','0','(.*)', '(.*)'\)`)
	ret := re.FindStringSubmatch(respStr)
	if len(ret) != 5 {
		return nil, errors.New("re.FindStringSubmatch error")
	}
	resp.Code = ret[1]
	resp.RedirectUrl = ret[2]
	resp.Msg = ret[3]
	resp.NickName = ret[4]
	log.Println(respStr)
	return resp, nil
}
