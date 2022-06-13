package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/warthecatalyst/douyin/api"
	"github.com/warthecatalyst/douyin/logx"
	"github.com/warthecatalyst/douyin/service"
)

type UserListResponse struct {
	api.Response
	UserList []*api.User `json:"user_list"`
}

func RelationAction(c *gin.Context) {
	userId, err := getUserId(c, FromCtx)
	if err != nil {
		logx.DyLogger.Errorf("Can't get userId from token")
		c.JSON(http.StatusOK, api.Response{StatusCode: api.TokenInvalidErr, StatusMsg: api.ErrorCodeToMsg[api.TokenInvalidErr]})
		return
	}
	actTyp := c.Query("action_type")
	actTypInt, err := strconv.Atoi(actTyp)
	if err != nil {
		c.JSON(http.StatusOK, api.Response{
			StatusCode: api.InputFormatCheckErr,
			StatusMsg:  fmt.Sprintf("strconv.Atoi error: %s", err)})
		return
	}
	toUserIdStr := c.Query("to_user_id")
	toUserId, err := strconv.Atoi(toUserIdStr)
	if err != nil {
		c.JSON(http.StatusOK, api.Response{
			StatusCode: api.InputFormatCheckErr,
			StatusMsg:  fmt.Sprintf("strconv.Atoi error: %s", err)})
		return
	}
	if err := service.FollowAction(userId, int64(toUserId), actTypInt); err != nil {
		c.JSON(http.StatusOK, api.Response{
			StatusCode: api.LogicErr,
			StatusMsg:  fmt.Sprintf("service.FollowAction error: %s", err)})
		return
	}
	c.JSON(http.StatusOK, api.Response{
		StatusCode: 0,
		StatusMsg:  ""})
	return
}

func FollowList(c *gin.Context) {
	userId, err := getUserId(c, FromQuery)
	if err != nil {
		logx.DyLogger.Errorf("Can't get userId from query")
		c.JSON(http.StatusOK, api.Response{StatusCode: api.InputFormatCheckErr, StatusMsg: api.ErrorCodeToMsg[api.InputFormatCheckErr]})
		return
	}
	users, err := service.GetFollowList(userId)
	if err != nil {
		c.JSON(http.StatusOK, api.Response{
			StatusCode: api.LogicErr,
			StatusMsg:  fmt.Sprintf("service.GetFollowList error: %s", err)})
		return
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: api.Response{
			StatusCode: 0,
		},
		UserList: users,
	})
}

func FollowerList(c *gin.Context) {
	userId, err := getUserId(c, FromQuery)
	if err != nil {
		logx.DyLogger.Errorf("Can't get userId from query")
		c.JSON(http.StatusOK, api.Response{StatusCode: api.InputFormatCheckErr, StatusMsg: api.ErrorCodeToMsg[api.InputFormatCheckErr]})
		return
	}
	users, err := service.GetFollowerList(userId)
	if err != nil {
		c.JSON(http.StatusOK, api.Response{
			StatusCode: api.LogicErr,
			StatusMsg:  fmt.Sprintf("service.GetFollowerList error: %s", err)})
		return
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: api.Response{
			StatusCode: 0,
		},
		UserList: users,
	})
}
