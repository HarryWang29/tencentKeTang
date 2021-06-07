package keTang

import (
	"crawler/tencentKeTang/internal/httplib"
	"encoding/json"
	"github.com/pkg/errors"
)

const (
	ItemsUri     = "https://ke.qq.com/cgi-bin/course/get_terms_detail?"
	TokenUri     = "https://ke.qq.com/cgi-bin/qcloud/get_token?"
	MediaUri     = "https://playvideo.qcloud.com/getplayinfo/v2/1258712167/"
	InfoUri      = "https://ke.qq.com/cgi-bin/identity/info?"
	BasicInfoUri = "https://ke.qq.com/cgi-bin/course/basic_info?"
)

type Config struct {
	Cookie string  `yaml:"cookie"`
	BKN    int64   `yaml:"bkn"`
	T      float32 `yaml:"t"`
	R      float32 `yaml:"r"`
}

type Api interface {
	Get(i *Items) (resp *ItemsResp, err error)
	Info() (result *InfoResult, err error)
	BasicInfo(cid string) (resp *BasicInfoResp, err error)
	Token(t *Token) (ret *TokenResult, err error)
	MediaInfo(m *MediaInfo) (info *MediaInfoResp, err error)
}

type api struct {
	c *Config
}

func New(c *Config) Api {
	return &api{c: c}
}

func (a *api) get(url string, resp interface{}, headers ...string) error {
	req := httplib.Get(url)
	if len(headers)%2 != 0 {
		return errors.New("headers error")
	}
	for i := 0; i < len(headers)-1; i += 2 {
		req.Header(headers[i], headers[i+1])
	}
	body, err := req.Bytes()
	if err != nil {
		return errors.Wrap(err, "httplib.Get.Bytes")
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return errors.Wrap(err, "json.Unmarshal")
	}
	return nil
}
