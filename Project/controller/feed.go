// feed 包，该包封装了视频流的接口
// 创建人：龚江炜
// 创建时间：2022-5-14

package controller

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

// FeedResponse 视频流接口响应对象，定义了视频流响应的基本内容
type FeedResponse struct {
	Response          // 状态码、状态描述
	NextTime  int64   `json:"next_time,omitempty"` // 本次返回的视频中，发布最早的时间。可以为空
	VideoList []Video `json:"video_list"`          // 本次返回的视频列表。可以为空
}

// Feed : 视频流接口，用于请求视频列表
// 参数 :
//      c : 返回的信息（状态和视频列表）

func Feed(c *gin.Context) {
	var result FeedResponse                                         // 结果
	db, _ := sql.Open("mysql", "root:root@(localhost:3306)/Tiktok") // 设置参数
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			result.Response.StatusCode = -1 // 更改状态码
			result.Response.StatusMsg = "Close database error!"
			result.NextTime = 0
			result.VideoList = nil
			c.JSON(http.StatusOK, result) // 设置返回的信息
			return
		}
	}(db) // 使用完毕后关闭数据库
	err := db.Ping() // 连接数据库
	if err != nil {  // 连接失败处理
		result.Response.StatusCode = -2
		result.Response.StatusMsg = "Connect database error!"
		result.NextTime = 0
		result.VideoList = nil
		c.JSON(http.StatusOK, result)
		return
	}

	// 以下为数据库连接测试代码，实际功能待实现
	// 预计完善时间：数据库创建完成后完善
	queryCommand := "SELECT video_id, video.user_id, name, play_url, cover_url FROM video, `user` WHERE video.user_id = `user`.user_id;" // 查询语句
	answer, _ := db.Query(queryCommand)                                                                                                  // 执行查询语句

	for answer.Next() {
		var video Video
		err := answer.Scan(&video.Id, &video.Author.Id, &video.Author.Name, &video.PlayUrl, &video.CoverUrl) // 获取查询结果
		if err != nil {                                                                                      // 读取失败处理
			result.Response.StatusCode = -3
			result.Response.StatusMsg = "Read video or user error!"
			result.NextTime = 0
			c.JSON(http.StatusOK, result)
			return
		}
		result.VideoList = append(result.VideoList, video) // 将该视频放入结果中
	}

	result.Response.StatusCode = 0 // 成功，设置状态码和描述
	result.Response.StatusMsg = "success"
	result.NextTime = 1
	c.JSON(http.StatusOK, result) // 设置返回的信息
}
