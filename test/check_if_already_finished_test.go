package test

import (
	"fmt"
	"online-judge/dao/mysql"
	"testing"
)

// CheckIfAlreadyFinished 检查这个题目是否已经被解决
func TestCheckIfAlreadyFinished(t *testing.T) {
	finished, err := mysql.CheckIfAlreadyFinished(1787406193556197376, "33637fe9-d549-45ca-9ab6-03cb0c96d1d3")
	fmt.Println(finished, err)

}
