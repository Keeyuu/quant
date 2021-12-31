package cache

import (
	"app/infrastructure/util/log"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/bluele/gcache"
)

var cache gcache.Cache
var mutex sync.Mutex

func Cache() gcache.Cache {
	return cache
}

func init() {
	mutex.Lock()
	defer mutex.Unlock()

	log.Info("start init gcache")
	cache = gcache.New(1000).LRU().AddedFunc(func(key, value interface{}) {
		//log.Debug(fmt.Sprintf("gcache added key:%s", key))
	}).EvictedFunc(func(key, value interface{}) {
		log.Debug(fmt.Sprintf("gcache evicted key:%s", key))
	}).PurgeVisitorFunc(func(key, value interface{}) {
		log.Debug(fmt.Sprintf("gcache purge key:%s", key))
	}).Build()
	log.Info("finish init gcache")
}

func PutStructValueJsonMemory(key string, value interface{}, expiration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return Cache().SetWithExpire(key, string(bytes), expiration)
}

func GetStructValueJsonMemory(key string, obj interface{}) (err error) {
	valueStr, err := Cache().Get(key)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(valueStr.(string)), obj)
	return err
}

func PutStructValueMemory(key string, value interface{}, expiration time.Duration) error {
	return Cache().SetWithExpire(key, value, expiration)
}

func GetStructValueMemory(key string) (interface{}, error) {
	return Cache().Get(key)
}

func PutStringValueMemory(key string, value string, expiration time.Duration) error {
	return Cache().SetWithExpire(key, value, expiration)
}

func GetStringValueMemory(key string) (value string, err error) {
	val, err := Cache().Get(key)
	if val == nil {
		return "", err
	}
	return val.(string), err
}
