package services

import (
	"go.uber.org/zap"
	"online-judge/dao/mysql"
	"online-judge/pkg/resp"
)

type Leaderboard struct{}

func (l *Leaderboard) GetUserLeaderboard() (response resp.ResponseWithData, err error) {
	data, err := mysql.GetUserLeaderboard()
	if err != nil {
		response.Code = resp.GetUserRankError
		zap.L().Error("services-GetUserLeaderboard-GetUserLeaderboard ", zap.Error(err))
		return response, err
	}
	response.Code = resp.Success
	response.Data = data
	//fmt.Println(data)
	return response, nil
}
