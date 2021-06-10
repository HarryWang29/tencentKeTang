package keTang

import (
	"fmt"
	"github.com/iris-contrib/schema"
	"github.com/pkg/errors"
	"net/url"
)

type info struct {
	BKN int64   `url:"bkn"`
	R   float32 `url:"r"`
}

type InfoResult struct {
	UID       interface{} `json:"uid"`
	RoleType  int         `json:"role_type"`
	UIDType   int         `json:"uid_type"`
	RoleInfo  interface{} `json:"role_info"`
	NickName  string      `json:"nick_name"`
	AidType   interface{} `json:"aid_type"`
	IsCreator int         `json:"is_creator"`
	FaceURL   string      `json:"face_url"`
	Aid       int         `json:"aid"`
	IsLogin   int         `json:"is_login"`
	IsLogout  int         `json:"is_logout"`
}

type InfoResp struct {
	Result  *InfoResult `json:"result"`
	Retcode int         `json:"retcode"`
}

func (a *api) Info() (result *InfoResult, err error) {
	i := &info{
		BKN: a.c.BKN,
		R:   a.c.R,
	}
	v := url.Values{}
	err = schema.NewEncoder().Encode(i, v)
	if err != nil {
		return nil, errors.Wrap(err, "schema.NewEncoder().Encode")
	}
	resp := &InfoResp{}
	_, err = a.get(fmt.Sprintf("%s%s", InfoUri, v.Encode()), &resp,
		"referer", "https://ke.qq.com/webcourse/index.html",
		"cookie", a.c.Cookie)
	if err != nil {
		return nil, errors.Wrap(err, "a.get")
	}
	return resp.Result, nil
}
