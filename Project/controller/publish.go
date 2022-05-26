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
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Publish 投稿接口
func Publish(context *gin.Context) {
	var result models.Response

	// 获取视频数据
	data, err := context.FormFile("data")
	if err != nil { // 如果获取视频失败，返回信息
		result.StatusCode = -1
		result.StatusMsg = "Get video error!"
		context.JSON(http.StatusBadRequest, result)
		return
	}
	// TODO: 文件大小限制？

	// 获取作者 token 和视频标题
	token := context.PostForm("token")
	title := context.PostForm("title")

	var authorId int64
	myClaims, err := ParseToken(token)
	if err != nil { // token 解析失败
		result.StatusCode = -2              // 失败，设置状态码和描述
		result.StatusMsg = "token error!"   // token 有误
		context.JSON(http.StatusOK, result) // 设置返回的信息
		return
	} else { // 如果 token 解析成功，获取 userId
		authorId, _ = strconv.ParseInt(myClaims.Uid, 10, 64)
	}

	var sb strings.Builder

	// 获取当前时间戳
	fileNameInt := time.Now().Unix()
	// 文件名是时间戳 + 作者 id
	fileNameStr := strconv.FormatInt(fileNameInt, 10) + strconv.FormatInt(authorId, 10)

	// 带扩展名的文件名
	sb.WriteString(fileNameStr)
	// 获取文件类型，获取拓展名
	fileSuffix := path.Ext(data.Filename)
	// 文件名 = 文件名 + 类型名
	sb.WriteString(fileSuffix)
	// 转换为 string 作为视频文件名
	saveFileName := sb.String()

	filePath := filepath.Join("./video", "/", saveFileName)

	fmt.Println("=============>", filePath, authorId)
	// 暂时先保存到 server
	if err := context.SaveUploadedFile(data, filePath); err != nil {
		result.StatusCode = -3                 // 失败，设置状态码和描述
		result.StatusMsg = "save video error!" // token 有误
		context.JSON(http.StatusOK, result)    // 设置返回的信息
		return
	}
	// FIXME: coverUrl

	// 播放地址为服务器地址 + 文件路径
	playUrl := "http://81.70.17.190:1080/" + filePath
	coverUrl := "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg"
	fmt.Println("playUrl: ", playUrl)

	err = dao.CreateVideoByData(title, authorId, playUrl, coverUrl)
	if err != nil { // 如果发布失败，返回信息
		result.StatusCode = -4                        // 失败，设置状态码和描述
		result.StatusMsg = "Insert into Mysql error!" // token 有误
		context.JSON(http.StatusOK, result)           // 设置返回的信息
		return
	}

	// 发布成功
	context.JSON(http.StatusOK, &models.Response{
		StatusCode: common.StatusOK,
		StatusMsg:  "success",
	})
}

// VideoVo 接受从数据库查出来的部分数据的结构体
type VideoVo struct {
	Id            int64  `json:"id" db:"id"`                          // 视频 ID
	UserId        int64  `json:"user_id" db:"user_id"`                // userId
	PlayUrl       string `json:"play_url" db:"play_url"`              // 视频播放地址
	CoverUrl      string `json:"cover_url" db:"cover_url"`            // 视频封面地址
	FavoriteCount int64  `json:"favorite_count" db:"favourite_count"` // 视频点赞总数
	CommentCount  int64  `json:"comment_count" db:"comment_count"`    // 视频评论总数
	IsFavorite    bool   `json:"is_favorite" db:"isfavourite"`        // 是否已点赞
}

// PublishList 视频发布列表
func PublishList(c *gin.Context) {
	// 返回的结果
	var result models.VideoListResponse
	// 获取请求的 token
	token := c.DefaultQuery("token", "")
	// 获取作者的 id
	authorId, err := strconv.ParseInt(c.DefaultQuery("user_id", "-1"), 10, 64)
	if err != nil {
		result.StatusCode = -1
		result.StatusMsg = "get author id error!"
		c.JSON(http.StatusOK, result)
		return
	}
	// 获取当前用户的 id
	var userId int64
	myClaims, err := ParseToken(token)
	if err != nil { // token 解析失败
		result.StatusCode = -2            // 失败，设置状态码和描述
		result.StatusMsg = "token error!" // token 有误
		c.JSON(http.StatusOK, result)     // 设置返回的信息
		return
	} else { // 如果 token 解析成功，获取 userId
		userId, _ = strconv.ParseInt(myClaims.Uid, 10, 64)
	}

	//获取 author 发布视频信息，userId 用于判断是否关注了
	videos := dao.GetVideoList(authorId, userId)

	//接口返回
	c.JSON(http.StatusOK, models.VideoListResponse{
		Response: models.Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
