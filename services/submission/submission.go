package submission

import (
	"context"
	"fmt"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go.uber.org/zap"
	"online_judge/consts"
	"online_judge/consts/resp_code"
	"online_judge/dao/mysql"
	"online_judge/models/common/response"
	"online_judge/models/submission/request"
	"online_judge/pkg/utils"
	pb "online_judge/proto"
	"online_judge/setting"
	"strconv"
	"strings"
	"sync"
)

type SubmissionService struct{}

func (s *SubmissionService) SubmitCode(request request.SubmissionReq) (response response.ResponseWithData) {
	// 直接在提交的时候通过外键判断 problemID 和 userID 是否存在
	err := mysql.SaveSubmitCode(&mysql.Submission{
		UserID:         request.UserID,
		SubmissionID:   request.SubmissionID,
		ProblemID:      request.ProblemID,
		Language:       request.Language,
		Code:           request.Code,
		SubmissionTime: request.SubmissionTime,
	})

	if err != nil {
		// 外键约束错误
		if mysql.IsForeignKeyConstraintError(err) {
			// userID 不存在
			if strings.Contains(err.Error(), "submission_user_id_fkey") {
				response.Code = resp_code.NotExistUserID
				zap.L().Error("service-SubmitCode-SaveSubmitCode",
					zap.String("message: user ID does not exist",
						strconv.FormatInt(request.UserID, 10)),
					zap.Error(err))
			}
			// problemID 不存在
			if strings.Contains(err.Error(), "submission_problem_id_fkey") {
				response.Code = resp_code.ProblemNotExist
				zap.L().Error("service-SubmitCode-SaveSubmitCode",
					zap.String("message: problem ID does not exist", request.ProblemID),
					zap.Error(err))
			}
		}
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-SaveSubmitCode-SaveSubmitCode ", zap.Error(err))
		return
	}

	// 获取全部的题目信息
	problemDetail, err := mysql.GetEntireProblem(request.ProblemID)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-SaveSubmitCode-SaveSubmitCode ", zap.Error(err))
		return
	}
	// 得到输入和输出
	var input, expected []string
	var total int
	for _, tc := range problemDetail.TestCases {
		input = append(input, tc.Input)
		expected = append(expected, tc.Expected)
	}
	// 获取总数
	total = len(input)
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
		Output:       resData.Output,
	})
	
	if err != nil {
		response.Code = resp_code.InsertToJudgementError
		zap.L().Error("services-SaveSubmitCode-InsertNewSubmission", zap.Error(err))
		return
	}

	// 如果AC，先判断是否已经完成，再直接增加通过题目数量
	finished, err := mysql.CheckIfAlreadyFinished(request.UserID, request.ProblemID)
	//fmt.Println("services-SaveSubmitCode-CheckIfAlreadyFinished:", finished, err)
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
		Output      string `json:"output"`
	}{
		UserId:      strconv.FormatInt(resData.UserId, 10),
		Status:      resData.Status,
		PassNum:     resData.PassNum,
		TotalNum:    resData.TotalNum,
		MemoryUsage: resData.MemoryUsage,
		Runtime:     resData.Runtime,
		Output:      resData.Output,
	}

	return
}

func (s *SubmissionService) Judgement(data *pb.SubmitRequest) (*pb.SubmitResponse, error) {
	// etcd 注册
	addr := setting.Conf.EtcdConfig.Host
	port := setting.Conf.EtcdConfig.Port
	etcdReg := etcd.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%d", addr, port)),
	)
	service := micro.NewService(
		micro.Registry(etcdReg),
	)
	service.Init()
	client := pb.NewSubmissionService("rpcSubmissionService", service.Client())

	resp, err := client.SubmitCode(context.Background(), data)
	if err != nil {
		zap.L().Error("services-SaveSubmitCode-SaveSubmitCode ", zap.Error(err))
		return nil, err
	}
	return resp, nil
}
