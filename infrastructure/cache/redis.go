package cache

import (
	"app/infrastructure/config"
	"app/infrastructure/util/log"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

const (
	GroupAccount = "account"
)

var once sync.Once

var redisInstance = struct {
	client *redis.Client
	err    error
	mutex  sync.Mutex
}{}

func init() {
	go once.Do(initRedis)
}

func initRedis() {
	if redisInstance.client == nil {
		redisInstance.mutex.Lock()
		defer redisInstance.mutex.Unlock()
		if redisInstance.client == nil {
			log.Info("start connect redis")
			redisOption := redis.Options{
				Addr:        config.Get().Redis.Addr,
				DialTimeout: config.Get().Redis.ConnectionTimeout,
			}
			if config.Get().Redis.Password != "" {
				redisOption.Password = config.Get().Redis.Password
			}
			client := redis.NewClient(&redisOption)
			redisInstance.client = client
			redisInstance.err = nil
			log.Info(fmt.Sprintf("finish connect redis:%v,err:%v", redisInstance.client, redisInstance.err))
			return
		}
	}
}

func redisClient() (*redis.Client, error) {
	initRedis()
	return redisInstance.client, redisInstance.err
}

func PutStructValue(key string, value interface{}, expiration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return PutValue(key, bytes, expiration)
}

func GetStructValue(key string, obj interface{}) (err error) {
	valueStr, err := GetValue(key)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(valueStr), obj)
	return err
}

func PutValue(key string, obj interface{}, expiration time.Duration) error {
	client, err := redisClient()
	if err != nil {
		return err
	}
	return client.Set(key, obj, expiration).Err()
}

func GetValue(key string) (value string, err error) {
	client, err := redisClient()
	if err != nil {
		return "", err
	}
	value, err = client.Get(key).Result()
	return
}

func DeleteKey(key string) (err error) {
	client, err := redisClient()
	if err != nil {
		return
	}
	_, err = client.Del(key).Result()
	return
}

func SetNXKey(key string, obj interface{}, expiration time.Duration) (bool, error) {
	client, err := redisClient()
	if err != nil {
		return false, err
	}
	return client.SetNX(key, obj, expiration).Result()
}
