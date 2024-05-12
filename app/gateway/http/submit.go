package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"online-judge/idl/pb"
)

func SubmissionHandler(c *gin.Context) {
	var taskReq pb.SubmitRequest
	c.JSON(http.StatusOK, gin.H{
		"uid":          taskReq.UserId,
		"code":         taskReq.Code,
		"input":        taskReq.Input,
		"expected":     taskReq.Expected,
		"time_limit":   taskReq.TimeLimit,
		"memory_limit": taskReq.MemoryLimit,
	})
}
