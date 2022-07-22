package project

import (
	"crawler/tencentKeTang/util"
	"fmt"
	"time"
)

type Catalogue struct {
	Name  string
	ID    []int64
	Depth int64
	Data  interface{}
}

func (a *api) GetCatalogue(cid string, tid int64) (list []*Catalogue, err error) {
	resp, err := a.keTang.BasicInfo(cid)
	if err != nil {
		return nil, err
	}
	a.catalogueName = util.ReplaceName(resp.Result.CourseDetail.Name)
	list = make([]*Catalogue, 0)
	//todo 更改为树结构
	for _, term := range resp.Result.CourseDetail.Terms {
		if tid != 0 && term.TermID != tid {
			continue
		}
		list = append(list, &Catalogue{
			Name:  fmt.Sprintf("%s[term_id:%d]", term.Name, term.TermID),
			ID:    []int64{term.TermID},
			Depth: 0,
			Data:  term,
		})
		for _, chapter := range term.ChapterInfo {
			list = append(list, &Catalogue{
				Name:  fmt.Sprintf("%s[ch_id:%d]", chapter.Name, chapter.ChID),
				ID:    []int64{chapter.ChID},
				Depth: 1,
				Data:  chapter,
			})
			for _, sub := range chapter.SubInfo {
				catalogue := &Catalogue{
					Depth: 2,
					Data:  sub,
				}
				bgTime := time.Unix(sub.Bgtime, 0)
				endTime := time.Unix(sub.Endtime, 0)
				catalogue.Name = fmt.Sprintf("%s[sub_id:%d]", sub.Name, sub.SubID)
				if endTime.After(time.Now()) {
					catalogue.Name += fmt.Sprintf("(%s~%s)",
						bgTime.Format("01月02日 15:04"),
						endTime.Format("15:04"),
					)
				}
				list = append(list, catalogue)
				for _, task := range sub.TaskInfo {
					catalogue := &Catalogue{
						Depth: 3,
						Data:  task,
					}
					catalogue.ID = a.string2list(task.ResidList)
					if len(catalogue.ID) == 0 {
						bgTime := time.Unix(task.Bgtime, 0)
						endTime := time.Unix(task.Endtime, 0)
						catalogue.Name = fmt.Sprintf("%s(%s~%s)",
							task.Name,
							bgTime.Format("01月02日 15:04"),
							endTime.Format("15:04"),
						)
					} else {
						catalogue.Name = fmt.Sprintf("%s[file_id:%v]", task.Name, catalogue.ID)
					}
					list = append(list, catalogue)
				}
			}
		}
	}
	a.catalogues = list
	return list, nil
}
