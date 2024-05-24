package cpp

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Compiler(path string, fileName string) error {
	c := fmt.Sprintf("%s/%s -std=c++11 -o %s/%s", path, fileName+".cpp", path, fileName+".exe")
	cmd := exec.Command("g++", strings.Split(c, " ")...)
	fmt.Println("c is ", c)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err

}
