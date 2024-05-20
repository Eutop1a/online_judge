package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"online-judge/pkg/resp"
	pb "online-judge/proto"
	"sync"
)

type SubmitSrv struct {
}

var UserSrvIns *SubmitSrv
var UserSrvOnce sync.Once

// GetSubmitSrv 懒汉式的单例模式 lazy-loading
func GetSubmitSrv() *SubmitSrv {
	UserSrvOnce.Do(func() {
		UserSrvIns = &SubmitSrv{}
	})
	return UserSrvIns
}

func (s SubmitSrv) SubmitCode(ctx context.Context, request *pb.SubmitRequest, response *pb.SubmitResponse) error {
	//TODO implement me
	uid := request.UserId
	code := request.Code
	language := request.Language
	input := request.Input
	expected := request.Expected
	timeLimit := request.TimeLimit
	memoryLimit := request.MemoryLimit
	fmt.Println(uid)
	fmt.Println(code)
	fmt.Println(language)
	fmt.Println(input)
	fmt.Println(expected)
	fmt.Println(timeLimit)
	fmt.Println(memoryLimit)
	response.Status = resp.Accepted
	response.PassNum = 10
	response.UserId = request.UserId
	err := LanguageCheck(request, response)
	if err != nil {
		zap.L().Error("judgement-service-language-check-failed", zap.Error(err))
	}
	return nil
}
