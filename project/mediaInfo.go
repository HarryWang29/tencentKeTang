package project

import (
	"crawler/tencentKeTang/internal/httplib"
	"encoding/json"
	"fmt"
	"github.com/iris-contrib/schema"
	"github.com/pkg/errors"
	"net/url"
	"strings"
)

type MediaInfo struct {
	Sign  string `url:"sign"`
	T     string `url:"t"`
	Exper int    `url:"exper"`
	Us    string `url:"us"`
	Vid   string `url:"-"`
}

type MediaInfoResp struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	RequestID  string `json:"requestId"`
	PlayerInfo struct {
		PlayerID                   string `json:"playerId"`
		Name                       string `json:"name"`
		DefaultVideoClassification string `json:"defaultVideoClassification"`
		VideoClassification        []struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			DefinitionList []int  `json:"definitionList"`
		} `json:"videoClassification"`
		LogoLocation string `json:"logoLocation"`
		LogoPic      string `json:"logoPic"`
		LogoURL      string `json:"logoUrl"`
		PatchInfo    []struct {
			Location int    `json:"location"`
			Link     string `json:"link"`
			Type     string `json:"type"`
			URL      string `json:"url"`
		} `json:"patchInfo"`
	} `json:"playerInfo"`
	CoverInfo struct {
		CoverURL string `json:"coverUrl"`
	} `json:"coverInfo"`
	VideoInfo struct {
		BasicInfo struct {
			Name        string        `json:"name"`
			Description string        `json:"description"`
			Tags        []interface{} `json:"tags"`
		} `json:"basicInfo"`
		Drm struct {
			Definition int `json:"definition"`
		} `json:"drm"`
		MasterPlayList struct {
			IdrAligned bool   `json:"idrAligned"`
			URL        string `json:"url"`
			Definition int    `json:"definition"`
			Md5        string `json:"md5"`
		} `json:"masterPlayList"`
		TranscodeList []struct {
			URL             string  `json:"url"`
			Definition      int     `json:"definition"`
			Duration        int     `json:"duration"`
			FloatDuration   float64 `json:"floatDuration"`
			Size            int     `json:"size"`
			TotalSize       int     `json:"totalSize"`
			Bitrate         int     `json:"bitrate"`
			Height          int     `json:"height"`
			Width           int     `json:"width"`
			Container       string  `json:"container"`
			Md5             string  `json:"md5"`
			VideoStreamList []struct {
				Bitrate int    `json:"bitrate"`
				Codec   string `json:"codec"`
				Fps     int    `json:"fps"`
				Height  int    `json:"height"`
				Width   int    `json:"width"`
			} `json:"videoStreamList"`
			AudioStreamList []struct {
				Bitrate      int    `json:"bitrate"`
				Codec        string `json:"codec"`
				SamplingRate int    `json:"samplingRate"`
			} `json:"audioStreamList"`
			TemplateName string `json:"templateName"`
		} `json:"transcodeList"`
	} `json:"videoInfo"`
}

func (m *MediaInfo) Get() (vodUrl string, err error) {
	v := url.Values{}
	err = schema.NewEncoder().Encode(m, v)
	if err != nil {
		return "", errors.Wrap(err, "schema.NewEncoder().Encode")
	}
	req := httplib.Get(fmt.Sprintf("%s%s?%s", MediaUri, m.Vid, v.Encode()))
	body, err := req.Bytes()
	if err != nil {
		return "", errors.Wrap(err, "httplib.Get.Bytes")
	}

	resp := &MediaInfoResp{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", errors.Wrap(err, "json.Unmarshal")
	}
	vodUrl = resp.VideoInfo.TranscodeList[len(resp.VideoInfo.TranscodeList)-1].URL
	i := strings.LastIndex(vodUrl, "/")
	vodUrl = vodUrl[:i+1] + "voddrm.token.dWluPTEwNjEyNjUwNjI7c2tleT1AZktsQkN6RGFGO3Bza2V5PVkqdXF4MDZ5cUtwWVN0YzNsM2R2WEFrVlp1QjJ0UkNtLTVEem5IVlp1VWtfO3Bsc2tleT0wMDA0MDAwMGVhMTk4YzIwYWM2MjYwNjllMmYxMmM2YTNiMzFjMTIyZTkyNWFjM2RmNjQ5YjFkYzM5ODM1YTBkOTkyZjZiNzVjYTBkYjg4YmFmOTBlNjA2O2V4dD07dWlkX3R5cGU9MDt1aWRfb3JpZ2luX3VpZF90eXBlPTA7Y2lkPTMxMzI4MTU7dGVybV9pZD0xMDMyNTYxOTQ7dm9kX3R5cGU9MA==." + vodUrl[i+1:]
	return vodUrl, nil
}
