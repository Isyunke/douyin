package service

import (
	"time"

	"github.com/warthecatalyst/douyin/api"
	"github.com/warthecatalyst/douyin/dao"
	"github.com/warthecatalyst/douyin/logx"
)

func Feed(userId int64, latestTime time.Time) (int64, []api.Video, error) {
	videos, err := dao.NewVideoDaoInstance().GetLatest(latestTime)
	v := newVideoList(videos)
	if err != nil {
		logx.DyLogger.Errorf("dao.NewVideoDaoInstance().GetLatest error: %s", err)
		return -1, v, err
	}
	if len(v) == 0 {
		logx.DyLogger.Debug("当前无视频！")
		return -1, v, err
	}

	for i := 0; i < len(v); i++ {
		//查询视频作者信息
		userModel, err := dao.NewUserDaoInstance().GetUserById(v[i].Author.Id)
		if err != nil {
			logx.DyLogger.Errorf("GetUserByUserId error(uid = %d): %s", v[i].Author.Id, err)
			continue
		}
		if userModel == nil {
			logx.DyLogger.Errorf("user not found(uid = %d)", v[i].Author.Id)
			continue
		}
		v[i].Author = api.User{
			Id:            int64(userModel.UserID),
			Name:          userModel.UserName,
			FollowCount:   userModel.FollowCount,
			FollowerCount: userModel.FollowerCount,
		}
		follows, err := dao.NewFollowDaoInstance().FindFollow(userId, v[i].Author.Id)
		if err != nil {
			logx.DyLogger.Errorf("GetUserByUserId error(uid1 = %d, uid2 = %d): %s", userId, v[i].Author.Id, err)
			continue
		}
		v[i].Author.IsFollow = len(follows) > 0
	}

	return videos[len(videos)-1].CreateAt.UnixMilli(), v, nil
}
