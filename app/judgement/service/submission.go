package service

import (
	"context"
	"fmt"
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
	input := request.Input
	expected := request.Expected
	timeLimit := request.TimeLimit
	memoryLimit := request.MemoryLimit
	fmt.Println(uid)
	fmt.Println(code)
	fmt.Println(input)
	fmt.Println(expected)
	fmt.Println(timeLimit)
	fmt.Println(memoryLimit)
	return nil
}
