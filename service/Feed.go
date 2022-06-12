package service

import (
	"time"

	"github.com/warthecatalyst/douyin/api"
	"github.com/warthecatalyst/douyin/dao"
	"github.com/warthecatalyst/douyin/logx"
)

func Feed(latestTime time.Time) (api.Feed, error) {
	videos, err := dao.NewVideoDaoInstance().GetLatest(latestTime)
	if err != nil {
		logx.DyLogger.Errorf("dao.NewVideoDaoInstance().GetLatest error: %s", err)
		return api.Feed{
			Response: api.Response{
				StatusCode: api.InnerErr,
				StatusMsg:  api.ErrorCodeToMsg[api.InnerErr],
			},
		}, err
	}
	if len(videos) == 0 {
		logx.DyLogger.Debug("当前无视频！")
		return api.Feed{Response: api.OK}, nil
	}

	v := newVideoList(videos)
	for i := 0; i < len(v); i++ {
		//查询视频作者信息
		resp, err := UserInfo(videos[i].UserID)
		if err != nil {
			return api.Feed{
				Response: api.Response{
					StatusCode: api.InnerErr,
					StatusMsg:  api.ErrorCodeToMsg[api.InnerErr],
				},
			}, err
		}
		v[i].Author = resp.User //作者信息
	}

	return api.Feed{VideoLists: v, Response: api.OK, NextTime: videos[len(videos)-1].CreateAt.UnixMilli()}, nil
}
