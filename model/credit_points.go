package model

type CreditSummary struct {
	System string `json:"system"` // 个人预约制度
	Remain string `json:"remain"`
	Total  string `json:"total"`
}

type CreditRecord struct {
	Title    string `json:"title"`    // 原因标题
	Subtitle string `json:"subtitle"` // 扣分及时间
	Location string `json:"location"` // 地点及备注
}

type CreditPoints struct {
	Summary CreditSummary  `json:"summary"`
	Records []CreditRecord `json:"records"`
}
