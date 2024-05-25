package utility

import (
	"fmt"
	"os"
)

func CodeSave(code string, path string, language string) error {
	//TODO:使用绝对路径存放代码路径
	err := os.MkdirAll(path, 0777)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Successfully created directories")
	}
	dirName := path + "/main" + language
	//TODO:以dirName为文件名创建文件
	problemFile, err := os.OpenFile(dirName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	// 写入CPP文件
	_, err = problemFile.WriteString(code)
	if err != nil {
		// 写入文件失败时，关闭文件并返回错误
		problemFile.Close()
		return err
	}
	return err
}
