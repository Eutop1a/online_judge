package services

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"online-judge/consts"
	"online-judge/dao/mq"
	"online-judge/dao/mysql"
	"online-judge/idl/pb"
	"online-judge/pkg/resp"
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

func (s *Submission) SubmitCode() (response resp.Response) {
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
	var input, expected string
	for _, tc := range problemDetail.TestCases {
		input += tc.Input + "\n"
		expected += tc.Expected + "\n"
	}

	// 将需要的内容序列化
	data := pb.SubmitRequest{
		UserId:      s.UserID,
		Code:        s.Code,
		Input:       input,
		Expected:    expected,
		TimeLimit:   int32(problemDetail.MaxRuntime),
		MemoryLimit: int32(problemDetail.MaxMemory),
	}

	dataBody, err := json.Marshal(data)
	if err != nil {
		response.Code = resp.JSONMarshalError
		zap.L().Error("services-SubmitCode-Marshal ", zap.Error(err))
		return
	}

	// 发送给MQ的生产者
	err = mq.SendMessage2MQ(dataBody)
	if err != nil {
		response.Code = resp.Send2MQError
		zap.L().Error("services-SubmitCode-SendMessage2MQ ", zap.Error(err))
		return
	}

	msgs, err := mq.ConsumeMessage(context.Background(), consts.RabbitMQProblemQueueName)
	if err != nil {
		response.Code = resp.RecvFromMQError
		zap.L().Error("services-SubmitCode-ConsumeMessage ", zap.Error(err))
		return
	}
	var forever = make(chan struct{})
	go func() {
		for d := range msgs {
			var submitRequest pb.SubmitRequest
			err := json.Unmarshal(d.Body, &submitRequest)
			if err != nil {
				log.Printf("Failed to unmarshal JSON: %v", err)
				continue
			}
			d.Ack(false)
			// 执行judgement函数
			wg.Add(1)
			go s.Judgement(&submitRequest)
			wg.Wait()
		}
	}()

	log.Printf("Waiting for messages. To exit press CTRL+C")

	<-forever
	response.Code = resp.Success
	return
}

var wg sync.WaitGroup

func (s *Submission) Judgement(data *pb.SubmitRequest) (*pb.SubmitResponse, error) {
	conn, err := grpc.Dial("127.0.0.1:4000/api/v1/submission", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := pb.NewSubmissionClient(conn)

	response, err := client.SubmitCode(context.Background(), data)
	if err != nil {
		log.Printf("Failed to call gRPC method: %v", err)
		return nil, err
	}
	// 处理 Judgement 函数的返回结果
	// 这里假设 Judgement 函数返回的是一个 Response 对象
	var Response resp.Response
	if response.Status == "pass" {
		Response.Code = resp.Success
		fmt.Println("SUCCESS")
	} else {
		fmt.Println("ERROR")
	}
	wg.Done()
	return response, nil
}
