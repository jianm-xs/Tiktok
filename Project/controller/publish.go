// publish 包，该包封装了投稿相关的借口
// 创建人：龚江炜
// 创建时间：2022-5-14

package controller

import (
	"Project/common"
	"Project/dao"
	"Project/models"
	"github.com/gin-gonic/gin"
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

// PublishList 登录用户的视频发布列表，直接列出用户所有投稿过的视频
func PublishList(context *gin.Context) {
	var request models.PublishListQuery
	if err := context.ShouldBindQuery(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_ = request.Token
	// TODO: JWT auth.

	videoList, err := dao.GetPublishListByUserId(request.UserID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, &models.Response{
		StatusCode: common.StatusOK,
		StatusMsg:  "success",
		Data: &models.PublishListResponse{
			VideoList: &videoList,
		},
	})
}
