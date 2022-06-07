// redis 定义了 Redis 的配置和连接
// 创建人：龚江炜
// 创建时间：2022-6-6

package dao

import (
	"Project/common"
	"Project/config"
	"Project/models"
	"errors"
	"github.com/go-redis/redis"
	"strconv"
)

var RedisDB *redis.Client

// InitRedis 初始化 Redis 客户端
// 返回值：
//		如果连接成功，返回 nil ，否则返回错误
func InitRedis() (err error) {
	// 创建一个 Redis 客户端
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     config.RedisCfg.Addr,
		Password: config.RedisCfg.Password,
		DB:       config.RedisCfg.DB,
	})
	// 检查是否成功连接到了 Redis 服务器
	_, err = RedisDB.Ping().Result()
	// 清理一下缓存。仅测试用，正式发布时要删除
	RedisDB.FlushAll()
	if err != nil {
		return err
	}
	return nil
}

// FindRedis 在 Redis 中寻找数据，如果没有该数据，则更新 Redis
// 参数 :
//	key: 键
//  val: 值，传指针，若 Redis 中有数据，则写入 val 中，否则将 val 写入 Redis 中
// 返回值：
//		错误信息
func FindRedis(key string, val *int64) error {
	// 查看 Redis 中是否存储了该视频的点赞数
	hasKey, err := RedisDB.Exists(key).Result()
	if err != nil {
		return err
	}
	switch hasKey {
	case common.HAS_KEY:
		// 如果存在对应数据，将该数据赋值给 val
		count, _ := RedisDB.Get(key).Result()
		*val, err = strconv.ParseInt(count, 10, 64)
	case common.HAS_NOT_KEY:
		// 如果不存在对应数据，将该数据存入 Redis
		err = RedisDB.Set(key, *val, 0).Err()
	default:
		return errors.New("invalid operation")
	}
	return err
}

// IncreaseValue 在 Redis 中将对应的值加一，如果 Redis 中没有该数据，从数据库中查询后修改
// 参数 :
//	key: 键
//  data: 结构体，只能为 User 或者 Video
// 	column: 字段名
//  table: 表名
// 返回值：
//		错误信息
func IncreaseValue[T models.User | models.Video](key string, data T, column string, table string) error {
	// 查看 Redis 中是否存在该键值
	hasKey, err := RedisDB.Exists(key).Result()
	if err != nil { // 查询出错
		return err
	}
	switch hasKey {
	case common.HAS_KEY:
		// 如果 Redis 中有该数据，直接操作 Redis 将该数据 + 1
		err = RedisDB.Incr(key).Err()
	case common.HAS_NOT_KEY:
		// 如果 Redis 中没有该数据，查询数据库，然后把结果写入 Redis

		// 用于获取数据库中的数据
		var count int64
		// 查询数据库中对应字段的值
		err = DB.Debug().Table(table).
			Select(column).
			Model(&data).
			Scan(&count).Error
		if err != nil { // 数据库查询失败
			return err
		}
		//	加一，再写入 Redis 中
		err = RedisDB.Set(key, count+1, 0).Err()
	default:
		// 防御性
		return errors.New("invalid operation")
	}
	return nil
}

// DecreaseValue 在 Redis 中将对应的值减一，如果 Redis 中没有该数据，从数据库中查询后修改
// 参数 :
//	key: 键
//  data: 结构体，只能为 User 或者 Video
// 	column: 字段名
//  table: 表名
// 返回值：
//		错误信息
func DecreaseValue[T models.User | models.Video](key string, data T, column string, table string) error {
	// 查看 Redis 中是否存在该键值
	hasKey, err := RedisDB.Exists(key).Result()
	if err != nil { // 查询出错
		return err
	}
	switch hasKey {
	case common.HAS_KEY:
		// 如果 Redis 中有该数据，直接操作 Redis 将该数据 - 1
		err = RedisDB.Decr(key).Err()
	case common.HAS_NOT_KEY:
		// 如果 Redis 中没有该数据，查询数据库，然后把结果写入 Redis

		// 用于获取数据库中的数据
		var count int64
		// 查询数据库中对应字段的值
		err = DB.Debug().Table(table).
			Select(column).
			Model(&data).
			Scan(&count).Error
		if err != nil { // 数据库查询失败
			return err
		}
		//	减一，再写入 Redis 中
		err = RedisDB.Set(key, count-1, 0).Err()
	default:
		// 防御性
		return errors.New("invalid operation")
	}
	return nil
}

// UpdateVideos 更新视频信息，如果更新出现问题，返回错误信息
// 参数 :
//	videos: 视频数组
// 返回值：
//		错误信息
func UpdateVideos(videos []models.Video) error {
	for i, _ := range videos {
		favoriteKey := "video_favoriteCount_" + strconv.FormatInt(videos[i].ID, 10)
		commentKey := "video_commentCount_" + strconv.FormatInt(videos[i].ID, 10)
		// 更新点赞数
		err := FindRedis(favoriteKey, &videos[i].FavoriteCount)
		if err != nil {
			return err
		}
		// 更新评论数
		err = FindRedis(commentKey, &videos[i].CommentCount)
		if err != nil {
			return err
		}
		// 更新用户信息
		err = UpdateUser(&videos[i].Author)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateUsers 更新用户信息，如果更新出现问题，返回错误信息
// 参数 :
//	users: 用户数组
// 返回值：
//		错误信息
func UpdateUsers(users []models.User) error {
	for i, _ := range users {
		followKey := "user_followCount_" + strconv.FormatInt(users[i].ID, 10)
		followerKey := "user_followerCount_" + strconv.FormatInt(users[i].ID, 10)
		// 更新关注数
		err := FindRedis(followKey, &users[i].FollowCount)
		if err != nil {
			return err
		}
		// 更新粉丝数
		err = FindRedis(followerKey, &users[i].FollowerCount)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateUser 更新用户信息，如果更新出现问题，返回错误信息
// 参数 :
//	user: 单个用户，传引用
// 返回值：
//		错误信息
func UpdateUser(user *models.User) error {
	followKey := "user_followCount_" + strconv.FormatInt(user.ID, 10)
	followerKey := "user_followerCount_" + strconv.FormatInt(user.ID, 10)
	// 更新关注数
	err := FindRedis(followKey, &user.FollowCount)
	if err != nil {
		return err
	}
	// 更新粉丝数
	err = FindRedis(followerKey, &user.FollowerCount)
	if err != nil {
		return err
	}
	return nil
}
