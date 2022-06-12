# douyin 抖音项目服务端
## 一、项目环境配置与运行

- MYSQL8.0
- Redis 3.2
- ffmpeg>=4.1.1

> 有需要本人阿里云oss账户进行测试的，可以邮件联系本人18317011442@163.com

### 项目依赖安装
```shell
go mod tidy
```
会下载如下的依赖：
```shell
go download
    github.com/aliyun/aliyun-oss-go-sdk v2.2.4+incompatible
    github.com/dgrijalva/jwt-go v3.2.0+incompatible
    github.com/gin-gonic/gin v1.7.7
    github.com/go-redis/redis v6.15.9+incompatible
    github.com/google/uuid v1.3.0
    github.com/ser163/WordBot v1.0.0
    github.com/u2takey/ffmpeg-go v0.4.1
    gopkg.in/ini.v1 v1.66.5
    gorm.io/driver/mysql v1.3.4
    gorm.io/gorm v1.23.5
```

### 项目运行
```shell
go run main.go
```

## 项目结构
### 项目整体结构设计图


### 项目

接口功能不完善，仅作为示例

* 用户登录数据保存在内存中，单次运行过程中有效
* 视频上传后会保存到本地 public 目录中，访问时用 127.0.0.1:8080/static/video_name 即可

### 测试数据

测试数据写在 demo_data.go 中，用于列表接口的 mock 测试