package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/warthecatalyst/douyin/api"
	"github.com/warthecatalyst/douyin/logx"
	"github.com/warthecatalyst/douyin/service"
)

type VideoListResponse struct {
	api.Response
	VideoList []api.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	userId, err := getUserId(c) //得到UserId
	data, err := c.FormFile("data")
	if err != nil {
		logx.DyLogger.Error("Can't get file from form using key = 'data'!")
		c.JSON(http.StatusOK, api.Response{
			StatusCode: api.InnerErr,
			StatusMsg:  api.ErrorCodeToMsg[api.InnerErr],
		})
		return
	}
	title := c.PostForm("title") //视频名称
	filename := data.Filename
	logx.DyLogger.Info(title)
	err = service.PublishVideoInfo(data, userId, title)
	if err != nil {
		c.JSON(http.StatusOK, api.Response{
			StatusCode: api.InnerErr,
			StatusMsg:  api.ErrorCodeToMsg[api.InnerErr],
		})
	} else {
		c.JSON(http.StatusOK, api.Response{
			StatusCode: 0,
			StatusMsg:  "文件上传成功 : " + filename,
		})
	}
}

// PublishList 返回用户发布的所有视频列表
func PublishList(c *gin.Context) {
	userId, err := getUserId(c) //得到UserId
	videolist, err := service.PublishListInfo(userId)
	if err != nil {
		logx.DyLogger.Errorf("Can't get videoList from userId")
		c.JSON(http.StatusOK, api.Response{StatusCode: api.RecordNotExistErr, StatusMsg: api.ErrorCodeToMsg[api.RecordNotExistErr]})
		return
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: api.Response{
			StatusCode: 0,
		},
		VideoList: *videolist,
	})
}
