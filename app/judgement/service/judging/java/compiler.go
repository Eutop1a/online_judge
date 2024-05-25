package java

import (
	"fmt"
	"os"
	"os/exec"
)

func Compiler(path string, fileName string) error {
	// 构建命令和参数
	command := "javac"
	arguments := []string{"-d", path, fmt.Sprintf("%s\\%s.java", path, fileName)}

	// 创建命令对象
	cmd := exec.Command(command, arguments...)

	// 设置输出和错误输出
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 运行命令
	err := cmd.Run()

	return err
}
