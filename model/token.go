package model

type TokenReq struct {
	TermID string
	FileID string
	BKN    string
	T      string
	Cookie string
}

type TokenResp struct {
	Result struct {
		Sign  string `json:"sign"`
		T     string `json:"t"`
		Exper int    `json:"exper"`
		Us    string `json:"us"`
	} `json:"result"`
	Retcode int `json:"retcode"`
}
