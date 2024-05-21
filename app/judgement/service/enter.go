package service

import (
	"fmt"
	"online-judge/app/judgement/service/judging/cpp"
	"online-judge/pkg/resp"
	pb "online-judge/proto"
)

var (
	status int
)

func LanguageCheck(request *pb.SubmitRequest, response *pb.SubmitResponse) (err error) {
	language := request.Language
	if language == resp.CPP {
		response, err = cpp.JudgeCpp(request, response)
	}
	fmt.Println("response: ", response.Status)
	return nil
}
