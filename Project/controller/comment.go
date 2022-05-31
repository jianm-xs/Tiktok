package controller

import (
	"Project/dao"
	"Project/models"
	"Project/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CommentAction 评论操作
func CommentAction(c *gin.Context) {
	var q models.CommentActionRequest
	q.UserID = utils.String2int64(c.DefaultPostForm("user_id", "-1"))
	q.Token = c.DefaultPostForm("token", "")
	q.VideoID = utils.String2int64(c.DefaultPostForm("video_id", ""))
	q.ActionType = int(utils.String2int64(c.DefaultPostForm("action_type", "-1")))
	q.CommentText = c.DefaultPostForm("comment_text", "")
	q.CommentID = utils.String2int64(c.DefaultPostForm("comment_id", "-1"))

	// Token 匹配？
	uidStr := c.DefaultPostForm("user_id", "-1")
	if ok, err := utils.CheckToken(q.Token, uidStr); err != nil || ok != true {
		c.JSON(http.StatusBadRequest, models.Response{
			StatusCode: -3,
			StatusMsg:  "authentication failed",
		})
		return
	}

	var comment *models.Comment
	comment, err := dao.PerformCommentAction(
		q.UserID,
		q.VideoID,
		q.ActionType,
		q.CommentText,
		q.CommentID,
	)
	fmt.Println(comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommentActionResponse{
		Response: models.Response{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
		Comment: models.Comment{
			ID:         comment.ID,
			Author:     comment.Author,
			AuthorID:   comment.AuthorID,
			VideoID:    comment.VideoID,
			Content:    comment.Content,
			CreateTime: comment.CreateTime,
			UpdateTime: comment.UpdateTime,
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
