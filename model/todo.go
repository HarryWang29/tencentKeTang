package model

type TaskInfo struct {
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
	Type        int     `json:"type"`
	Bgtime      int     `json:"bgtime"`
	ExprFlag    int     `json:"expr_flag"`
	TeList      []int64 `json:"te_list"`
	Name        string  `json:"name"`
	TaskBitFlag int     `json:"task_bit_flag"`
	ResidList   string  `json:"resid_list"`
	ExprRange   int     `json:"expr_range"`
	AppendFlag  int     `json:"append_flag"`
	Aid         int     `json:"aid"`
	Taid        string  `json:"taid"`
	Cid         int     `json:"cid"`
}

type Todo struct {
	Csid      int         `json:"csid"`
	SubID     int         `json:"sub_id"`
	Introduce string      `json:"introduce"`
	Name      string      `json:"name"`
	Endtime   int         `json:"endtime"`
	TermID    int         `json:"term_id"`
	TaskInfo  []*TaskInfo `json:"task_info"`
	Bgtime    int         `json:"bgtime"`
	Cid       int         `json:"cid"`
}
