package services

import (
	"context"
	"fmt"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go.uber.org/zap"
	"online-judge/dao/mysql"
	"online-judge/pkg/resp"
	pb "online-judge/proto"
	"online-judge/setting"
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
	var idNum int64
	err := mysql.CheckUserID(s.UserID, &idNum)
	switch {
	case err != nil: // 搜索数据库错误
		response.Code = resp.SearchDBError
		zap.L().Error("services-SubmitCode-CheckUserID ", zap.Error(err))
		return
	case idNum == 0: // 用户id不存在
		response.Code = resp.NotExistUserID
		zap.L().Error("services-SubmitCode-CheckUserID " +
			fmt.Sprintf(" user_id %d do not exist ", s.UserID))
		return
	}
	idNum = 0
	// 检查题目id是否存在
	err = mysql.CheckProblemID(s.ProblemID, &idNum)
	switch {
	case err != nil: // 搜索数据库错误
		response.Code = resp.SearchDBError
		zap.L().Error("services-SubmitCode-CheckProblemID ", zap.Error(err))
		return
	case idNum == 0: // 题目id不存在
		response.Code = resp.ProblemNotExist
		zap.L().Error("services-SubmitCode-CheckProblemID " +
			fmt.Sprintf("problemID %s do not exist ", s.ProblemID))
		return
	}
	tmp := mysql.Submission{
		UserID:         s.UserID,
		SubmissionID:   s.SubmissionID,
		ProblemID:      s.ProblemID,
		Language:       s.Language,
		Code:           s.Code,
		SubmissionTime: s.SubmissionTime,
	}

	err = mysql.SubmitCode(&tmp)
	if err != nil {
		response.Code = resp.SearchDBError
		zap.L().Error("services-SubmitCode-SubmitCode ", zap.Error(err))
		return
	}
	// 获取全部的题目信息
	problemDetail, err := mysql.GetEntireProblem(s.ProblemID)
	if err != nil {
		response.Code = resp.SearchDBError
		zap.L().Error("services-SubmitCode-SubmitCode ", zap.Error(err))
		return
	}
	// 得到输入和输出
	var input, expected []string
	for _, tc := range problemDetail.TestCases {
		input = append(input, tc.Input)
		expected = append(expected, tc.Expected)
	}
	var language int32
	switch s.Language {
	case "Go":
		language = resp.GO
	case "Java":
		language = resp.JAVA
	case "C++":
		language = resp.CPP
	case "Python":
		language = resp.PYTHON
	default:
		response.Code = resp.UnsupportedLanguage
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
			response.Code = resp.InternalServerError
			return
		}
		wg.Done()
	}()
	wg.Wait()
	response.Code = resp.Success
	response.Data = resData
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
