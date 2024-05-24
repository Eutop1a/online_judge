package golang

import (
	"fmt"
	"os"
	"os/exec"
)

func Compiler(path, fileName string) error {
	// 使用 go build 命令编译 Go 代码
	cmd := exec.Command("go", "build", "-o",
		fmt.Sprintf("%s\\%s", path, fileName)+".exe", fmt.Sprintf("%s\\%s.go", path, fileName))
	cmd.Env = append(os.Environ(), "GO111MODULE=off") // 禁用 Go modules，如果不是使用 Go modules 的项目
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	return err
}
