package leaderboard

import (
	"go.uber.org/zap"
	"online_judge/consts/resp_code"
	"online_judge/dao/mysql"
	"online_judge/models/common/response"
)

type LeaderboardService struct{}

func (l *LeaderboardService) GetUserLeaderboard() (response response.ResponseWithData, err error) {
	data, err := mysql.GetUserLeaderboard()
	if err != nil {
		response.Code = resp_code.GetUserRankError
		zap.L().Error("services-GetUserLeaderboard-GetUserLeaderboard ", zap.Error(err))
		return response, err
	}
	response.Code = resp_code.Success
	response.Data = data
	//fmt.Println(data)
	return response, nil
}
