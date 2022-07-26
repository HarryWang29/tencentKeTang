package keTang

import (
	"encoding/json"
	"fmt"
	"github.com/iris-contrib/schema"
	"github.com/pkg/errors"
	"net/url"
	"sort"
)

type DescribeRecVideoReq struct {
	CourseID string                  `url:"course_id"`
	FileID   string                  `url:"file_id"`
	Header   string                  `url:"header"`
	TermID   string                  `url:"term_id"`
	VodType  int                     `url:"vod_type"`
	BKN      int64                   `url:"bkn"`
	R        float32                 `url:"r"`
	Headers  DescribeRecvVideoHeader `url:"-"`
}

type DescribeRecvVideoHeader struct {
	SrvAppid int    `json:"srv_appid"`
	CliAppid string `json:"cli_appid"`
	Uin      string `json:"uin"`
	CliInfo  struct {
		CliPlatform int `json:"cli_platform"`
	} `json:"cli_info"`
}

type DescribeRecVideoResp struct {
	Result struct {
		Header struct {
			Code   int    `json:"code"`
			Msg    string `json:"msg"`
			ExtMsg string `json:"ext_msg"`
		} `json:"header"`
		RecVideoInfo *struct {
			FileId         string            `json:"file_id"`
			Dk             string            `json:"dk"`
			Infos          RecVideoInfosList `json:"infos"`
			MasterPlayList string            `json:"master_play_list"`
			PSign          string            `json:"p_sign"`
			Subtitles      []struct {
				Url  string `json:"url"`
				Type string `json:"type"`
			} `json:"subtitles"`
		} `json:"rec_video_info"`
	} `json:"result"`
	Retcode int `json:"retcode"`
}

type RecVideoInfos struct {
	Url          string `json:"url"`
	Duration     int    `json:"duration"`
	Expire       int    `json:"expire"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	AudioCodec   string `json:"audio_codec"`
	VideoCodec   string `json:"video_codec"`
	TemplateId   int    `json:"template_id"`
	IsSpeedHd    int    `json:"is_speed_hd"`
	Size         int    `json:"size"`
	AudioBitrate int    `json:"audio_bitrate"`
	VideoBitrate int    `json:"video_bitrate"`
	SizeByte     string `json:"size_byte"`
	TsDecodeIv   string `json:"ts_decode_iv"`
	TsList       []struct {
		Url      string  `json:"url"`
		Duration float32 `json:"duration"`
	} `json:"ts_list"`
}

type RecVideoInfosList []*RecVideoInfos

func (r RecVideoInfosList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RecVideoInfosList) Len() int {
	return len(r)
}

func (r RecVideoInfosList) Less(i, j int) bool {
	return r[i].VideoBitrate > r[j].VideoBitrate
}

func (a *api) DescribeRecVideo(req *DescribeRecVideoReq) (info *RecVideoInfos, dk string, err error) {
	req.BKN = a.c.BKN
	req.R = a.c.R
	bs, _ := json.Marshal(req.Headers)
	req.Header = string(bs)
	v := url.Values{}
	err = schema.NewEncoder().Encode(req, v)
	if err != nil {
		return nil, "", errors.Wrap(err, "schema.NewEncoder().Encode")
	}
	resp := &DescribeRecVideoResp{}
	_, err = a.get(fmt.Sprintf("%s%s", DescribeRecVideoUri, v.Encode()), &resp,
		"referer", "https://ke.qq.com/webcourse/index.html",
		"cookie", a.c.Cookie)
	if err != nil {
		return nil, "", errors.Wrap(err, "a.get")
	}

	if resp.Result.Header.Code != 0 {
		return nil, "", errors.New(resp.Result.Header.Msg)
	}
	if resp.Result.RecVideoInfo == nil {
		return nil, "", errors.New("describeRecVideoResp.Result.RecVideoInfo is empty")
	}
	if len(resp.Result.RecVideoInfo.Infos) == 0 {
		return nil, "", errors.New("describeRecVideoResp.Result.RecVideoInfo.Infos is empty")
	}
	sort.Sort(resp.Result.RecVideoInfo.Infos)
	return resp.Result.RecVideoInfo.Infos[0], resp.Result.RecVideoInfo.Dk, nil
}
