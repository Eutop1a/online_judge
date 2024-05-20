package cpp

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Complier(path string, file_name string) error {
	c := fmt.Sprintf("%s\\%s -std=c++11 -o %s\\%s", path, file_name+".cpp", path, file_name+".exe")
	cmd := exec.Command("g++", strings.Split(c, " ")...)
	fmt.Println("c is ", c)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err

}
