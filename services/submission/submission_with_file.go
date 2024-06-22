package submission

import (
	"fmt"
	"go.uber.org/zap"
	"online_judge/consts"
	"online_judge/consts/resp_code"
	"online_judge/dao/mysql"
	"online_judge/models/common/response"
	"online_judge/models/submission/request"
	"online_judge/pkg/utils"
	pb "online_judge/proto"
	"strconv"
	"sync"
)

func (s *SubmissionService) SubmitCodeWithFile(request request.SubmissionReq) (response response.ResponseWithData) {
	// 检验是否有这个用户ID
	exists, err := mysql.CheckUserID(request.UserID)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-DeleteUser-CheckUserID ", zap.Error(err))
		return
	}
	if !exists {
		response.Code = resp_code.NotExistUserID
		zap.L().Error("services-DeleteUser-CheckUserID "+
			fmt.Sprintf("do not have this userID %d ", request.UserID), zap.Error(err))
		return
	}

	// 检查题目id是否存在
	exists, err = mysql.CheckProblemIDExists(request.ProblemID)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-UpdateProblem-CheckProblemID ", zap.Error(err))
		return
	}
	if !exists {
		response.Code = resp_code.ProblemNotExist
		zap.L().Error("services-UpdateProblem-CheckProblemID " +
			fmt.Sprintf("problemID %s do not exist", request.ProblemID))
		return
	}

	err = mysql.SaveSubmitCode(&mysql.Submission{
		UserID:         request.UserID,
		SubmissionID:   request.SubmissionID,
		ProblemID:      request.ProblemID,
		Language:       request.Language,
		Code:           request.Code,
		SubmissionTime: request.SubmissionTime,
	})

	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-SaveSubmitCode-SaveSubmitCode ", zap.Error(err))
		return
	}
	// 获取全部的题目信息
	problemDetail, err := mysql.GetEntireProblemWithFile(request.ProblemID)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-SaveSubmitCode-SaveSubmitCode ", zap.Error(err))
		return
	}
	// 得到输入和输出
	input, expected, err := utils.GetInputAndExpectedFromFile(problemDetail.InputPath, problemDetail.ExpectedPath)
	if err != nil {
		response.Code = resp_code.ReadTestFileError
		zap.L().Error("services-SubmitCodeWithFile-GetInputAndExpectedFromFile ", zap.Error(err))
		return
	}
	// 获取总数
	total := len(input)
	var language int32
	switch request.Language {
	case "Go":
		language = consts.GO
	case "Java":
		language = consts.JAVA
	case "C++":
		language = consts.CPP
	case "Python":
		language = consts.PYTHON
	default:
		response.Code = resp_code.UnsupportedLanguage
		return
	}
	// 将需要的内容序列化
	data := pb.SubmitRequest{
		UserId:      request.UserID,
		Language:    language,
		Code:        request.Code,
		Input:       input,
		Expected:    expected,
		TimeLimit:   int32(problemDetail.MaxRuntime),
		MemoryLimit: int32(problemDetail.MaxMemory),
		TotalNum:    int32(total),
	}
	//
	//dataBody, err := json.Marshal(data)
	//if err != nil {
	//	response.Code = response.JSONMarshalError
	//	zap.L().Error("services-SaveSubmitCode-Marshal ", zap.Error(err))
	//	return
	//}
	//
	//// 发送给MQ的生产者
	//err = mq.SendMessage2MQ(dataBody)
	//if err != nil {
	//	response.Code = response.Send2MQError
	//	zap.L().Error("services-SaveSubmitCode-SendMessage2MQ ", zap.Error(err))
	//	return
	//}
	//// 消费者
	//msgs, err := mq.ConsumeMessage(context.Background(), consts.RabbitMQProblemQueueName)
	//if err != nil {
	//	response.Code = response.RecvFromMQError
	//	zap.L().Error("services-SaveSubmitCode-ConsumeMessage ", zap.Error(err))
	//	return
	//}

	//var resData *pb.SubmitResponse
	//var forever = make(chan struct{})
	//go func() {
	//	for d := range msgs {
	//		var submitRequest pb.SubmitRequest
	//		err := json.Unmarshal(d.Body, &submitRequest)
	//		if err != nil {
	//			zap.L().Error("services-SaveSubmitCode-Unmarshal ", zap.Error(err))
	//			continue
	//		}
	//		// 执行judgement函数
	//		resData, err = s.Judgement(&submitRequest)
	//		if err != nil {
	//			response.Code = response.InternalServerError
	//			return
	//		}
	//		forever <- struct{}{}
	//		// 确认ACK
	//		d.Ack(false)
	//	}
	//}()
	//<-forever
	var resData *pb.SubmitResponse
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// 执行judgement函数
		resData, err = s.Judgement(&data)
		if err != nil {
			response.Code = resp_code.InternalServerError
			return
		}
		wg.Done()
	}()
	wg.Wait()

	var verdict string

	switch resData.Status {
	case resp_code.Accepted:
		verdict = "accepted"
	case resp_code.WrongAnswer:
		verdict = "wrong answer"
	case resp_code.ComplierError:
		verdict = "compiler error"
	case resp_code.TimeLimited:
		verdict = "time limited"
	case resp_code.MemoryLimited:
		verdict = "memory limited"
	case resp_code.RuntimeError:
		verdict = "runtime error"
	case resp_code.SystemError:
		verdict = "system error"
	default:
		verdict = "unknown"
	}

	// 先向数据库中添加这一次的提交记录
	err = mysql.InsertNewSubmission(&mysql.Judgement{
		UID:          request.UserID,
		JudgementID:  utils.GetUUID(),
		SubmissionID: request.SubmissionID,
		ProblemID:    request.ProblemID,
		MemoryUsage:  int(resData.MemoryUsage),
		Verdict:      verdict,
		Runtime:      int(resData.Runtime),
	})
	if err != nil {
		response.Code = resp_code.InsertToJudgementError
		zap.L().Error("services-SaveSubmitCode-InsertNewSubmission", zap.Error(err))
		return
	}

	// 如果AC，先判断是否已经完成，再直接增加通过题目数量
	finished, err := mysql.CheckIfAlreadyFinished(request.UserID, request.ProblemID)
	fmt.Println("services-SaveSubmitCode-CheckIfAlreadyFinished:", finished, err)
	if err != nil { // 查询数据库错误
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-SaveSubmitCode-CheckIfAlreadyFinished ", zap.Error(err))
		return
	}
	if finished { // 题目已经被完成
		zap.L().Error("services-SaveSubmitCode-CheckIfAlreadyFinished " +
			fmt.Sprintf("%d had finished this problem %s", request.UserID, request.ProblemID))
	} else { // 题目还没有被完成过
		if resData.Status == resp_code.Accepted {
			err = mysql.AddPassNum(resData.UserId)
			if err != nil {
				zap.L().Error("services-SaveSubmitCode-AddPassNum", zap.Error(err))
				response.Code = resp_code.SearchDBError
				return
			}
		}
	}

	response.Code = resp_code.Success
	response.Data = struct {
		UserId      string `json:"user_id"`
		Status      int32  `json:"status"`
		PassNum     int32  `json:"pass_num"`
		TotalNum    int32  `json:"total_num"`
		MemoryUsage int32  `json:"memory_usage"`
		Runtime     int32  `json:"runtime"`
	}{
		UserId:      strconv.FormatInt(resData.UserId, 10),
		Status:      resData.Status,
		PassNum:     resData.PassNum,
		TotalNum:    resData.TotalNum,
		MemoryUsage: resData.MemoryUsage,
		Runtime:     resData.Runtime,
	}

	return
}
