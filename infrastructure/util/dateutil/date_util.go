package dateutil

import "time"

// 获取当前时间 - 毫秒
func GetNowByMs() int64 {
	return time.Now().UnixNano() / 1000_000
}

// 毫秒转纳秒
func ConvertMs2Ns(ms int64) int64 {
	return ms * 1000_000
}

// 纳秒转毫秒
func ConvertNs2Ms(ns int64) int64 {
	return ns / 1000_000
}

// 由于原生的time.RF3339不支持2021-03-30T21:11:39+0800这种格式，so需要如下转换
func ConvertISO8601DateStr(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05-0700", dateStr)
}
