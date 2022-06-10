package main

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/warthecatalyst/douyin/api"
	"github.com/warthecatalyst/douyin/common"
	"github.com/warthecatalyst/douyin/controller"
	"github.com/warthecatalyst/douyin/dao"
	"github.com/warthecatalyst/douyin/service"
	"github.com/warthecatalyst/douyin/tokenx"
	"net/http"
	"os"
	"strings"
)

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		userId, username := tokenx.ParseToken(token)
		if username == "" {
			// TODO: 端上应该重定向到登录界面
			c.AbortWithStatusJSON(http.StatusOK, api.Response{StatusCode: api.LogicErr, StatusMsg: "非法token"})
			return
		}
		if user, err := service.NewUserServiceInstance().GetUserByUserId(userId); err != nil {
			c.AbortWithStatusJSON(http.StatusOK, api.Response{StatusCode: api.LogicErr, StatusMsg: "内部错误"})
			return
		} else if user == nil {
			c.AbortWithStatusJSON(http.StatusOK, api.Response{StatusCode: api.LogicErr, StatusMsg: "非法用户"})
			return
		}
		c.Set("user_id", userId)
		c.Next()
	}
}

func initData() {
	path := "data/AccessKey.txt"
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	for i := 0; i < 2; i++ {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\r\n")
		if common.UploadAccessKeyID == "" {
			common.UploadAccessKeyID = input
		} else {
			common.UploadAccessKeySecret = input
		}
	}
	fmt.Println(common.UploadAccessKeyID)
	fmt.Println(common.UploadAccessKeySecret)
}

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/user/", controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", controller.Publish)
	apiRouter.GET("/publish/list/", controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
	apiRouter.POST("/comment/action/", controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)
}

func initAll() {
	initData()
	dao.InitDB()
	//rdb.InitRdb()

}

func main() {
	initAll()
	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
