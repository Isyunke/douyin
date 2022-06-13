package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/warthecatalyst/douyin/config"
	"github.com/warthecatalyst/douyin/logx"
)

var (
	bucket *oss.Bucket
)

func Init() {
	client, err := oss.New("https://"+config.OssConf.Url, config.OssConf.AccessKeyID, config.OssConf.AccessKeySecret)
	if err != nil {
		logx.DyLogger.Panicf("OSS初始化client失败：%s", err)
	}
	bucket, err = client.Bucket(config.OssConf.Bucket)
	if err!=nil{
		logx.DyLogger.Panicf("OSS初始化bucket失败：%s", err)
	}
}
