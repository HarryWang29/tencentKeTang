package project

import (
	"crawler/tencentKeTang/config"
	"github.com/iris-contrib/schema"
	"github.com/pkg/errors"
	"net/url"
	"strings"
)

const (
	ItemsUri = "https://ke.qq.com/cgi-bin/course/get_terms_detail?"
	TokenUri = "https://ke.qq.com/cgi-bin/qcloud/get_token?"
	MediaUri = "https://playvideo.qcloud.com/getplayinfo/v2/1258712167/"
)

type Project struct {
	c      *config.Config
	CID    string `url:"course_id"`
	TaID   string `url:"taid"`
	TermID string `url:"term_id"`
	Type   string `url:"type"`
	VID    string `url:"vid"`
}

func New(c *config.Config) *Project {
	return &Project{c: c}
}

func (p *Project) LoadTaskUrl(taskUrl string) error {
	taskUrl = strings.Replace(taskUrl, "#", "?", 1)
	u, err := url.Parse(taskUrl)
	if err != nil {
		return errors.Wrap(err, "url.Parse")
	}
	values := u.Query()
	err = schema.DecodeQuery(values, p)
	if err != nil {
		return errors.Wrap(err, "schema.DecodeQuery")
	}
	return nil
}
