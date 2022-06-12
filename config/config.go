package config

import (
	"strings"

	"github.com/warthecatalyst/douyin/logx"
	"gopkg.in/ini.v1"
	"gorm.io/gorm"
)

const (
	configFilePath = "./config/config.ini"
)

type OssConfig struct {
	Url             string
	Bucket          string
	BucketDirectory string
	AccessKeyID     string
	AccessKeySecret string
}

type VideoConfig struct {
	SavePath      string
	AllowedExts   []string
	UploadMaxSize int64
}

type UserConfig struct {
	PasswordEncrpted bool
}

// 解析配置文件
var (
	AppMode    string // 服务器启动模式默认 debug 模式
	Port       string //服务启动端口
	Dbtype     string //数据库类型
	DbHost     string //数据库服务器主机
	DbPort     string //数据服务器端口
	DbUser     string //数据库用户
	DbPassWord string //数据库密码
	DbName     string //数据库名
	DbLogLevel string //日志打印级别

	RdbHost string
	RdbPort string

	FeedListLength int

	OssConf OssConfig

	VideoConf VideoConfig

	UserConf UserConfig
)

func init() {
	f, err := ini.Load(configFilePath)
	if err != nil {
		logx.DyLogger.Panic("配置文件初始化失败")
	}

	loadServer(f)
	loadDb(f)
	loadRdb(f)
	loadFeed(f)
	loadOss(f)
	loadVideo(f)
	loadUser(f)
}

// loadServer 加载服务器配置
func loadServer(file *ini.File) {
	s := file.Section("server")
	AppMode = s.Key("AppMode").MustString("debug")
	Port = s.Key("Port").MustString("8080")
}

// loadDb 加载数据库相关配置
func loadDb(file *ini.File) {
	s := file.Section("database")
	Dbtype = s.Key("Dbtype").MustString("mysql")
	DbName = s.Key("DbName").MustString("douyin")
	DbPort = s.Key("DbPort").MustString("3306")
	DbHost = s.Key("DbHost").MustString("127.0.0.1")
	DbUser = s.Key("DbUser").MustString("")
	DbPassWord = s.Key("DbPassWord").MustString("")
	DbLogLevel = s.Key("LogLevel").MustString("error")
}

func loadRdb(file *ini.File) {
	s := file.Section("redis")
	RdbHost = s.Key("Host").MustString("127.0.0.1")
	RdbPort = s.Key("Port").MustString("6379")
}

func loadFeed(file *ini.File) {
	s := file.Section("feed")
	FeedListLength = s.Key("ListLength").MustInt(30)
}

func loadOss(file *ini.File) {
	s := file.Section("oss")
	OssConf.Url = s.Key("Url").MustString("")
	OssConf.Bucket = s.Key("Bucket").MustString("")
	OssConf.BucketDirectory = s.Key("BucketDirectory").MustString("")
	OssConf.AccessKeyID = s.Key("AccessKeyID").MustString("")
	OssConf.AccessKeySecret = s.Key("AccessKeySecret").MustString("")
}

func loadVideo(file *ini.File) {
	s := file.Section("video")
	VideoConf.SavePath = s.Key("SavePath").MustString("../userdata/")
	videoExts := s.Key("AllowedExts").MustString("mp4,wmv,avi")
	VideoConf.AllowedExts = strings.Split(videoExts, ",")
	VideoConf.UploadMaxSize = s.Key("UploadMaxSize").MustInt64(1024)
}

func loadUser(file *ini.File) {
	s := file.Section("user")
	UserConf.PasswordEncrpted=s.Key("PasswordEncrypted").MustBool(false)
}

var Db *gorm.DB

func GetDB() *gorm.DB {
	return Db
}
