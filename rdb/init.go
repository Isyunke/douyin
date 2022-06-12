package rdb

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/ser163/WordBot/generate"
	"github.com/warthecatalyst/douyin/config"
	"github.com/warthecatalyst/douyin/logx"
)

var rdb *redis.Client

func Init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RdbHost, config.RdbPort),
		Password: "",
		DB:       0,
		PoolSize: 100,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		logx.DyLogger.Panicf("connect redis error, err=%+v", err)
	}

	setSalts()
	return
}

func createRandomString(count int) []string {
	var randStrs []string
	for i := 1; i <= count; i++ {
		wordList, _ := generate.GenRandomMix(10)
		randStrs = append(randStrs, wordList.Word)
	}
	return randStrs
}

func setSalts() {
	salts := GetAllSalts()
	if len(salts) != 0 {
		logx.DyLogger.Infof("salts already exist!")
		return
	}
	err := rdb.SAdd(keySalt, createRandomString(10)).Err()
	if err != nil {
		logx.DyLogger.Panicf("set salts error, err=%+v", err)
	}
	return
}

func GetRdb() *redis.Client {
	return rdb
}
