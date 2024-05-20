package cpp

import (
	"fmt"
	"online-judge/pkg/resp"
	pb "online-judge/proto"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
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
	path := `..\..\..\temp`
	err := Complier(path, UID)
	if err != nil {
		fmt.Println("Complier Error")
		response.Status = resp.ComplierError
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

	for i := 0; i < len(input); i++ {
		go func() {
			// /path/uid.exe
			cmd := exec.Command(path, "\\"+strconv.FormatInt(uid, 10), ".exe")
			cmd.Stdin = strings.NewReader(input[i])
			fmt.Println("test Input: ", input[i])

			//TODO:根据输入样例运行，拿到输出结果和标准输出结果进行比对
			var startMem runtime.MemStats //开始时内存情况
			runtime.ReadMemStats(&startMem)
			output, err := cmd.CombinedOutput()
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
				fmt.Println("out: ", string(output))
				WA <- 1
				return
			}
			//TODO:运行超内存
			// ÷1024是为了转化为KB
			if endMem.Alloc/1024-(startMem.Alloc/1024) > uint64(memoryLimit) {

				MLE <- 1
				return
			}
			lock.Lock() //该测试样例通过
			passCount++
			lock.Unlock()
		}()
	}

	select {
	case <-RE:
		response.Status = resp.RuntimeError
	case <-WA:
		response.Status = resp.WrongAnswer
	case <-time.After(time.Millisecond * time.Duration(timeLimit)):
		if passCount == len(input) {
			response.Status = resp.Accepted //测试样例全部通过，表示正确
		} else {
			response.Status = resp.TimeLimited //超时
		}
	case <-MLE:
		response.Status = resp.MemoryLimited
	}
	fmt.Println("status: ", response.Status)
	return response, nil
}
