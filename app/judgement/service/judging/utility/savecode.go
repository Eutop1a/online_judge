package utility

import (
	"os"
	"strconv"
)

func CodeSave(code string, UID int64) (string, error) {
	//TODO:使用绝对路径存放代码路径
	path := `..\..\..\temp`
	dirName := path + "\\" + strconv.FormatInt(UID, 10) + ".cpp"
	//TODO:以dirName为文件名创建文件
	problemFile, err := os.OpenFile(dirName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return dirName, err
	}

	// 写入CPP文件
	_, err = problemFile.WriteString(code)
	if err != nil {
		// 写入文件失败时，关闭文件并返回错误
		problemFile.Close()
		return dirName, err
	}
	return dirName, err
}
