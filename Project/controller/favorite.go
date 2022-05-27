package controller

import (
	"Project/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// FavoriteAction 点赞接口
func FavoriteAction(c *gin.Context) {
	c.JSON(http.StatusOK, models.Response{
		StatusCode: 0,
		StatusMsg:  "success1",
	})
}

// FavoriteList 点赞列表接口
func FavoriteList(c *gin.Context) {
	c.JSON(http.StatusOK, models.VideoListResponse{
		Response: models.Response{
			StatusCode: 0,
			StatusMsg:  "success!"},
		VideoList: nil,
	})
}
