package keTang

import (
	"fmt"
	"github.com/iris-contrib/schema"
	"github.com/pkg/errors"
	"net/url"
)

type MediaInfo struct {
	Sign   string `url:"sign"`
	T      string `url:"t"`
	Exper  int    `url:"exper"`
	Us     string `url:"us"`
	Cookie string `url:"-"`
	Vid    string `url:"-"`
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

func (a *api) MediaInfo(m *MediaInfo) (info *MediaInfoResp, err error) {
	m.Cookie = a.c.Cookie
	v := url.Values{}
	err = schema.NewEncoder().Encode(m, v)
	if err != nil {
		return nil, errors.Wrap(err, "schema.NewEncoder().Encode")
	}
	_, err = a.get(fmt.Sprintf("%s%s?%s", MediaUri, m.Vid, v.Encode()), &info)
	if err != nil {
		return nil, errors.Wrap(err, "get")
	}
	return info, nil
}
