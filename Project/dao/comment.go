package dao

import (
	"Project/models"
	"Project/utils"
	"errors"
	"gorm.io/gorm"
	"time"
)

const (
	PUBLISH = 1
	DELETE  = 2
)

// PerformCommentAction 登录用户对视频评论进行操作
// param:
//		userID: 请求的用户 id
// 		videoID: 评论所属的视频 id
//		actionType: 操作类型：1-发布评论，2-删除评论
// 		commentText: 用户填写的评论内容，在 action_type=1 的时候使用
//		commentID: 要删除的评论id，在 action_type=2 的时候使用
// 返回值：
//		操作成功：返回评论相关信息，nil
//		否则返回 nil, error
func PerformCommentAction(userID int64, videoID int64, actionType int,
	commentText string, commentID int64) (*models.Comment, error) {
	// 有这个视频吗？
	if err := DB.Debug().Where("video_id = ?", videoID).Error; err != nil {
		return nil, errors.New("video does not exist")
	}
	switch actionType {
	case PUBLISH:
		return CreateComment(userID, videoID, commentText)
	case DELETE:
		return DeleteComment(userID, commentID)
	default:
		// 防御性
		return nil, errors.New("invalid operation")
	}
}

// DeleteComment 删除评论
// param:
//		userID: 请求的用户 id
//		commentID: 要删除的评论id
// 返回值：
//		操作成功：返回评论相关信息，nil
//		否则返回 nil, error
func DeleteComment(userID int64, commentID int64) (*models.Comment, error) {
	var comment models.Comment
	err := DB.Debug().Where("id = ?", commentID).Find(&comment).Error
	if err != nil {
		return nil,
			errors.New("comment does not exist or has been deleted")
	}
	// 作者是否是请求发起人？
	if comment.AuthorID != userID {
		return nil,
			errors.New("you do not have permission to perform this action")
	}
	// commentID 字段不合法
	if err := DB.Debug().Where("comment_id = ?", commentID).Error; err != nil {
		return nil,
			errors.New("invalid comment_id")
	}
	if comment.IsDelete {
		return nil, errors.New("comment has been deleted")
	}
	DB.Model(&comment).Update("is_delete", 1)
	return &comment, nil
}

// CreateComment 删除评论
// param:
//		userID: 请求的用户 id
// 		videoID: 评论所属的视频 id
//		ctx(context): 评论内容
// 返回值：
//		操作成功：返回评论相关信息，nil
//		否则返回 nil, error
func CreateComment(userID int64, videoID int64,
	ctx string) (*models.Comment, error) {

	// 不能有空评论（大概吧）
	if ctx == "" {
		return nil, errors.New("comment must not be null or empty string")
	}
	id, err := utils.CommentIdWorker.NextId()
	if err != nil {
		// 生成 ID 异常
		return nil, err
	}
	// 忽略 `is_follow` 字段，默认值 false
	// 评论是用户自己发布的，自己不能关注自己
	var (
		comment models.Comment
		author  *models.User
	)
	err = DB.Debug().Omit("is_follow").
		Where("user_id = ?", userID).Find(&author).Error
	if err != nil {
		return nil, err
	}
	comment = models.Comment{
		ID:         id,
		Author:     *author,
		AuthorID:   userID,
		VideoID:    videoID,
		Content:    ctx,
		IsDelete:   false,
		CreateTime: utils.JsonTime(time.Now()),
		UpdateTime: time.Now(),
	}
	result := DB.Debug().Create(&comment)
	if result.Error != nil {
		// 插入异常
		return nil, result.Error
	}
	return &comment, nil
}

// GetCommentList 返回评论列表
// param:
//		userID: 请求的用户 id
// 		videoID: 评论所属的视频 id
// 返回值：
//		操作成功：返回评论列表
// 		操作失败：返回报错信息
func GetCommentList(userId, videoId int64) ([]models.Comment, error) {
	var comments []models.Comment
	// 查询 follow
	queryFollow := DB.Raw("? UNION ALL ?",
		DB.Select("? as user_id, 1 as is_follow", userId).Table("follow"),                            // 自己不能关注自己
		DB.Select("follow.user_id, 1 as is_follow").Where("follower_id = ?", userId).Table("follow"), // 查找当前用户关注的所有用户
	)
	err := DB.Debug().Table("comment").
		// 预加载 User，给 user 表加上 is_follow 字段再查找
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Select("user.*, is_follow").
				Joins("LEFT JOIN (?) AS fo ON user.user_id = fo.user_id", queryFollow).
				Table("user")
		}).
		// 选择返回的字段
		Select("comment.*").
		// 按照创建时间降序排列，即时间最晚的在前面
		Order("comment.create_time DESC").
		// 筛选条件，video_id
		Where("comment.video_id = ?", videoId).
		Find(&comments).Error
	if err != nil { // 查询出错
		return nil, err
	}
	return comments, nil // 查询成功
}
