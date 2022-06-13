package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/warthecatalyst/douyin/api"
	"github.com/warthecatalyst/douyin/logx"
	"github.com/warthecatalyst/douyin/rdb"
	"github.com/warthecatalyst/douyin/service"
	"github.com/warthecatalyst/douyin/tokenx"
)

type UIDSrc int

const (
	FromCtx UIDSrc = iota
	FromQuery
)

func getUserId(c *gin.Context, src UIDSrc) (int64, error) {
	var userId int64
	if src == FromCtx {
		userIdInterface, _ := c.Get("user_id")
		var ok bool
		userId, ok = userIdInterface.(int64)
		if !ok {
			logx.DyLogger.Errorf("context中user_id(%v)非int！", userIdInterface)
			c.JSON(http.StatusOK, api.Response{
				StatusCode: api.InputFormatCheckErr,
				StatusMsg:  "参数错误"})
			return tokenx.InvalidUserId, errors.New("参数错误")
		}
	} else if src == FromQuery {
		userIdStr := c.Query("user_id")
		var err error
		userId, err = strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			logx.DyLogger.Errorf("strconv.Atoi error: %s", err)
			c.JSON(http.StatusOK, api.Response{
				StatusCode: api.InputFormatCheckErr,
				StatusMsg:  "参数错误"})
			return tokenx.InvalidUserId, errors.New("参数错误")
		}
	}

	return userId, nil
}

type UserLoginResponse struct {
	api.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	api.Response
	User api.User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	if username == "" || password == "" {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: api.Response{
				StatusCode: api.InputFormatCheckErr,
				StatusMsg:  "用户名和密码不能为空，请重新输入",
			},
		})
		return
	}

	userId, token, err := service.NewUserServiceInstance().CreateUser(username, password)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: api.Response{StatusCode: api.LogicErr, StatusMsg: fmt.Sprintf("注册失败：%s", err.Error())},
		})
		return
	}
	logx.DyLogger.Debugf("[Register] userId=%+v, token=%+v", userId, token)
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: api.OK,
		UserId:   userId,
		Token:    token,
	})

	return
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	if username == "" || password == "" {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: api.Response{StatusCode: api.InputFormatCheckErr, StatusMsg: "用户名和密码不能为空"},
		})
		return
	}

	user, err := service.NewUserServiceInstance().LoginCheck(username, password)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: api.Response{StatusCode: api.LogicErr, StatusMsg: "内部错误"},
		})
		return
	}
	if user == nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: api.Response{StatusCode: api.LogicErr, StatusMsg: "用户名/密码错误"},
		})
		return
	}
	token := rdb.GetToken(user.Id)
	if token == "" {
		token = tokenx.CreateToken(user.Id, username)
		if err := rdb.AddToken(user.Id, token); err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: api.Response{StatusCode: api.InnerErr, StatusMsg: "签发token失败"},
			})
			return
		}
	}
	logx.DyLogger.Debugf("[Login] userId=%+v, token=%+v", user.Id, token)
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: api.OK,
		UserId:   user.Id,
		Token:    token,
	})

	return
}

func UserInfo(c *gin.Context) {
	userId, err := getUserId(c, FromQuery)
	if err != nil {
		logx.DyLogger.Errorf("Can't get userId from query")
		c.JSON(http.StatusOK, api.Response{StatusCode: api.InputFormatCheckErr, StatusMsg: api.ErrorCodeToMsg[api.InputFormatCheckErr]})
		return
	}
	loginUserId, err := getUserId(c, FromCtx)
	if err != nil {
		logx.DyLogger.Errorf("Can't get userId from context")
		c.JSON(http.StatusOK, api.Response{StatusCode: api.InputFormatCheckErr, StatusMsg: api.ErrorCodeToMsg[api.InputFormatCheckErr]})
		return
	}

	user, err := service.NewUserServiceInstance().GetUserByUserId(loginUserId, userId)
	if user == nil {
		c.JSON(http.StatusOK, api.Response{
			StatusCode: api.LogicErr,
			StatusMsg:  "用户不存在",
		})
		return
	}
	logx.DyLogger.Infof("[UserInfo] user=%+v", *user)
	c.JSON(http.StatusOK, UserResponse{
		Response: api.Response{StatusCode: 0},
		User:     *user,
	})

	return
}
