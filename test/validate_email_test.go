package test

import (
	"fmt"
	"regexp"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	email := `1905700640@qq.com`
	fmt.Println(len(email))
	regex := regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)
	//_ = regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)
	ok := regex.MatchString(email)
	if !ok {
		t.Errorf("Email does not match %s", email)
	} else {
		fmt.Println("Email matched")
	}
}
