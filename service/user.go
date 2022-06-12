package service

import (
	"errors"

	"github.com/warthecatalyst/douyin/api"
	"github.com/warthecatalyst/douyin/dao"
	"github.com/warthecatalyst/douyin/idgenerator"
	"github.com/warthecatalyst/douyin/logx"
	"github.com/warthecatalyst/douyin/model"
	"github.com/warthecatalyst/douyin/rdb"
	"github.com/warthecatalyst/douyin/tokenx"
)

type UserService struct{}

var (
	userService = &UserService{}
)

func NewUserServiceInstance() *UserService {
	return userService
}

func (u *UserService) CreateUser(username string, password string) (int64, string, error) {
	userInfo, err := dao.NewUserDaoInstance().GetUserByUsername(username)
	if err != nil {
		logx.DyLogger.Errorf("GetUserByUsername error: %s", err)
		return -1, "", err
	}
	if userInfo != nil {
		return -1, "", errors.New("当前用户名已存在")
	}
	userId := idgenerator.GenerateUid()
	token := tokenx.CreateToken(userId, username)
	rdb.AddToken(userId, token)
	logx.DyLogger.Debugf("gen token=%v", token)

	user := &model.User{
		UserID:   userId,
		UserName: username,
		PassWord: password,
	}
	err = dao.NewUserDaoInstance().AddUser(user)
	if err != nil {
		logx.DyLogger.Errorf("AddUser error: %s", err)
		return -1, "", err
	}

	return userId, token, nil
}

func (u *UserService) GetUserByUserId(userId int64) (*api.User, error) {
	userModel, err := dao.NewUserDaoInstance().GetUserById(userId)
	if err != nil {
		return nil, err
	}
	if userModel == nil {
		return nil, nil
	}

	return &api.User{
		Id:            userId,
		Name:          userModel.UserName,
		FollowCount:   userModel.FollowCount,
		FollowerCount: userModel.FollowerCount,
	}, nil

}

func (u *UserService) LoginCheck(username, password string) (*api.User, error) {
	user, err := dao.NewUserDaoInstance().GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		logx.DyLogger.Errorf("没有该用户！（username = %s)", username)
		return nil, nil
	}

	if password != user.PassWord {
		logx.DyLogger.Errorf("密码不对！")
		return nil, nil
	}

	return &api.User{
		Id:            user.UserID,
		Name:          username,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
	}, nil
}

// UserInfo 获取用户信息
// 包括 user_id,name,follow_count,follower_count,is_favorite
// 最后一个字段应该和具体业务有关（我暂时还不太理解）
func UserInfo(id int64) (api.UserInfo, error) {
	// 获取到 user_id,follow_count,follower_count
	u, err := dao.NewUserDaoInstance().GetUserById(id)
	if err != nil {
		return api.UserInfo{
			Response: api.Response{StatusCode: api.InnerErr, StatusMsg: api.ErrorCodeToMsg[api.InnerErr]},
		}, err
	}

	return api.UserInfo{
		Response: api.Response{StatusCode: 0, StatusMsg: "success"},
		User:     model.UserQuery{ID: u.UserID, FollowCount: u.FollowCount, FollowerCount: u.FollowerCount, Name: u.UserName},
	}, nil
}
