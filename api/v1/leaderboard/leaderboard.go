package leaderboard

import (
	"github.com/gin-gonic/gin"
	"online_judge/consts/resp_code"
	response2 "online_judge/models/common/response"
)

type ApiLeaderboard struct{}

func (l *ApiLeaderboard) GetLeaderboard(c *gin.Context) {

}

func (l *ApiLeaderboard) GetProblemLeaderboard(c *gin.Context) {

}

// GetUserLeaderboard 获取用户题解排名接口
// @Tags Rank API
// @Summary 获取用户题解排名
// @Description 获取用户题解排名接口
// @Accept multipart/form-data
// @Produce json
// @Success 200 {object} common.GetUserLeaderboardResponse "1000 获取用户题解排名成功"
// @Failure 200 {object} common.GetUserLeaderboardResponse "1022 获取用户题解排名失败"
// @Failure 200 {object} common.GetUserLeaderboardResponse "1014 服务器内部错误"
// @Router /leaderboard/user [GET]
func (l *ApiLeaderboard) GetUserLeaderboard(c *gin.Context) {
	response, err := LeaderboardService.GetUserLeaderboard()
	if err != nil {
		response2.ResponseError(c, response2.CodeGetUserRankError)
		return
	}
	switch response.Code {
	case resp_code.Success:
		response2.ResponseSuccess(c, response.Data)
	default:
		response2.ResponseError(c, response2.CodeInternalServerError)
	}
	return
}
