package service

import (
	"fmt"
	"online-judge/app/judgement/service/judging/cpp"
	"online-judge/app/judgement/service/judging/golang"
	"online-judge/pkg/resp"
	pb "online-judge/proto"
)

func LanguageCheck(request *pb.SubmitRequest, response *pb.SubmitResponse) (err error) {
	language := request.Language
	if language == resp.CPP {
		response, err = cpp.JudgeCpp(request, response)
	}
	if language == resp.GO {
		response, err = golang.JudgeGO(request, response)
	}
	response.TotalNum = request.TotalNum
	fmt.Println("response: ", response.Status)
	return nil
}
