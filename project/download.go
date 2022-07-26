package project

import (
	"crawler/tencentKeTang/keTang"
	"crawler/tencentKeTang/util"
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
	err = a.dealData(catalogue.Data, []string{util.ReplaceName(a.catalogueName)})
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
		err = a.dealData(catalogue.Data, []string{util.ReplaceName(a.catalogueName)})
		if err != nil {
			log.Printf("dealCatalogue err:%s", err)
			continue
		}
	}
	return nil
}

//递归处理
func (a *api) dealData(data interface{}, pathList []string) error {
	if pathList == nil {
		pathList = []string{}
	}
	switch v := data.(type) {
	case *keTang.BasicTerm:
		if v.Name != "" {
			path := fmt.Sprintf("%d.", v.TermNo+1)
			pathList = append(pathList, util.ReplaceName(path+v.Name))
		}
		for _, chapter := range v.ChapterInfo {
			err := a.dealData(chapter, pathList)
			if err != nil {
				log.Printf("dealChapterInfo err:%s", err)
				continue
			}
		}
	case *keTang.BasicChapter:
		if v.Name != "" {
			path := fmt.Sprintf("%d.", v.ChNo+1)
			pathList = append(pathList, util.ReplaceName(path+v.Name))
		}
		for _, sub := range v.SubInfo {
			err := a.dealData(sub, pathList)
			if err != nil {
				log.Printf("dealSubInfo err:%s", err)
				continue
			}
		}
	case *keTang.BasicSub:
		if v.Name != "" {
			path := fmt.Sprintf("%d.", v.SubID+1)
			pathList = append(pathList, util.ReplaceName(path+v.Name))
		}
		for _, task := range v.TaskInfo {
			err := a.dealData(task, pathList)
			if err != nil {
				log.Printf("dealTask err:%s", err)
				continue
			}
		}
	case *keTang.BasicTask:
		ids := a.string2list(v.ResidList)
		for i, id := range ids {
			//简单粗暴的解决方案，每次启动不会下载同样的视频
			_, loaded := a.vodUrlMap.LoadOrStore(fmt.Sprintf("%d|%d|%d", v.Cid, v.TermID, id), "")
			if loaded {
				continue
			}
			vodUrl, dk, bitrate, err := a.getVodUrl(fmt.Sprint(v.Cid), fmt.Sprint(v.TermID), fmt.Sprint(id))
			if err != nil {
				log.Printf("getVodUrl err: %s", err)
				continue
			}
			name := v.Name
			if len(ids) > 1 {
				//当出现task中有多个视频文件时，会出现覆盖问题，此处增加序号
				name = fmt.Sprintf("%s%d", name, i+1)
			}
			pathList = append(pathList, util.ReplaceName(name))
			err = a.downLoad(vodUrl, dk, bitrate, pathList)
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

func (a *api) getVodUrl(cid, termID, vID string) (vodUrl, dk string, bitrate int, err error) {
	//获取下载连接
	info, dk, err := a.keTang.DescribeRecVideo(&keTang.DescribeRecVideoReq{
		CourseID: cid,
		FileID:   vID,
		TermID:   termID,
		VodType:  0,
		Headers: keTang.DescribeRecvVideoHeader{
			SrvAppid: 201,
			CliAppid: "ke",
			Uin:      a.getUin(),
			CliInfo: struct {
				CliPlatform int `json:"cli_platform"`
			}{
				CliPlatform: 3,
			},
		},
	})
	if err != nil {
		return "", "", 0, errors.Wrap(err, "keTang.MediaInfo")
	}

	//拼接视频真实地址
	vodUrl = info.Url
	i := strings.LastIndex(vodUrl, "/")
	vodUrl = vodUrl[:i+1] + "voddrm.token." + a.getMediaToken(cid, termID) + "." + vodUrl[i+1:]
	return vodUrl, dk, info.VideoBitrate, nil
}

func (a *api) downLoad(vodUrl, dk string, bitrate int, path []string) (err error) {
	err = a.ffmpeg.Do(vodUrl, dk, bitrate, path)
	if err != nil {
		return errors.Wrap(err, "ffmpeg.Do")
	}
	return nil
}

func (a *api) getUin() (uin string) {
	v, ok := a.cookie.Load("uid_type")
	if !ok {
		v = ""
	}
	switch v.(string) {
	case "":
		//发现有没有"uid_type"的情况
		uin = gjson.Get(a.getCookieByKey("tdw_data_new_2"), "uin").String()
	default:
		uin = a.getCookieByKey("uin")
	}
	return
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
	default:
		origin = fmt.Sprintf("uin=%s;skey=%s;pskey=%s;plskey=%s;ext=;uid_type=%s;uid_origin_uid_type=%s;uid_origin_auth_type=0;cid=%s;term_id=%s;vod_type=0;platform=3",
			a.getCookieByKey("uin"),
			a.getCookieByKey("skey"),
			a.getCookieByKey("p_skey"),
			a.getCookieByKey("p_lskey"),
			a.getCookieByKey("uid_type"),
			a.getCookieByKey("uid_origin_uid_type"),
			cID,
			termID,
		)
	}
	return base64.StdEncoding.EncodeToString([]byte(origin))
}
