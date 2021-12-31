package stringutil

import (
	"fmt"
	"regexp"
	"strings"
)

// 字符串的占位符格式是 #{参数名}，该方法会获取参数名进行值替换
// eg，a=1&b=2&paymentId=#{paymentId}, 则会解析出paymentId并进行替换
func ReplacePlaceholder(str string, valueMap map[string]interface{}) (newStr string, err error) {
	newStr = str
	urlRegexp, err := regexp.Compile(`#{[0-9A-Za-z_]*}`)
	if err != nil {
		return
	}
	allPlaceholder := urlRegexp.FindAllString(str, -1)
	if allPlaceholder == nil || len(allPlaceholder) == 0 {
		return
	}
	for _, itm := range allPlaceholder {
		paramKey := itm[2 : len(itm)-1]
		value := ""
		if v, ok := valueMap[paramKey]; ok {
			if v != nil {
				value = fmt.Sprintf("%v", v)
			} else {
				value = "null"
			}
		}
		str = strings.Replace(str, itm, value, -1)
	}
	newStr = str
	return
}
