package controller

import (
	"github.com/gin-gonic/gin"
	"online-judge/consts/resp_code"
	"online-judge/pkg/resp"
	"online-judge/services"
)

func GetLeaderboard(c *gin.Context) {

}

func GetProblemLeaderboard(c *gin.Context) {

}

// GetUserLeaderboard 获取用户题解排名接口
// @Tags Rank API
// @Summary 获取用户题解排名
// @Description 获取用户题解排名接口
// @Accept multipart/form-data
// @Produce json
// @Success 200 {object} models.GetUserLeaderboardResponse "1000 获取用户题解排名成功"
// @Failure 200 {object} models.GetUserLeaderboardResponse "1022 获取用户题解排名失败"
// @Failure 200 {object} models.GetUserLeaderboardResponse "1014 服务器内部错误"
// @Router /leaderboard/user [GET]
func GetUserLeaderboard(c *gin.Context) {
	var leaderboard services.Leaderboard
	response, err := leaderboard.GetUserLeaderboard()
	if err != nil {
		resp.ResponseError(c, resp.CodeGetUserRankError)
		return
	}
	switch response.Code {
	case resp_code.Success:
		resp.ResponseSuccess(c, response.Data)
	default:
		resp.ResponseError(c, resp.CodeInternalServerError)
	}
	return
}
