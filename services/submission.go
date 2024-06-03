package services

import (
	"context"
	"fmt"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go.uber.org/zap"
	"online-judge/consts"
	"online-judge/consts/resp_code"
	"online-judge/dao/mysql"
	"online-judge/pkg/resp"
	"online-judge/pkg/utils"
	pb "online-judge/proto"
	"online-judge/setting"
	"strconv"
	"sync"
	"time"
)

// Submission 提交记录
type Submission struct {
	SubmissionID   string    `form:"submission_id" json:"submission_id"`     // 提交ID
	UserID         int64     `form:"user_id" json:"user_id"`                 //用户ID
	ProblemID      string    `form:"problem_id" json:"problem_id"`           //题目ID
	Language       string    `form:"language" json:"language"`               //编程语言
	Code           string    `form:"code" json:"code"`                       // 代码
	SubmissionTime time.Time `form:"submission_time" json:"submission_time"` // 提交时间
}

func (s *Submission) SubmitCode() (response resp.ResponseWithData) {
	// 检查用户id是否存在
	// 检验是否有这个用户ID
	exist, err := mysql.CheckUserID(s.UserID)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-DeleteUser-CheckUserID ", zap.Error(err))
		return
	}
	if !exist {
		response.Code = resp_code.NotExistUserID
		zap.L().Error("services-DeleteUser-CheckUserID "+
			fmt.Sprintf("do not have this userID %d ", s.UserID), zap.Error(err))
		return
	}
	// 检查题目id是否存在
	exists, err := mysql.CheckProblemIDExists(s.ProblemID)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-UpdateProblem-CheckProblemID ", zap.Error(err))
		return
	}
	if !exists {
		response.Code = resp_code.ProblemNotExist
		zap.L().Error("services-UpdateProblem-CheckProblemID " +
			fmt.Sprintf("problemID %s do not exist", s.ProblemID))
		return
	}

	err = mysql.SubmitCode(&mysql.Submission{
		UserID:         s.UserID,
		SubmissionID:   s.SubmissionID,
		ProblemID:      s.ProblemID,
		Language:       s.Language,
		Code:           s.Code,
		SubmissionTime: s.SubmissionTime,
	})

	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-SubmitCode-SubmitCode ", zap.Error(err))
		return
	}
	// 获取全部的题目信息
	problemDetail, err := mysql.GetEntireProblem(s.ProblemID)
	if err != nil {
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-SubmitCode-SubmitCode ", zap.Error(err))
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
	switch s.Language {
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
		UserId:      s.UserID,
		Language:    language,
		Code:        s.Code,
		Input:       input,
		Expected:    expected,
		TimeLimit:   int32(problemDetail.MaxRuntime),
		MemoryLimit: int32(problemDetail.MaxMemory),
		TotalNum:    int32(total),
	}
	//
	//dataBody, err := json.Marshal(data)
	//if err != nil {
	//	response.Code = resp.JSONMarshalError
	//	zap.L().Error("services-SubmitCode-Marshal ", zap.Error(err))
	//	return
	//}
	//
	//// 发送给MQ的生产者
	//err = mq.SendMessage2MQ(dataBody)
	//if err != nil {
	//	response.Code = resp.Send2MQError
	//	zap.L().Error("services-SubmitCode-SendMessage2MQ ", zap.Error(err))
	//	return
	//}
	//// 消费者
	//msgs, err := mq.ConsumeMessage(context.Background(), consts.RabbitMQProblemQueueName)
	//if err != nil {
	//	response.Code = resp.RecvFromMQError
	//	zap.L().Error("services-SubmitCode-ConsumeMessage ", zap.Error(err))
	//	return
	//}

	//var resData *pb.SubmitResponse
	//var forever = make(chan struct{})
	//go func() {
	//	for d := range msgs {
	//		var submitRequest pb.SubmitRequest
	//		err := json.Unmarshal(d.Body, &submitRequest)
	//		if err != nil {
	//			zap.L().Error("services-SubmitCode-Unmarshal ", zap.Error(err))
	//			continue
	//		}
	//		// 执行judgement函数
	//		resData, err = s.Judgement(&submitRequest)
	//		if err != nil {
	//			response.Code = resp.InternalServerError
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
		UID:          s.UserID,
		JudgementID:  utils.GetUUID(),
		SubmissionID: s.SubmissionID,
		ProblemID:    s.ProblemID,
		MemoryUsage:  int(resData.MemoryUsage),
		Verdict:      verdict,
		Runtime:      int(resData.Runtime),
	})
	if err != nil {
		response.Code = resp_code.InsertToJudgementError
		zap.L().Error("services-SubmitCode-InsertNewSubmission", zap.Error(err))
		return
	}

	// 如果AC，先判断是否已经完成，再直接增加通过题目数量
	finished, err := mysql.CheckIfAlreadyFinished(s.UserID, s.ProblemID)
	//fmt.Println("services-SubmitCode-CheckIfAlreadyFinished:", finished, err)
	if err != nil { // 查询数据库错误
		response.Code = resp_code.SearchDBError
		zap.L().Error("services-SubmitCode-CheckIfAlreadyFinished ", zap.Error(err))
		return
	}
	if finished { // 题目已经被完成
		zap.L().Error("services-SubmitCode-CheckIfAlreadyFinished " +
			fmt.Sprintf("%d had finished this problem %s", s.UserID, s.ProblemID))
	} else { // 题目还没有被完成过
		if resData.Status == resp_code.Accepted {
			err = mysql.AddPassNum(resData.UserId)
			if err != nil {
				zap.L().Error("services-SubmitCode-AddPassNum", zap.Error(err))
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

func (s *Submission) Judgement(data *pb.SubmitRequest) (*pb.SubmitResponse, error) {
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

	response, err := client.SubmitCode(context.Background(), data)
	if err != nil {
		zap.L().Error("services-SubmitCode-SubmitCode ", zap.Error(err))
		return nil, err
	}
	return response, nil
}
