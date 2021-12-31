package stringutil

import (
	"strconv"
	"strings"
)

func Slice2Map(strSlice []string) (strMap map[string]string) {
	strMap = make(map[string]string, len(strSlice))
	for _, str := range strSlice {
		strMap[str] = str
	}
	return
}

// index 下标从0开始
func SplitByAndGet(str string, separator string, index int, isDef bool) string {
	if index < 0 {
		return ""
	}
	strSlice := strings.Split(str, separator)
	if len(strSlice) > index {
		return strSlice[index]
	} else if len(strSlice) > 0 && len(strSlice) <= index && isDef {
		return strSlice[0]
	}
	return ""
}

func FormDecInt2Str(n int64) string {
	return strconv.FormatInt(n, 10)
}

func StrJoin(strs ...string) string {
	var builder strings.Builder
	for i := 0; i < len(strs); i++ {
		builder.WriteString(strs[i])
	}
	return builder.String()
}
