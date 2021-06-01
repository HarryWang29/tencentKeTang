package project

import (
	"crawler/tencentKeTang/ffmpeg"
	"fmt"
	"github.com/pkg/errors"
	"log"
)

func (p *Project) Do(taskUrl string) error {
	//解析url参数
	err := p.LoadTaskUrl(taskUrl)
	if err != nil {
		return errors.Wrap(err, "p.LoadTaskUrl")
	}
	//加载目录
	items := &Items{
		CID:        p.CID,
		TermIDList: fmt.Sprintf("[%s]", p.TermID),
		BKN:        p.c.Http.BKN,
		T:          p.c.Http.T,
	}
	itemsResp, err := items.Get()
	if err != nil {
		return errors.Wrap(err, "items.Get")
	}
	//todo 整理目录，用户可自动选择下载目标
	//目前将目录所有文件遍历
	for _, term := range itemsResp.Result.Terms {
		for _, chapter := range term.ChapterInfo {
			for _, sub := range chapter.SubInfo {
				for _, info := range sub.TaskInfo {
					//获取文件的token
					token := &Token{
						TermID: p.TermID,
						FileID: info.Video.Vid,
						BKN:    p.c.Http.BKN,
						T:      p.c.Http.T,
						Cookie: p.c.Http.Cookie,
					}
					tokenRet, err := token.Get()
					if err != nil {
						log.Printf("get token err:%s, fileID:%s", err, info.Video.Vid)
						continue
					}
					//获取视频信息
					media := MediaInfo{
						Sign:  tokenRet.Sign,
						T:     tokenRet.T,
						Exper: 0,
						Us:    tokenRet.Us,
						Vid:   info.Video.Vid,
					}
					vodUrl, err := media.Get()
					if err != nil {
						log.Printf("get vodUrl err:%s, fileID:%s", err, info.Video.Vid)
						continue
					}
					//下载视频，由于m3u8转mp4主要消耗的是cpu/gpu资源，于是此处不考虑开启并发
					err = ffmpeg.New(&p.c.Ffmpeg).Do(vodUrl, info.Video.Name)
					if err != nil {
						log.Printf("save mp4 err:%s, fileID:%s", err, info.Video.Vid)
						continue
					}
				}

			}
		}

	}
	return nil
}
