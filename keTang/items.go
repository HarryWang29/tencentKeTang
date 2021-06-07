package keTang

import (
	"fmt"
	"github.com/iris-contrib/schema"
	"github.com/pkg/errors"
	"net/url"
)

type Items struct {
	CID        string  `url:"cid"`
	TermIDList string  `url:"term_id_list"`
	BKN        int64   `url:"bkn"`
	T          float32 `url:"t"`
}

type ItemsResp struct {
	Result struct {
		ServerTime int `json:"server_time"`
		Terms      []struct {
			RoomID         int `json:"room_id"`
			Sponsor        int `json:"sponsor"`
			RecentLiveTask struct {
			} `json:"recent_live_task"`
			TermRefundTime int `json:"term_refund_time"`
			SignEndtime    int `json:"sign_endtime"`
			Cycle          int `json:"cycle"`
			Bgtime         int `json:"bgtime"`
			RoomScale      int `json:"room_scale"`
			TermNo         int `json:"term_no"`
			ChapterInfo    []struct {
				ChID      int    `json:"ch_id"`
				Introduce string `json:"introduce"`
				ChNo      int    `json:"ch_no"`
				Name      string `json:"name"`
				SubInfo   []struct {
					Csid      int    `json:"csid"`
					SubID     int    `json:"sub_id"`
					Introduce string `json:"introduce"`
					Name      string `json:"name"`
					Endtime   int    `json:"endtime"`
					TermID    int    `json:"term_id"`
					TaskInfo  []struct {
						RestrictFlag int    `json:"restrict_flag"`
						CreateTime   int    `json:"create_time"`
						Csid         int    `json:"csid"`
						Introduce    string `json:"introduce"`
						Timelen      int    `json:"timelen"`
						SpecialFlag  int    `json:"special_flag"`
						Endtime      int    `json:"endtime"`
						ResidExt     struct {
							Times   int    `json:"times"`
							Txcloud int    `json:"txcloud"`
							Vid     string `json:"vid"`
						} `json:"resid_ext"`
						TermID int `json:"term_id"`
						Video  struct {
							Vid      string `json:"vid"`
							SetVid   int    `json:"set_vid"`
							TimeLen  int    `json:"time_len"`
							CoverURL string `json:"cover_url"`
							Name     string `json:"name"`
							State    int    `json:"state"`
							Aid      int    `json:"aid"`
						} `json:"video"`
						Type        int         `json:"type"`
						Bgtime      int         `json:"bgtime"`
						ExprFlag    int         `json:"expr_flag"`
						TeList      []int64     `json:"te_list"`
						Name        string      `json:"name"`
						TaskBitFlag int         `json:"task_bit_flag"`
						ResidList   interface{} `json:"resid_list"`
						ExprRange   int         `json:"expr_range"`
						AppendFlag  int         `json:"append_flag"`
						Aid         int         `json:"aid"`
						Taid        string      `json:"taid"`
						Cid         int         `json:"cid"`
					} `json:"task_info"`
					Bgtime int `json:"bgtime"`
					Cid    int `json:"cid"`
				} `json:"sub_info"`
				TermID int `json:"term_id"`
				Type   int `json:"type"`
				Aid    int `json:"aid"`
				Cid    int `json:"cid"`
			} `json:"chapter_info"`
			Price         int     `json:"price"`
			TranscodeFlag int     `json:"transcode_flag"`
			TermWarnType  int     `json:"term_warn_type"`
			TermBitFlag   int     `json:"term_bit_flag"`
			SignBgtime    int     `json:"sign_bgtime"`
			Introduce     string  `json:"introduce"`
			PubTime       int     `json:"pub_time"`
			ExprNum       int     `json:"expr_num"`
			Endtime       int     `json:"endtime"`
			TermID        int     `json:"term_id"`
			LiveRoom      int     `json:"live_room"`
			SignMax       int     `json:"sign_max"`
			TeList        []int64 `json:"te_list"`
			RoomURL4C     string  `json:"room_url4c"`
			Name          string  `json:"name"`
			ApplyNum      int     `json:"apply_num"`
			RoomURL       string  `json:"room_url"`
			LiveVid       int     `json:"live_vid"`
			Aid           int     `json:"aid"`
			Cid           int     `json:"cid"`
		} `json:"terms"`
	} `json:"result"`
	Retcode int `json:"retcode"`
}

func (a *api) Get(i *Items) (resp *ItemsResp, err error) {
	if i == nil {
		return nil, errors.New("param is nil")
	}
	i.BKN = a.c.BKN
	i.T = a.c.T
	v := url.Values{}
	err = schema.NewEncoder().Encode(i, v)
	if err != nil {
		return nil, errors.Wrap(err, "schema.NewEncoder().Encode")
	}
	err = a.get(fmt.Sprintf("%s%s", ItemsUri, v.Encode()), &resp,
		"referer", "https://ke.qq.com/webcourse/index.html")
	if err != nil {
		return nil, errors.Wrap(err, "a.get")
	}

	return resp, nil
}
