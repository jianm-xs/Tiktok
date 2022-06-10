// favorite 包，该包包含了点赞的数据库操作
// 创建人：龚江炜
// 创建时间：2022-5-15

package dao

import (
	"Project/models"
	"Project/utils"
	"errors"
	"log"
	"strconv"
	"time"
)

// FavoriteAction 点赞操作
// 参数 :
//		userId : 请求的用户 id
// 		videoId : 点赞的视频 id
//		actionType : 操作类型： 1 -> 点赞， 2 -> 取消点赞
// 返回值：
//		如果操作成功，返回 nil， 否则返回错误信息
func FavoriteAction(userId, videoId, actionType int64) error {
	var err error
	// 通过 VideoId 获取对应的 Key 值，方便 Redis 操作
	videoIdStr := strconv.FormatInt(videoId, 10)
	key := "video:favoriteCount"
	favoriteID := strconv.FormatInt(userId, 10) + ":" + strconv.FormatInt(videoId, 10)
	switch actionType {
	case PUBLISH:
		// 点赞操作，插入这条记录
		favorite := models.Favorite{
			UserID:     userId,
			VideoID:    videoId,
			CreateTime: time.Now(),
		}
		// 雪花算法生成 ID
		if favorite.ID, err = utils.FavoriteIDWorker.NextID(); err != nil {
			log.Println(err)
			return errors.New("create ID error") //	ID 生成失败返回-1
		}
		// 插入数据到 Redis 中
		err = CreateData("favorite", favorite, favoriteID)
		if err != nil { // 数据写入失败
			return err
		}
		// 该视频的点赞数 + 1
		err = IncreaseValue(key, models.Video{ID: videoId}, "favorite_count", "video", videoIdStr)
		return err
	case DELETE:
		// 把这条数据从 Redis 中删除
		err = DeleteData("favorite", favoriteID)
		if err != nil { // 删除点赞记录失败
			return err
		}
		// 该视频的点赞数 - 1
		err = DecreaseValue(key, models.Video{ID: videoId}, "favorite_count", "video", videoIdStr)
		return err
	default:
		// 防御性
		return errors.New("invalid operation")
	}
}

// GetFavoriteList 点赞列表
// 参数 :
//		userId : 请求的用户 id
// 返回值：
//		[]models.Video 成功返回点赞列表，失败返回nil
//		error 成功，返回 nil， 否则返回错误信息
func GetFavoriteList(authorId, userId int64) ([]models.Video, error) {
	var videos []models.Video
	// 更新点赞数据到 MySQL 中
	err := updateFavoriteData()
	if err != nil {
		// 如果更新失败，返回错误
		return nil, err
	}

	// 查询点赞
	queryFavorite := DB.Select("video_id,create_time,1 as is_favorite").
		Where("favorite_id = ?", authorId).
		Table("favorite")

	err = DB.Debug().Table("video").
		// 预加载 User
		Preload("Author").
		// 联结点赞视频
		Joins("JOIN (?) AS fa ON fa.video_id = video.video_id", queryFavorite).
		// 按照点赞时间降序排列，即时间最晚的在前面
		Order("fa.create_time DESC").
		Find(&videos).Error
	if err != nil {
		// 如果查询失败，返回错误信息
		return nil, err
	}
	// 更新视频的点赞数和评论数
	userIdStr := strconv.FormatInt(userId, 10)
	err = UpdateVideos(videos[:], userIdStr)
	if err != nil {
		// 如果更新出现问题，返回错误
		return nil, err
	}
	return videos, nil
}
