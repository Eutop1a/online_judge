package leaderboard

import (
	"github.com/gin-gonic/gin"
	v1 "online_judge/api/v1"
)

type ApiLeaderboard struct{}

func (l *ApiLeaderboard) InitLeaderboard(RouterGroup *gin.RouterGroup) {
	leaderboardApi := v1.ApiGroupApp.ApiLeaderboard

	RouterGroup.GET("/user", leaderboardApi.GetUserLeaderboard)
}
