package keTang

import (
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

func (a *api) Token(t *Token) (ret *TokenResult, err error) {
	t.BKN = a.c.BKN
	t.T = a.c.T
	t.Cookie = a.c.Cookie
	v := url.Values{}
	err = schema.NewEncoder().Encode(t, v)
	if err != nil {
		return nil, errors.Wrap(err, "schema.NewEncoder().Encode")
	}
	resp := &TokenResp{}
	err = a.get(fmt.Sprintf("%s%s", TokenUri, v.Encode()), &resp,
		"referer", "https://ke.qq.com/webcourse/index.html",
		"cookie", t.Cookie)
	if err != nil {
		return nil, errors.Wrap(err, "a.get")
	}

	if resp.Result == nil {
		return nil, errors.New("tokenResp.Result is empty")
	}
	return resp.Result, nil
}
