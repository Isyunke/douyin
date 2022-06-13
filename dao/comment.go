package dao

import (
	"errors"
	"sync"

	"github.com/warthecatalyst/douyin/model"
	"gorm.io/gorm"
)

type CommentDao struct{}

var (
	commentDao  *CommentDao
	commentOnce sync.Once
)

func NewCommentDaoInstance() *CommentDao {
	commentOnce.Do(
		func() {
			commentDao = &CommentDao{}
		})
	return commentDao
}

//QueryCommentCountOfVideo 方法 从Video表中查询评论的数据
func (*CommentDao) QueryCommentCountOfVideo(conditions map[string]interface{}) (int32, error) {
	var video model.Video
	err := db.Where(conditions).First(&video).Error
	if err != nil {
		return 0, err
	}
	return video.CommentCount, err
}

//IsComment 查询 是否存在CommentId
func (*CommentDao) IsComment(commentID int64) bool {
	var com model.Comment
	result := db.Where("id = ?", commentID).First(&com)
	return result.RowsAffected != 0

}

//Add 向数据库中增加一条评论记录
func (*CommentDao) Add(userID, videoID int64, content string) error {
	c := model.Comment{
		UserID:  userID,
		VideoID: videoID,
		Content: content,
	}
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
	//评论表中写入数据
	err = trans.Create(&c).Error
	if err != nil {
		trans.Rollback()
		return err
	}

	//在Video表中更新评论记录
	err = trans.Model(&model.Video{}).Where("video_id = ?", videoID).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error
	if err != nil {
		trans.Rollback()
		return err
	}
	return trans.Commit().Error
}

//Del 从数据库中删除一条评论记录
func (*CommentDao) Del(commentId, videoID int64) error {
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
	//评论表中删除数据
	err = trans.Where("id = ?", commentId).Delete(&model.Comment{}).Error
	if err != nil {
		trans.Rollback()
		return err
	}

	//在Video表中更新点赞记录
	err = trans.Model(&model.Video{}).Where("video_id = ?", videoID).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error
	if err != nil {
		trans.Rollback()
		return err
	}
	return trans.Commit().Error

}

//CommentListByVideoID 获取视频的所有评论ID
func (*CommentDao) CommentListByVideoID(videoID int64) ([]int64, error) {
	var c []model.Comment
	err := db.Model(&model.Comment{}).Order("create_at desc").
		Select("id").
		Where("video_id = ?", videoID).
		Find(&c).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	var res []int64
	for _, i := range c {
		res = append(res, int64(i.ID))
	}
	return res, nil
}

func (*CommentDao) GetCommentFromId(commentId int64) (*model.Comment, error) {
	comment := &model.Comment{}
	if err := db.Where("id = ?", commentId).First(comment).Error; err != nil {
		return nil, err
	}

	return comment, nil
}

func (c *CommentDao) GetUserFromCommentId(commentId int64) (int64, error) {
	var comment model.Comment
	err := db.Select("user_id").Where("id = ?", commentId).First(&comment).Error
	if err != nil {
		return 0, err
	}
	return comment.UserID, nil

}
