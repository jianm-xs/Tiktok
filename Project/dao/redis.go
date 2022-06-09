// redis 定义了 Redis 的配置和连接
// 创建人：龚江炜
// 创建时间：2022-6-6

package dao

import (
	"Project/config"
	"Project/models"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"strconv"
)

var RedisDB *redis.Client

// Redis 相关常量定义

const (
	HasKey    = true
	HasNotKey = false
)

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
func FindRedis(key string, val *int64, id string) error {
	// 查看 Redis 中是否存储了该值
	hasKey, err := RedisDB.HExists(key, id).Result()
	if err != nil {
		return err
	}
	switch hasKey {
	case HasKey:
		// 如果存在对应数据，将该数据赋值给 val
		count, _ := RedisDB.HGet(key, id).Result()
		*val, err = strconv.ParseInt(count, 10, 64)
	case HasNotKey:
		// 如果不存在对应数据，将该数据存入 Redis
		err = RedisDB.HSet(key, id, *val).Err()
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
func IncreaseValue[T models.User | models.Video](key string, data T, column string, table, id string) error {
	// 查看 Redis 中是否存在该键值
	hasKey, err := RedisDB.HExists(key, id).Result()
	if err != nil { // 查询出错
		return err
	}
	switch hasKey {
	case HasKey:
		// 如果 Redis 中有该数据，直接操作 Redis 将该数据 + 1
		err = RedisDB.HIncrBy(key, id, 1).Err()
	case HasNotKey:
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
		err = RedisDB.HSet(key, id, count+1).Err()
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
func DecreaseValue[T models.User | models.Video](key string, data T, column string, table string, id string) error {
	// 查看 Redis 中是否存在该键值
	hasKey, err := RedisDB.HExists(key, id).Result()
	if err != nil { // 查询出错
		return err
	}
	switch hasKey {
	case HasKey:
		// 如果 Redis 中有该数据，直接操作 Redis 将该数据 - 1
		err = RedisDB.HIncrBy(key, id, -1).Err()
	case HasNotKey:
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
		err = RedisDB.HSet(key, id, count-1).Err()
	default:
		// 防御性
		return errors.New("invalid operation")
	}
	return nil
}

// CreateData 在 Redis 中创建数据
// 参数 :
//	key: 键
//  data: 数据结构体
//  id: ID 值
// 返回值：
//		错误信息
func CreateData[T models.Favorite | models.Follow](key string, data T, id string) error {
	// 数据解析为字符串
	val, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// 读取 Redis 中是否有该数据，如果有，可以获取值，如果没有，只能获取 nil
	redisData, _ := RedisDB.HGet(key, id).Result()
	if redisData == "" || redisData == "nil" {
		// 说明不存在该数据，创建数据
		err = RedisDB.HSet(key, id, val).Err()
	} else {
		// 如果有该数据，无法继续创建
		err = errors.New("don't repeat the operation")
	}
	return err
}

// DeleteData 在 Redis 中将对应的值加一，如果 Redis 中没有该数据，从数据库中查询后修改
// 参数 :
//	key: 键
//  id: ID 值
// 返回值：
//		错误信息
func DeleteData(key string, id string) error {
	// 查看 Redis 中是否存在该键值
	hasKey, err := RedisDB.HExists(key, id).Result()
	if err != nil { // 查询出错
		return err
	}
	switch hasKey {
	case HasKey:
		// 如果 Redis 中有该数据，获取值
		data, _ := RedisDB.HGet(key, id).Result()
		// 如果值为 nil， 说明不存在该值
		if data == "nil" {
			err = errors.New("no such data")
		} else {
			// 如果有该数据，删除：值改为 nil
			err = RedisDB.HSet(key, id, "nil").Err()
		}
	case HasNotKey:
		// 如果没有该数据，为了防止缓存击穿，增加该键值到 Redis
		err = RedisDB.HSet(key, id, "nil").Err()
	}
	return err
}

// UpdateVideos 更新视频信息，如果更新出现问题，返回错误信息
// 参数 :
//	videos: 视频数组
// 	id: 用户 id ，用于查看是否点赞该视频
// 返回值：
//		错误信息
func UpdateVideos(videos []models.Video, id string) error {
	for i, _ := range videos {
		videoId := strconv.FormatInt(videos[i].ID, 10)
		favoriteKey := "video:favoriteCount"
		commentKey := "video:commentCount"
		// 更新点赞数
		err := FindRedis(favoriteKey, &videos[i].FavoriteCount, videoId)
		if err != nil {
			return err
		}
		// 更新评论数
		err = FindRedis(commentKey, &videos[i].CommentCount, videoId)
		if err != nil {
			return err
		}

		// 更新用户信息
		err = UpdateUser(&videos[i].Author, id)
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
func UpdateUsers(users []models.User, id string) error {
	for i, _ := range users {
		userId := strconv.FormatInt(users[i].ID, 10)
		followKey := "user:followCount"
		followerKey := "user:followerCount"
		// 更新关注数
		err := FindRedis(followKey, &users[i].FollowCount, userId)
		if err != nil {
			return err
		}
		// 查看是否点赞了

		// 更新粉丝数
		err = FindRedis(followerKey, &users[i].FollowerCount, userId)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateUser 更新用户信息，如果更新出现问题，返回错误信息
// 参数 :
//	user: 单个用户，传引用
// 	id: 用户 id ，用于查看是否点赞该视频
// 返回值：
//		错误信息
func UpdateUser(user *models.User, id string) error {
	userId := strconv.FormatInt(user.ID, 10)
	followKey := "user:followCount"
	followerKey := "user:followerCount"
	// 更新关注数
	err := FindRedis(followKey, &user.FollowCount, userId)
	if err != nil {
		return err
	}
	// 更新粉丝数
	err = FindRedis(followerKey, &user.FollowerCount, userId)
	if err != nil {
		return err
	}
	// 查看是否点赞了
	isFollowKey := userId + ":" + id
	// 读取 Redis 中的该数据
	isFollow, _ := RedisDB.HGet("favorite", isFollowKey).Result()
	if userId != id && (isFollow == "nil" || isFollow == "") {
		// 没有关注
		user.IsFollow = false
	} else {
		user.IsFollow = true
	}
	return nil
}
