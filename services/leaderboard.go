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
		zap.L().Error("services-GetUserLeaderboard ", zap.Error(err))
		response.Code = resp.GetUserRankError
		return response, err
	}
	response.Code = resp.Success
	response.Data = data
	//fmt.Println(data)
	return response, nil
}
