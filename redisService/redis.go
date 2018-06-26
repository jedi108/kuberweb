package redisService

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
)

type RedisCache struct {
	client map[string]*redis.Client
}

func NewRedisClients() *RedisCache {
	return &RedisCache{
		client: make(map[string]*redis.Client),
	}
}

func (rc *RedisCache) AddRedisAndInit(addr string, dbNum int) {
	index := fmt.Sprintf("%v,db:%v", addr, dbNum)
	rc.client[index] = redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   dbNum,
	})

	pong, err := rc.client[index].Ping().Result()
	if err != nil {
		logger.Fatalf("redis init failse - instanse addr:%v, pong %v", err, pong)
	}
}

type RedisInfos struct {
	RedisInfos []*redisInfo
}

type redisInfo struct {
	Err   string
	Res   string
	Names string
}

func (rc *RedisCache) FlushAll() *RedisInfos {

	rIs := &RedisInfos{}
	var errString = ""

	for k, v := range rc.client {
		result, err := v.FlushAll().Result()
		if err != nil {
			errString = fmt.Sprintf("%v", err)
		} else {
			errString = ""
		}

		rIs.RedisInfos = append(rIs.RedisInfos, &redisInfo{
			Err:   errString,
			Res:   fmt.Sprintf("%v", result),
			Names: k,
		})

	}
	return rIs
}

func (rc *RedisCache) GetRedisInfo() *RedisInfos {
	rIs := &RedisInfos{}
	var errString = ""
	for k, v := range rc.client {
		result, err := v.DBSize().Result()
		if err != nil {
			errString = fmt.Sprintf("%v", err)
		} else {
			errString = ""
		}
		rIs.RedisInfos = append(rIs.RedisInfos, &redisInfo{
			Err:   errString,
			Res:   fmt.Sprintf("%v", result),
			Names: k,
		})
	}
	return rIs
}
