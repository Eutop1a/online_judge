package service

import (
	"online-judge/app/judgement/service/judging"
	"online-judge/pkg/resp"
	pb "online-judge/proto"
)

var (
	status int
)

func LanguageCheck(request *pb.SubmitRequest, response *pb.SubmitResponse) (err error) {
	language := request.Language
	if language == resp.CPP {
		status, err = judging.JudgeCpp(request)
	}

}
