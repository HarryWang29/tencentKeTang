package project

import (
	"crawler/tencentKeTang/keTang"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"log"
	"strings"
)

func (a *api) DownLoadByIndex(i int64) (err error) {
	if len(a.catalogues) <= int(i) {
		return errors.New("index error")
	}
	catalogue := a.catalogues[i]
	err = a.dealData(catalogue.Data)
	if err != nil {
		return errors.Wrap(err, "dealCatalogue")
	}
	return nil
}

func (a *api) DownLoadByCID(cid string) (err error) {
	list, err := a.GetCatalogue(cid, 0)
	if err != nil {
		return errors.Wrap(err, "getCatalogue")
	}
	for _, catalogue := range list {
		err = a.dealData(catalogue.Data)
		if err != nil {
			log.Printf("dealCatalogue err:%s", err)
			continue
		}
	}
	return nil
}

//递归处理
func (a *api) dealData(data interface{}) error {
	switch v := data.(type) {
	case *keTang.BasicTerm:
		for _, chapter := range v.ChapterInfo {
			err := a.dealData(chapter)
			if err != nil {
				log.Printf("dealChapterInfo err:%s", err)
				continue
			}
		}
	case *keTang.BasicChapter:
		for _, sub := range v.SubInfo {
			err := a.dealData(sub)
			if err != nil {
				log.Printf("dealSubInfo err:%s", err)
				continue
			}
		}
	case *keTang.BasicSub:
		for _, task := range v.TaskInfo {
			err := a.dealData(task)
			if err != nil {
				log.Printf("dealTask err:%s", err)
				continue
			}
		}
	case *keTang.BasicTask:
		ids := a.string2list(v.ResidList)
		for i, id := range ids {
			vodUrl, err := a.getVodUrl(fmt.Sprint(v.Cid), fmt.Sprint(v.TermID), fmt.Sprint(id))
			if err != nil {
				log.Printf("getVodUrl err: %s", err)
				continue
			}
			//下载视频，由于m3u8转mp4主要消耗的是cpu/gpu资源，于是此处不考虑开启并发
			name := v.Name
			if len(ids) > 1 {
				//当出现task中有多个视频文件时，会出现覆盖问题，此处增加序号
				name = fmt.Sprintf("%s%d", name, i+1)
			}
			err = a.downLoad(vodUrl, name)
			if err != nil {
				log.Printf("download err:%s", err)
				continue
			}
		}
	default:
		return errors.New("unknown type")
	}
	return nil
}

func (a *api) getVodUrl(cid, termID, vID string) (vodUrl string, err error) {
	//获取文件token
	ret, err := a.keTang.Token(&keTang.Token{
		TermID: termID,
		FileID: vID,
	})
	if err != nil {
		return "", errors.Wrap(err, "keTang.Token")
	}
	//获取下载连接
	mediaInfo, err := a.keTang.MediaInfo(&keTang.MediaInfo{
		Sign:  ret.Sign,
		T:     ret.T,
		Exper: ret.Exper,
		Us:    ret.Us,
		Vid:   vID,
	})
	if err != nil {
		return "", errors.Wrap(err, "keTang.MediaInfo")
	}

	//拼接视频真实地址
	vodUrl = mediaInfo.VideoInfo.TranscodeList[len(mediaInfo.VideoInfo.TranscodeList)-1].URL
	i := strings.LastIndex(vodUrl, "/")
	vodUrl = vodUrl[:i+1] + "voddrm.token." + a.getMediaToken(cid, termID) + "." + vodUrl[i+1:]
	return vodUrl, nil
}

func (a *api) downLoad(vodUrl, name string) (err error) {
	err = a.ffmpeg.Do(vodUrl, name)
	if err != nil {
		return errors.Wrap(err, "ffmpeg.Do")
	}
	return nil
}

func (a *api) getMediaToken(cID, termID string) string {
	var origin string
	v, ok := a.cookie.Load("uid_type")
	if !ok {
		v = ""
	}
	switch v.(string) {
	case "":
		//发现有没有"uid_type"的情况
		origin = fmt.Sprintf("uin=%s;skey=%s;pskey=%s;plskey=%s;ext=;cid=%s;term_id=%s;vod_type=0",
			gjson.Get(a.getCookieByKey("tdw_data_new_2"), "uin").String(),
			a.getCookieByKey("skey"),
			a.getCookieByKey("p_skey"),
			a.getCookieByKey("p_lskey"),
			cID,
			termID,
		)
	case "0":
		//qq扫码与qq帐号登录都是0
		origin = fmt.Sprintf("uin=%s;skey=%s;pskey=%s;plskey=%s;ext=;uid_type=%s;uid_origin_uid_type=%s;cid=%s;term_id=%s;vod_type=0",
			a.getCookieByKey("uin"),
			a.getCookieByKey("skey"),
			a.getCookieByKey("p_skey"),
			a.getCookieByKey("p_lskey"),
			a.getCookieByKey("uid_type"),
			a.getCookieByKey("uid_origin_uid_type"),
			cID,
			termID,
		)
	case "2":
		//微信扫码登录
		origin = fmt.Sprintf("uin=%s;skey=;pskey=;plskey=;ext=%s;uid_appid=%s;uid_type=2;uid_origin_uid_type=%s;cid=%s;term_id=%s;vod_type=0",
			a.getCookieByKey("uin"),
			a.getCookieByKey("uid_a2"),
			a.getCookieByKey("uid_appid"),
			a.getCookieByKey("uid_origin_uid_type"),
			cID,
			termID,
		)
	}
	return base64.StdEncoding.EncodeToString([]byte(origin))
}
