package dao

import (
	"errors"
	"github.com/warthecatalyst/douyin/model"
	"gorm.io/gorm"
	"sync"
)

type FavouriteDao struct{}

var (
	favoriteDao  *FavouriteDao
	favoriteOnce sync.Once
)

func NewFavoriteDaoInstance() *FavouriteDao {
	favoriteOnce.Do(
		func() {
			favoriteDao = &FavouriteDao{}
		})
	return favoriteDao
}

//QueryCountOfVideo 方法 从Video表中查询点赞的数据
func (*FavouriteDao) QueryCountOfVideo(conditions map[string]interface{}) (int32, error) {
	var video model.Video
	err := db.Where(conditions).First(&video).Error
	if err != nil {
		return 0, err
	}
	return video.FavoriteCount, err
}

//IsFavourite 查询 userID的用户是否对videoID的视频进行点赞
func (*FavouriteDao) IsFavourite(userID, videoID int64) bool {
	var fav model.Favourite
	result := db.Where("user_id = ? AND video_id = ?", userID, videoID).First(&fav)
	return result.RowsAffected != 0
}

//Add 向数据库中增加一条点赞记录
func (*FavouriteDao) Add(userID, videoID int64) error {
	f := model.Favourite{
		UserID:  userID,
		VideoID: videoID,
	}
	//通过事务实现，由于事务具有ACID特性
	trans := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			trans.Rollback()
		}
	}()
	err := trans.Error
	if err != nil {
		return err
	}
	//点赞表中写入数据
	err = trans.Create(&f).Error
	if err != nil {
		trans.Rollback()
		return err
	}

	//在Video表中更新点赞记录
	err = trans.Model(&model.Video{}).Where("video_id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
	if err != nil {
		trans.Rollback()
		return err
	}
	return trans.Commit().Error
}

//Del 从数据库中删除一条点赞记录
func (*FavouriteDao) Del(userID, videoID int64) error {
	trans := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			trans.Rollback()
		}
	}()
	err := trans.Error
	if err != nil {
		return err
	}
	//点赞表中删除数据
	err = trans.Where("video_id = ? AND user_id = ?", videoID, userID).Delete(&model.Favourite{}).Error
	if err != nil {
		trans.Rollback()
		return err
	}

	//在Video表中更新点赞记录
	err = trans.Model(&model.Video{}).Where("video_id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error
	if err != nil {
		trans.Rollback()
		return err
	}
	return trans.Commit().Error
}

//VideoIDListByUserID 获取某用户点赞的所有视频的ID列表
func (*FavouriteDao) VideoIDListByUserID(userID int64) ([]int64, error) {
	var f []model.Favourite
	err := db.Model(&model.Favourite{}).
		Select("video_id").
		Where("user_id = ?", userID).
		Find(&f).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	var res []int64
	for _, i := range f {
		res = append(res, i.VideoID)
	}
	return res, nil
}
