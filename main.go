package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/warthecatalyst/douyin/api"
	"github.com/warthecatalyst/douyin/controller"
	"github.com/warthecatalyst/douyin/dao"
	"github.com/warthecatalyst/douyin/oss"
	"github.com/warthecatalyst/douyin/rdb"
	"github.com/warthecatalyst/douyin/service"
	"github.com/warthecatalyst/douyin/tokenx"
)

func CheckLogin(mustLogin bool, getTokenFromUrl bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ""
		if getTokenFromUrl {
			token = c.Query("token")
		} else {
			token = c.PostForm("token")
		}
		if token == "" && !mustLogin {
			c.Set("user_id", tokenx.InvalidUserId)
			return
		}
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

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", CheckLogin(false, true), controller.Feed)
	apiRouter.GET("/user/", CheckLogin(true, true), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", CheckLogin(true, false), controller.Publish)
	apiRouter.GET("/publish/list/", CheckLogin(true, true), controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", CheckLogin(true, true), controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", CheckLogin(true, true), controller.FavoriteList)
	apiRouter.POST("/comment/action/", CheckLogin(true, true), controller.CommentAction)
	apiRouter.GET("/comment/list/", CheckLogin(true, true), controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", CheckLogin(true, true), controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", CheckLogin(true, true), controller.FollowList)
	apiRouter.GET("/relation/follower/list/", CheckLogin(true, true), controller.FollowerList)
}

func initAll() {
	dao.InitDB()
	rdb.Init()
	oss.Init()
}

func main() {
	initAll()
	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
