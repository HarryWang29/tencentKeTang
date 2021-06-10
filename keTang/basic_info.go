package keTang

import (
	"fmt"
	"github.com/iris-contrib/schema"
	"github.com/pkg/errors"
	"net/url"
)

type basicInfo struct {
	CID string  `url:"cid"`
	BKN int64   `url:"bkn"`
	T   float32 `url:"t"`
}

type BasicTask struct {
	RestrictFlag int64   `json:"restrict_flag"`
	CreateTime   int64   `json:"create_time"`
	Csid         int64   `json:"csid"`
	Introduce    string  `json:"introduce"`
	SpecialFlag  int64   `json:"special_flag"`
	Endtime      int64   `json:"endtime"`
	ResidExt     string  `json:"resid_ext"`
	TermID       int64   `json:"term_id"`
	Type         int64   `json:"type"`
	Bgtime       int64   `json:"bgtime"`
	ExprFlag     int64   `json:"expr_flag"`
	TeList       []int64 `json:"te_list"`
	Name         string  `json:"name"`
	TaskBitFlag  int64   `json:"task_bit_flag"`
	ResidList    string  `json:"resid_list"`
	ExprRange    int64   `json:"expr_range"`
	AppendFlag   int64   `json:"append_flag"`
	Aid          int64   `json:"aid"`
	Taid         string  `json:"taid"`
	Cid          int64   `json:"cid"`
}

type BasicSub struct {
	Csid      int64        `json:"csid"`
	SubID     int64        `json:"sub_id"`
	Introduce string       `json:"introduce"`
	Name      string       `json:"name"`
	Endtime   int64        `json:"endtime"`
	TermID    int64        `json:"term_id"`
	TaskInfo  []*BasicTask `json:"task_info"`
	Bgtime    int64        `json:"bgtime"`
	Cid       int64        `json:"cid"`
}

type BasicChapter struct {
	ChID      int64       `json:"ch_id"`
	Introduce string      `json:"introduce"`
	ChNo      int64       `json:"ch_no"`
	Name      string      `json:"name"`
	SubInfo   []*BasicSub `json:"sub_info"`
	TermID    int64       `json:"term_id"`
	Type      int64       `json:"type"`
	Aid       int64       `json:"aid"`
	Cid       int64       `json:"cid"`
}

type BasicTerm struct {
	RoomID         int64           `json:"room_id"`
	Sponsor        int64           `json:"sponsor"`
	RecentLiveTask struct{}        `json:"recent_live_task"`
	TermRefundTime int64           `json:"term_refund_time"`
	SignEndtime    int64           `json:"sign_endtime"`
	Cycle          int64           `json:"cycle"`
	Bgtime         int64           `json:"bgtime"`
	RoomScale      int64           `json:"room_scale"`
	TermNo         int64           `json:"term_no"`
	ChapterInfo    []*BasicChapter `json:"chapter_info"`
	Price          int64           `json:"price"`
	IsCanDegrade   bool            `json:"is_can_degrade"`
	TranscodeFlag  int64           `json:"transcode_flag"`
	TermWarnType   int64           `json:"term_warn_type"`
	TermBitFlag    int64           `json:"term_bit_flag"`
	SignBgtime     int64           `json:"sign_bgtime"`
	Introduce      string          `json:"introduce"`
	PubTime        int64           `json:"pub_time"`
	ExprNum        int64           `json:"expr_num"`
	Endtime        int64           `json:"endtime"`
	TermID         int64           `json:"term_id"`
	LiveRoom       int64           `json:"live_room"`
	SignMax        int64           `json:"sign_max"`
	TeList         []int64         `json:"te_list"`
	RoomURL4C      string          `json:"room_url4c"`
	Name           string          `json:"name"`
	ApplyNum       int64           `json:"apply_num"`
	RoomURL        string          `json:"room_url"`
	LiveVid        int64           `json:"live_vid"`
	Aid            int64           `json:"aid"`
	Cid            int64           `json:"cid"`
}

type BasicInfoResp struct {
	Result struct {
		CourseDetail struct {
			AgencyCoverURL string       `json:"agency_cover_url"`
			Recordtime     int64        `json:"recordtime"`
			CoverURLColor  string       `json:"cover_url_color"`
			CategoryLabel  []int64      `json:"category_label"`
			LevelLabel     int64        `json:"level_label"`
			IosPrice       int64        `json:"ios_price"`
			PasscardNum    int64        `json:"passcard_num"`
			Bgtime         int64        `json:"bgtime"`
			AgencyDomain   string       `json:"agency_domain"`
			Score          int64        `json:"score"`
			Terms          []*BasicTerm `json:"terms"`
			Price          int64        `json:"price"`
			PasscardPrice  int64        `json:"passcard_price"`
			ServiceQq      []struct {
				Nick string `json:"nick"`
				URL  string `json:"url"`
				Uin  int64  `json:"uin"`
			} `json:"service_qq"`
			CourseBitFlag int64  `json:"course_bit_flag"`
			Details       string `json:"details"`
			Payid         int64  `json:"payid"`
			SalerUin      int64  `json:"saler_uin"`
			Vip           int64  `json:"vip"`
			AgencyType    int64  `json:"agency_type"`
			Train         int64  `json:"train"`
			Summary       string `json:"summary"`
			Industry1     int64  `json:"industry1"`
			Passcard      int64  `json:"passcard"`
			CoverURL      string `json:"cover_url"`
			Goal          string `json:"goal"`
			RefundType    int64  `json:"refund_type"`
			AgencyName    string `json:"agency_name"`
			Industry3     int64  `json:"industry3"`
			Industry2     int64  `json:"industry2"`
			AgencyID      int64  `json:"agency_id"`
			Endtime       int64  `json:"endtime"`
			AgencySummay  string `json:"agency_summay"`
			Pinyin        string `json:"pinyin"`
			Showid        int64  `json:"showid"`
			BlackFlag     int64  `json:"black_flag"`
			StoreNum      int64  `json:"store_num"`
			Name          string `json:"name"`
			ApplyNum      int64  `json:"apply_num"`
			Aid           int64  `json:"aid"`
			Cid           int64  `json:"cid"`
			Status        int64  `json:"status"`
			IsAgencySaler int64  `json:"is_agency_saler"`
		} `json:"course_detail"`
	} `json:"result"`
	Retcode int64 `json:"retcode"`
}

func (a *api) BasicInfo(cid string) (resp *BasicInfoResp, err error) {
	i := basicInfo{
		CID: cid,
		BKN: a.c.BKN,
		T:   a.c.T,
	}
	v := url.Values{}
	err = schema.NewEncoder().Encode(i, v)
	if err != nil {
		return nil, errors.Wrap(err, "schema.NewEncoder().Encode")
	}
	_, err = a.get(fmt.Sprintf("%s%s", BasicInfoUri, v.Encode()), &resp,
		"referer", "https://ke.qq.com/webcourse/index.html",
		"cookie", a.c.Cookie)
	if err != nil {
		return nil, errors.Wrap(err, "a.get")
	}

	return resp, nil

}
