package service

import (
	"fmt"
	"online_judge/app/judgement/service/judging/cpp"
	"online_judge/app/judgement/service/judging/golang"
	"online_judge/app/judgement/service/judging/java"
	"online_judge/app/judgement/service/judging/python"
	"online_judge/consts"
	pb "online_judge/proto"
)

func LanguageCheck(request *pb.SubmitRequest, response *pb.SubmitResponse) (err error) {
	language := request.Language
	if language == consts.CPP {
		response, err = cpp.JudgeCpp(request, response)
	}
	if language == consts.GO {
		response, err = golang.JudgeGO(request, response)
	}
	if language == consts.JAVA {
		response, err = java.JudgeJAVA(request, response)
	}
	if language == consts.PYTHON {
		response, err = python.JudgePYTHON(request, response)
	}
	response.TotalNum = request.TotalNum
	fmt.Println("response: ", response.Status)
	return nil
}
