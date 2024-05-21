package cpp

import (
	"fmt"
	"io"
	"online-judge/app/judgement/responses"
	pb "online-judge/proto"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"
)

func JudgeCpp(request *pb.SubmitRequest, response *pb.SubmitResponse) (*pb.SubmitResponse, error) {
	uid := request.UserId
	input := request.Input
	//code := request.Code 暂时不测存代码
	expected := request.Expected
	timeLimit := request.TimeLimit
	memoryLimit := request.MemoryLimit
	UID := strconv.FormatInt(uid, 10)
	/*savepath, err := utility.CodeSave(code, uid)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("savepath: ", savepath) 暂时不测存代码*/

	err := Complier(responses.Path, UID)
	if err != nil {
		fmt.Println("Complier Error: %v", err)
		response.Status = responses.ComplierError
		response.PassNum = 0
		return response, err
	}
	fmt.Println("Complier success")
	var WA = make(chan int)  //wrong answer
	var RE = make(chan int)  //Runtime Error
	var MLE = make(chan int) //Memory Limit Exceeded

	var lock sync.Mutex
	//var ALL = make(chan int) //
	var passCount = 0 //统计通过的样例个数
	//fmt.Println(len(input))
	for i := 0; i < len(input); i++ {
		go func() {
			// /path/uid.exe
			cmd := exec.Command(responses.Path + "/" + strconv.FormatInt(uid, 10) + ".exe")
			fmt.Println(responses.Path + "/" + strconv.FormatInt(uid, 10) + ".exe")
			stdin, err := cmd.StdinPipe()
			if err != nil {
				fmt.Println("Error creating stdin pipe:", err)
				return
			}
			defer stdin.Close()
			io.WriteString(stdin, input[i]+"\n")
			fmt.Println("test Input: ", input[i])
			//TODO:根据输入样例运行，拿到输出结果和标准输出结果进行比对
			var startMem runtime.MemStats //开始时内存情况
			runtime.ReadMemStats(&startMem)
			output, err := cmd.CombinedOutput()
			fmt.Println("out: ", string(output))
			if err != nil {
				fmt.Println("Runtime Error: ", err)
				RE <- 1
				return
			}
			var endMem runtime.MemStats //结束时内存情况
			runtime.ReadMemStats(&endMem)
			//TODO:处理答案错误

			if expected[i] != string(output) {
				fmt.Println("test Output: ", expected[i])
				WA <- 1
				return
			}
			//TODO:运行超内存
			// ÷1024是为了转化为KB
			if endMem.Alloc/1024-(startMem.Alloc/1024) > uint64(memoryLimit) {
				fmt.Println("Memory Usage: ", (endMem.TotalAlloc-startMem.TotalAlloc)/(1024))
				MLE <- 1
				return
			}
			lock.Lock() //该测试样例通过
			passCount++
			lock.Unlock()
			fmt.Println(passCount)
		}()
	}

	select {
	case <-RE:
		response.Status = responses.RuntimeError
	case <-WA:
		response.Status = responses.WrongAnswer
	case <-time.After(time.Millisecond * time.Duration(timeLimit)):
		if passCount == len(input) {
			response.Status = responses.Accepted //测试样例全部通过，表示正确
		} else {
			response.Status = responses.TimeLimited //超时
		}
	case <-MLE:
		response.Status = responses.MemoryLimited
	}
	fmt.Println("status: ", response.Status)
	return response, nil
}
