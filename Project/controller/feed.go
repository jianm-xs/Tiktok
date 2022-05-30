// feed 包，该包封装了视频流的接口
// 创建人：龚江炜
// 创建时间：2022-5-14

package controller

import (
	"Project/dao"
	"Project/models"
	"Project/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
)

// Feed : 视频流接口，用于请求视频列表
// 参数 :
//      c : 返回的信息（状态和视频列表）

func Feed(c *gin.Context) {
	var result models.FeedResponse // 响应结果
	var lastTime string
	// 限制返回视频的最新投稿时间，可能为空。默认为 -1
	// 客户端有 bug，所以暂时不使用请求时附带的时间戳
	//lastTimeInt, err := strconv.ParseInt(c.DefaultQuery("latest_time", string(-1)), 10, 64)
	//if err == nil && lastTimeInt != -1 { // 如果有返回时间戳，转为 string 类型
	//	tm := time.Unix(lastTimeInt, 0)             // 转为 Time 类型
	//	lastTime = tm.Format("2006-01-02 15:04:05") // 格式化为字符串
	//}
	token := c.DefaultQuery("token", "") // 用户的鉴权 token，可能为空
	var userId int64
	myClaims, err := utils.ParseToken(token)
	if err != nil { // token 解析失败
		userId = -1 // 当做没有登录处理
	} else { // 如果 token 解析成功，获取 userId
		userId, _ = strconv.ParseInt(myClaims.Uid, 10, 64)
	}
	result.VideoList = dao.GetVideos(lastTime, userId) // 执行数据库查询，获取结果
	result.Response.StatusCode = 0                     // 成功，设置状态码和描述
	result.Response.StatusMsg = "success"
	if len(result.VideoList) != 0 { // 如果有返回视频，更新 nextTime。方便下次获取视频列表时使用
		length := len(result.VideoList)                                // 获取视频数
		result.NextTime = result.VideoList[length-1].CreateTime.Unix() // 获取最后一个视频的时间，转为 int
	}
	c.JSON(http.StatusOK, result) // 设置返回的信息
}
