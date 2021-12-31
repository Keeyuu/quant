package request

type CalcOrderStatsReq struct {
	FromTs    int64    `json:"from_ts"`
	ToTs      int64    `json:"to_ts"`
	OrderIds  []string `json:"order_ids"`
	UserIds   []string `json:"user_ids"`
	DeviceIds []string `json:"device_ids"`
	StoreIds  []string `json:"store_ids"`
	AllowFail bool     `json:"allow_fail"` // 出现失败后是否继续执行剩下的处理逻辑
}

type ReportOrderStatsReq struct {
	Date       string `form:"date"`        // 统计日期："2006-01-02"
	DateOffset int    `form:"date_offset"` // 日期偏移量，如：-1, 0, 2。可与date结合使用
}
