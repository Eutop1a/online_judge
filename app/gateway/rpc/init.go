package rpc

import (
	"go-micro.dev/v4"
	"online-judge/app/gateway/wrappers"
	"online-judge/idl/pb"
)

var (
	SubmissionService pb.SubmissionService
)

func InitRPC() {
	userMicroService := micro.NewService(
		micro.Name("submissionService.client"),
		micro.WrapClient(wrappers.NewSubmitWrapper),
	)
	// 用户服务调用实例
	submissionService := pb.NewSubmissionService("rpcSubmissionService", userMicroService.Client())

	SubmissionService = submissionService
}
