package test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"online_judge/dao/mysql"
	"online_judge/models/admin/request"
	"online_judge/services"
	"online_judge/setting"
	"testing"
)

func init() {
	// 1. loading config files
	setting.Init()

	// 3. init mysql connection
	mysql.Init(setting.Conf.MySQLConfig)
}

func TestDeleteIfUidNotExist(t *testing.T) {
	err := mysql.DeleteUser(123456456)
	require.NoError(t, err)
}

func TestSetAdmin(t *testing.T) {
	var AdminService = services.ServiceGroupApp.AdminService

	var request = request.AdminAddSuperAdminReq{
		Username: "2021211705",
		Secret:   "eutop1a",
	}
	for i := 0; i < 10; i++ {
		resp := AdminService.AddSuperAdmin(request)
		fmt.Println(resp.Code)
	}

}

func TestGetProblemsByCategoryName(t *testing.T) {
	var ProblemService = services.ServiceGroupApp.ProblemService

	var categoryName = "aaaa"

	data, err := ProblemService.GetProblemListByCategory(categoryName)
	require.NoError(t, err)
	fmt.Println(data)

	categoryName = "string"

	data, err = ProblemService.GetProblemListByCategory(categoryName)
	require.NoError(t, err)
	fmt.Println(data)

}
