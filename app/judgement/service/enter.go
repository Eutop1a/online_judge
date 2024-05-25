package service

import (
	"fmt"
	"online-judge/app/judgement/service/judging/cpp"
	"online-judge/app/judgement/service/judging/golang"
	"online-judge/consts"
	pb "online-judge/proto"
)

func LanguageCheck(request *pb.SubmitRequest, response *pb.SubmitResponse) (err error) {
	language := request.Language
	if language == consts.CPP {
		response, err = cpp.JudgeCpp(request, response)
	}
	if language == consts.GO {
		response, err = golang.JudgeGO(request, response)
	}
	response.TotalNum = request.TotalNum
	fmt.Println("response: ", response.Status)
	return nil
}
