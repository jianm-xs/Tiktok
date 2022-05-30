package controller

import (
	"Project/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CommentAction 评论操作
func CommentAction(c *gin.Context) {
	c.JSON(http.StatusOK, models.CommentActionResponse{
		Response: models.Response{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
	})
}

// CommentList 评论列表
func CommentList(c *gin.Context) {
	c.JSON(http.StatusOK, models.CommentListResponse{
		Response: models.Response{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
		CommentList: nil,
	})
}
