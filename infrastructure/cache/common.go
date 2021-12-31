package cache

import (
	"app/infrastructure/config"
	"crypto/md5"
	"fmt"
	"time"
)

func ExpireInterval() time.Duration {
	return time.Duration(config.Get().Cache.ExpireInterval) * time.Minute
}

func CreateKeyWithMd5(group string, params ...interface{}) string {
	data := []byte(fmt.Sprintf("%v", params))
	return fmt.Sprintf("%s:%s:%x", config.Get().Name, group, md5.Sum(data))
}

func CreateKey(group string, params ...interface{}) string {
	key := fmt.Sprintf("%s:%s", config.Get().Name, group)
	if params == nil || len(params) <= 0 {
		return key
	}
	if len(params) >= 1 {
		key += ":"
	}
	for i := 0; i < len(params); i++ {
		key += fmt.Sprintf("%v", params[i])
		if i+1 < len(params) {
			key += ":"
		}
	}
	return key
}
