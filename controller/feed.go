package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/warthecatalyst/douyin/api"
	"github.com/warthecatalyst/douyin/service"
)

// TODO 需要考虑用户登录时和未登录时处理方式的区别
// Feed 推送视频流
func Feed(c *gin.Context) {
	latestTimeStr := c.Query("latest_time")
	latestTime, err := strconv.ParseInt(latestTimeStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, api.Response{
			StatusCode: api.InputFormatCheckErr,
			StatusMsg:  "latest_time不是整数时间戳格式！"})
		return
	}
	resp, err := service.Feed(time.Unix(0, latestTime * (int64(time.Millisecond))))
	if err != nil {
		c.JSON(http.StatusOK, resp.Response)
		return
	}

	c.JSON(http.StatusOK, resp)
}
