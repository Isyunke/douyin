package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"

	"github.com/warthecatalyst/douyin/api"
	"github.com/warthecatalyst/douyin/config"
	"github.com/warthecatalyst/douyin/dao"
	"github.com/warthecatalyst/douyin/idgenerator"
	"github.com/warthecatalyst/douyin/logx"
	"github.com/warthecatalyst/douyin/model"
	"github.com/warthecatalyst/douyin/rdb"
	"github.com/warthecatalyst/douyin/tokenx"
)

func getMd5(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

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
		return tokenx.InvalidUserId, "", err
	}
	if userInfo != nil {
		return tokenx.InvalidUserId, "", errors.New("当前用户名已存在")
	}
	userId := idgenerator.GenerateUid()
	token := tokenx.CreateToken(userId, username)
	if err := rdb.AddToken(userId, token); err != nil {
		return tokenx.InvalidUserId, "", err
	}
	logx.DyLogger.Debugf("gen token=%v", token)

	user := &model.User{
		UserID:   userId,
		UserName: username,
	}
	if config.UserConf.PasswordEncrpted {
		user.PassWord = getMd5(password)
	} else {
		user.PassWord = password
	}
	err = dao.NewUserDaoInstance().AddUser(user)
	if err != nil {
		logx.DyLogger.Errorf("AddUser error: %s", err)
		return tokenx.InvalidUserId, "", err
	}

	return userId, token, nil
}

func (u *UserService) GetUserByUserId(loginUserId, userId int64) (*api.User, error) {
	userModel, err := dao.NewUserDaoInstance().GetUserById(userId)
	if err != nil {
		return nil, err
	}
	if userModel == nil {
		return nil, nil
	}

	user := &api.User{
		Id:            userId,
		Name:          userModel.UserName,
		FollowCount:   userModel.FollowCount,
		FollowerCount: userModel.FollowerCount,
	}

	if loginUserId == tokenx.InvalidUserId {
		user.IsFollow = false
	} else {
		follow, err := dao.NewFollowDaoInstance().FindFollow(loginUserId, userId)
		if err != nil {
			return nil, err
		}
		user.IsFollow = len(follow) > 0
	}

	return user, nil
}

func (u *UserService) UserExistByUserId(userId int64) (bool, error) {
	userModel, err := dao.NewUserDaoInstance().GetUserById(userId)
	if err != nil {
		return false, err
	}

	return userModel != nil, nil
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

	w := password
	if config.UserConf.PasswordEncrpted {
		w = getMd5(password)
	}
	if w != user.PassWord {
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
