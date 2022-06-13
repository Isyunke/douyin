package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/warthecatalyst/douyin/api"
	"github.com/warthecatalyst/douyin/service"
)

type FeedResp struct {
	api.Response
	NextTime   int64       `json:"next_time"`
	VideoLists []api.Video `json:"video_list"`
}

// Feed 推送视频流
func Feed(c *gin.Context) {
	userId, err := getUserId(c, FromCtx)
	if err != nil {
		return
	}
	latestTimeStr := c.Query("latest_time")
	latestTime, err := strconv.ParseInt(latestTimeStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, FeedResp{
			Response: api.Response{
				StatusCode: api.InputFormatCheckErr,
				StatusMsg:  "latest_time不是整数时间戳格式！"}})
		return
	}
	nextTime, videoList, err := service.Feed(userId, time.Unix(0, latestTime*(int64(time.Millisecond))))
	if err != nil {
		c.JSON(http.StatusOK, FeedResp{
			Response: api.Response{
				StatusCode: api.InnerErr,
				StatusMsg:  api.ErrorCodeToMsg[api.InnerErr]}})
		return
	}

	c.JSON(http.StatusOK, FeedResp{
		Response:   api.OK,
		NextTime:   nextTime,
		VideoLists: videoList})
	return
}
