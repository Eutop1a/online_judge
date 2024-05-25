package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mime/multipart"
	"online-judge/pkg/resp"
	"os"
	"path/filepath"
)

func SaveFile(c *gin.Context, fileHeader []*multipart.FileHeader, dstDir string) {
	// 保存输出文件
	err := os.MkdirAll(dstDir, os.ModePerm)
	if err != nil {
		zap.L().Error("controller-SaveFile-MkdirAll", zap.Error(err))
		resp.ResponseError(c, resp.CodeInternalServerError)
		return
	}
	for _, file := range fileHeader {
		dst := filepath.Join(dstDir, file.Filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			zap.L().Error("controller-SaveFile-SaveUploadedFile", zap.Error(err))
			resp.ResponseError(c, resp.CodeInternalServerError)
			return
		}
	}
}
