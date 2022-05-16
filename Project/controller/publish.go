// publish 包，该包封装了投稿相关的借口
// 创建人：龚江炜
// 创建时间：2022-5-14

package controller

import (
	"Project/common"
	"Project/dao"
	"Project/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
)

// Publish 投稿接口
func Publish(context *gin.Context) {
	var request models.PublishVideoRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_ = request.Token
	// TODO: JWT auth.

	// Data 需要: `PlayUrl`, `CoverUrl`，其余默认即可
	err := dao.CreateVideoByUserId(request.UserID, request.Data)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, &models.Response{
		StatusCode: common.StatusOK,
		StatusMsg:  "success",
	})
}

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}
type VideoVo struct {
	Id            int64  `json:"id" db:"id"`                          // 视频 ID
	UserId        int64  `json:"user_id" db:"user_id"`                // userId
	PlayUrl       string `json:"play_url" db:"play_url"`              // 视频播放地址
	CoverUrl      string `json:"cover_url" db:"cover_url"`            // 视频封面地址
	FavoriteCount int64  `json:"favorite_count" db:"favourite_count"` // 视频点赞总数
	CommentCount  int64  `json:"comment_count" db:"comment_count"`    // 视频评论总数
	IsFavorite    bool   `json:"is_favorite" db:"isfavourite"`        // 是否已点赞
}

var Db *sqlx.DB

func getUserVideoInfoByToken(token string) []Video {

	database, err := sqlx.Open("mysql", "root:root@tcp(110.42.225.78:3306)/test")
	if err != nil {
		fmt.Println("连接数据库失败：" + err.Error())
		return []Video{}
	}
	Db = database
	var videos []Video
	var videoVos []VideoVo
	var user []User

	err = Db.Select(&videoVos, "select * from video where user_id=?", 1)
	if err != nil {
		fmt.Println("exec failed, ", err)
		return []Video{}
	}
	err = Db.Select(&user, "select * from user where id=?", 1)
	if err != nil {
		fmt.Println("exec failed, ", err)
		return []Video{}
	}

	defer Db.Close()
	for _, v := range videoVos {
		t := Video{
			Id:            v.Id,
			Author:        user[0],
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
		}
		videos = append(videos, t)
	}

	for _, v := range videos {
		fmt.Println(v)
	}

	return videos
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	// 声明接收的变量
	var token string = "qqq"
	// 将request的body中的数据，自动按照json格式解析到结构体
	//if err := c.ShouldBindJSON(&token); err != nil {
	//	// 返回错误信息
	//	// gin.H封装了生成json数据的工具
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	videos := getUserVideoInfoByToken(token)

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
