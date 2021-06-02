package project

import (
	"crawler/tencentKeTang/internal/httplib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/iris-contrib/schema"
	"github.com/pkg/errors"
	"net/url"
	"strings"
)

type MediaInfo struct {
	Sign   string `url:"sign"`
	T      string `url:"t"`
	Exper  int    `url:"exper"`
	Us     string `url:"us"`
	Vid    string `url:"-"`
	Cookie string `url:"-"`
	CID    int    `url:"-"`
	TermID string `url:"-"`
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
	vodUrl = vodUrl[:i+1] + "voddrm.token." + m.getMediaToken() + "." + vodUrl[i+1:]
	return vodUrl, nil
}

func (m *MediaInfo) getMediaToken() string {
	cm := m.cookie2Map()
	origin := fmt.Sprintf("uin=%s;skey=%s;pskey=%s;plskey=%s;ext=;uid_type=%s;uid_origin_uid_type=%s;cid=%d;term_id=%s;vod_type=0",
		cm["uin"],
		cm["skey"],
		cm["p_skey"],
		cm["p_lskey"],
		cm["uid_type"],
		cm["uid_origin_uid_type"],
		m.CID,
		m.TermID,
	)
	return base64.StdEncoding.EncodeToString([]byte(origin))
}

func (m *MediaInfo) cookie2Map() map[string]string {
	list := strings.Split(m.Cookie, "; ")
	cm := make(map[string]string)
	for _, s := range list {
		i := strings.Index(s, "=")
		cm[s[:i]] = s[i+1:]
	}
	return cm
}
