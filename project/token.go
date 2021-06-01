package project

import (
	"crawler/tencentKeTang/internal/httplib"
	"encoding/json"
	"fmt"
	"github.com/iris-contrib/schema"
	"github.com/pkg/errors"
	"net/url"
)

type Token struct {
	TermID string  `url:"term_id"`
	FileID string  `url:"fileId"`
	BKN    int64   `url:"bkn"`
	T      float32 `url:"t"`
	Cookie string  `url:"-"`
}

type TokenResult struct {
	Sign  string `json:"sign"`
	T     string `json:"t"`
	Exper int    `json:"exper"`
	Us    string `json:"us"`
}

type TokenResp struct {
	Result  *TokenResult `json:"result"`
	Retcode int          `json:"retcode"`
}

func (t *Token) Get() (ret *TokenResult, err error) {
	v := url.Values{}
	err = schema.NewEncoder().Encode(t, v)
	if err != nil {
		return nil, errors.Wrap(err, "schema.NewEncoder().Encode")
	}
	req := httplib.Get(nil, fmt.Sprintf("%s%s", TokenUri, v.Encode()))
	req.Header("referer", "https://ke.qq.com/webcourse/index.html")
	req.Header("cookie", t.Cookie)
	body, err := req.Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "httplib.Get.Bytes")
	}

	resp := &TokenResp{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	if resp.Result == nil {
		return nil, errors.New("tokenResp.Result is empty")
	}
	return resp.Result, nil
}
