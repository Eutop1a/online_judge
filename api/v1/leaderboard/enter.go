package leaderboard

import "online_judge/services"

type ApiGroup struct {
	ApiLeaderboard
}

var (
	LeaderboardService = services.ServiceGroupApp.LeaderboardService
)
