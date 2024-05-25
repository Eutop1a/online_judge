package services

import (
	"fmt"
	"go.uber.org/zap"
	"online-judge/consts"
	"online-judge/dao/mysql"
	"online-judge/pkg/resp"
	"online-judge/pkg/utils"
	"os"
	"path/filepath"
)

// DeleteUser 删除用户
func (u *UserService) DeleteUser() (response resp.Response) {
	// 检验是否有这个用户ID
	var UserIDCount int64
	err := mysql.CheckUserID(u.UserID, &UserIDCount)
	if err != nil {
		response.Code = consts.SearchDBError
		zap.L().Error("services-DeleteUser-CheckUserID ", zap.Error(err))
		return
	}
	if UserIDCount == 0 {
		response.Code = consts.NotExistUserID
		zap.L().Error("services-DeleteUser-CheckUserID "+
			fmt.Sprintf("do not have this userID %d ", u.UserID), zap.Error(err))
		return
	}
	// 删除用户
	err = mysql.DeleteUser(u.UserID)
	if err != nil {
		response.Code = consts.DBDeleteError
		zap.L().Error("services-DeleteUser-DeleteUser "+
			fmt.Sprintf("delete userID %d failed ", u.UserID), zap.Error(err))
		return
	}
	// 删除成功
	response.Code = consts.Success
	return
}

const SECRETCIPHER = "afd372788c1f7f646a678654901ce041ecc9012487dc0055b932cac9acaca27b6cf0488a3b5d0aa05022ab9a51e54b0e54e8188beaf4ad9cef517c0c76641f21"

func (u *UserService) AddSuperAdmin(secret string) (response resp.Response) {

	if utils.CryptoSecret(secret) != SECRETCIPHER {
		response.Code = consts.SecretError
		zap.L().Error("services-AddSuperAdmin-CryptoSecret " +
			fmt.Sprintf("secret error %d:%s", u.UserID, secret))
		return
	}
	var userIDCount int64
	err := mysql.CheckUserID(u.UserID, &userIDCount)
	if err != nil {
		response.Code = consts.SearchDBError
		zap.L().Error("services-AddSuperAdmin-CheckUserID ", zap.Error(err))
		return
	}
	if userIDCount == 0 {
		response.Code = consts.NotExistUserID
		zap.L().Error("services-AddSuperAdmin-CheckUserID "+
			fmt.Sprintf("do not have this userID %d ", u.UserID), zap.Error(err))
		return
	}
	// 检查这个ID是否已经存在于Admin数据库中
	userIDCount = 0
	err = mysql.CheckAdminUserID(u.UserID, &userIDCount)
	if err != nil {
		response.Code = consts.SearchDBError
		zap.L().Error("services-AddSuperAdmin-CheckAdminUserID ", zap.Error(err))
		return
	}
	if userIDCount != 0 {
		response.Code = consts.UserIDAlreadyExist
		zap.L().Error("services-AddSuperAdmin-CheckAdminUserID "+
			fmt.Sprintf("this userID %d already exist", u.UserID), zap.Error(err))
		return
	}
	// 添加到数据库中
	err = mysql.AddAdminUser(u.UserID)
	if err != nil {
		response.Code = consts.SearchDBError
		zap.L().Error("services-AddAdmin-AddAdminUser ", zap.Error(err))
		return
	}
	response.Code = consts.Success

	return
}

// AddAdmin 添加管理员
func (u *UserService) AddAdmin() (response resp.Response) {
	var UserIDCount int64
	err := mysql.CheckUserID(u.UserID, &UserIDCount)
	if err != nil {
		response.Code = consts.SearchDBError
		zap.L().Error("services-AddAdmin-CheckUserID ", zap.Error(err))
		return
	}
	if UserIDCount == 0 {
		response.Code = consts.NotExistUserID
		zap.L().Error("services-AddAdmin-CheckUserID "+
			fmt.Sprintf("do not have this userID %d ", u.UserID), zap.Error(err))
		return
	}
	err = mysql.AddAdminUser(u.UserID)
	if err != nil {
		response.Code = consts.SearchDBError
		zap.L().Error("services-AddAdmin-AddAdminUser ", zap.Error(err))
		return
	}
	response.Code = consts.Success

	return
}

// CreateProblem 创建题目
func (p *Problem) CreateProblem() (response resp.Response) {
	// 检查题目标题是否重复
	var problemNum int64
	err := mysql.CheckProblemTitle(p.Title, &problemNum)
	switch {
	case err != nil: // 搜索数据库错误
		response.Code = consts.SearchDBError
		zap.L().Error("services-CreateProblem-CheckProblemTitle ", zap.Error(err))
		return
	case problemNum > 0: // 题目已经存在
		response.Code = consts.ProblemAlreadyExist
		zap.L().Error("services-CreateProblem-CheckProblemTitle " +
			fmt.Sprintf("title %s aleardy exist", p.Title))
		return
	}
	var problems mysql.Problems
	problems.ProblemID = p.ProblemID
	problems.Title = p.Title
	problems.Content = p.Content
	problems.Difficulty = p.Difficulty
	problems.MaxRuntime = p.MaxRuntime
	problems.MaxMemory = p.MaxMemory

	// 提前转换类型
	var convertedTestCases []*mysql.TestCase
	for _, tc := range p.TestCases {
		// 进行类型转换
		convertedTestCases = append(convertedTestCases, &mysql.TestCase{
			TID:      tc.TID,
			PID:      tc.PID,
			Input:    tc.Input,
			Expected: tc.Expected,
		})
	}

	problems.TestCases = convertedTestCases
	// 创建题目
	err = mysql.CreateProblem(&problems)
	if err != nil {
		response.Code = consts.CreateProblemError
		zap.L().Error("services-CreateProblem-CreateProblem ", zap.Error(err))
		return
	}
	response.Code = consts.Success
	return
}

// UpdateProblem 更新题目
func (p *Problem) UpdateProblem() (response resp.Response) {
	// 检查题目id是否存在
	var idNum int64
	err := mysql.CheckProblemID(p.ProblemID, &idNum)
	switch {
	case err != nil: // 搜索数据库错误
		response.Code = consts.SearchDBError
		zap.L().Error("services-UpdateProblem-CheckProblemID ", zap.Error(err))
		return
	case idNum == 0: // 题目id不存在
		response.Code = consts.ProblemNotExist
		zap.L().Error("services-UpdateProblem-CheckProblemID " +
			fmt.Sprintf("problemID %s do not exist", p.ProblemID))
		return
	}

	// 检查题目标题是否存在
	var titleNum int64
	err = mysql.CheckProblemTitle(p.Title, &titleNum)
	switch {
	case err != nil: // 搜索数据库错误
		response.Code = consts.SearchDBError
		zap.L().Error("services-UpdateProblem-CheckProblemTitle", zap.Error(err))
		return
	case titleNum != 0: // 题目title已存在
		response.Code = consts.ProblemAlreadyExist
		zap.L().Error("services-UpdateProblem-CheckProblemTitle" +
			fmt.Sprintf("problem title %s already exist", p.Title))
		return
	}

	var problems mysql.Problems
	problems.ProblemID = p.ProblemID
	problems.Title = p.Title
	problems.Content = p.Content
	problems.Difficulty = p.Difficulty
	problems.MaxRuntime = p.MaxRuntime
	problems.MaxMemory = p.MaxMemory

	// 提前转换类型
	var convertedTestCases []*mysql.TestCase
	for _, tc := range p.TestCases {

		// 进行类型转换
		convertedTestCases = append(convertedTestCases, &mysql.TestCase{
			TID:      tc.TID,
			PID:      tc.PID,
			Input:    tc.Input,
			Expected: tc.Expected,
		})

	}
	problems.TestCases = convertedTestCases

	err = mysql.UpdateProblem(&problems)
	if err != nil {
		zap.L().Error("services-UpdateProblem-UpdateProblem ", zap.Error(err))
		response.Code = consts.InternalServerError
		return
	}
	response.Code = consts.Success
	return
}

// DeleteProblem 删除题目
func (p *Problem) DeleteProblem() (response resp.Response) {
	// 检查题目id是否存在
	var idNum int64
	err := mysql.CheckProblemID(p.ProblemID, &idNum)
	switch {
	case err != nil: // 搜索数据库错误
		response.Code = consts.SearchDBError
		zap.L().Error("services-DeleteProblem-CheckProblemID ", zap.Error(err))
		return
	case idNum == 0: // 题目id不存在
		response.Code = consts.ProblemNotExist
		zap.L().Error("services-DeleteProblem-CheckProblemID " +
			fmt.Sprintf("problemID %s do not exist", p.ProblemID))
		return
	}
	// 删除题目
	err = mysql.DeleteProblem(p.ProblemID)
	if err != nil {
		response.Code = consts.SearchDBError
		zap.L().Error("services-DeleteProblem-DeleteProblem  ", zap.Error(err))
		return
	}
	response.Code = consts.Success
	return
}

// CreateProblemWithFile 创建测试用例为file的题目
func (p *Problem) CreateProblemWithFile() (response resp.Response) {
	// 检查题目标题是否重复
	var problemNum int64
	err := mysql.CheckProblemIDWithFile(p.Title, &problemNum)
	switch {
	case err != nil: // 搜索数据库错误
		response.Code = consts.SearchDBError
		zap.L().Error("services-CreateProblemWithFile-CheckProblemTitle ", zap.Error(err))
		return
	case problemNum > 0: // 题目已经存在
		response.Code = consts.ProblemAlreadyExist
		zap.L().Error("services-CreateProblemWithFile-CheckProblemTitle " +
			fmt.Sprintf("title %s aleardy exist", p.Title))
		return
	}

	var problems mysql.ProblemWithFile
	problems.ProblemID = p.ProblemID
	problems.Title = p.Title
	problems.Content = p.Content
	problems.Difficulty = p.Difficulty
	problems.MaxRuntime = p.MaxRuntime
	problems.MaxMemory = p.MaxMemory
	problems.InputPath = p.InputDst
	problems.ExpectedPath = p.ExpectedDst

	// 创建题目
	err = mysql.CreateProblemWithFile(&problems)
	if err != nil {
		response.Code = consts.CreateProblemError
		zap.L().Error("services-CreateProblemWithFile-CreateProblemWithFile ", zap.Error(err))
		return
	}
	response.Code = consts.Success
	return
}

// DeleteProblemWithFile 删除题目
func (p *Problem) DeleteProblemWithFile() (response resp.Response) {
	// 检查题目id是否存在
	var idNum int64
	err := mysql.CheckProblemIDWithFile(p.ProblemID, &idNum)
	switch {
	case err != nil: // 搜索数据库错误
		response.Code = consts.SearchDBError
		zap.L().Error("services-DeleteProblemWithFile-CheckProblemIDWithFile ", zap.Error(err))
		return
	case idNum == 0: // 题目id不存在
		response.Code = consts.ProblemNotExist
		zap.L().Error("services-DeleteProblemWithFile-CheckProblemIDWithFile " +
			fmt.Sprintf("problemID %s do not exist", p.ProblemID))
		return
	}
	// 删除题目
	problemID, err := mysql.DeleteProblemWithFile(p.ProblemID)
	if err != nil {
		response.Code = consts.SearchDBError
		zap.L().Error("services-DeleteProblemWithFile-DeleteProblemWithFile ", zap.Error(err))
		return
	}
	//fmt.Println("path: ", filepath.Join(consts.FilePath, problemID))
	err = os.RemoveAll(filepath.Join(consts.FilePath, problemID))

	if err != nil {
		response.Code = consts.RemoveTestFileError
		zap.L().Error("services-DeleteProblemWithFile-Remove ", zap.Error(err))
		return
	}
	response.Code = consts.Success
	return
}

// DeleteProblemTestCaseWithFile 删除题目测试用例文件
func (p *Problem) DeleteProblemTestCaseWithFile() (response resp.Response) {
	// 检查题目id是否存在
	var idNum int64
	err := mysql.CheckProblemIDWithFile(p.ProblemID, &idNum)
	switch {
	case err != nil: // 搜索数据库错误
		response.Code = consts.SearchDBError
		zap.L().Error("services-DeleteProblemTestCaseWithFile-CheckProblemID ", zap.Error(err))
		return
	case idNum == 0: // 题目id不存在
		response.Code = consts.ProblemNotExist
		zap.L().Error("services-DeleteProblemTestCaseWithFile-CheckProblemID " +
			fmt.Sprintf("problemID %s do not exist", p.ProblemID))
		return
	}

	// 检查题目标题是否存在
	var titleNum int64
	err = mysql.CheckProblemTitleWithFile(p.Title, &titleNum)
	switch {
	case err != nil: // 搜索数据库错误
		response.Code = consts.SearchDBError
		zap.L().Error("services-UpdateProblemWithFile-CheckProblemTitle", zap.Error(err))
		return
	case titleNum != 0: // 题目title已存在
		response.Code = consts.ProblemAlreadyExist
		zap.L().Error("services-UpdateProblemWithFile-CheckProblemTitle" +
			fmt.Sprintf("problem title %s already exist", p.Title))
		return
	}
	err = os.RemoveAll(filepath.Join(consts.FilePath, p.ProblemID))

	if err != nil {
		response.Code = consts.RemoveTestFileError
		zap.L().Error("services-DeleteProblemTestCaseWithFile-Remove ", zap.Error(err))
		return
	}
	response.Code = consts.Success
	return
}

// UpdateProblemWithFile 更新题目
func (p *Problem) UpdateProblemWithFile() (response resp.Response) {

	var problems mysql.ProblemWithFile
	problems.ProblemID = p.ProblemID
	problems.Title = p.Title
	problems.Content = p.Content
	problems.Difficulty = p.Difficulty
	problems.MaxRuntime = p.MaxRuntime
	problems.MaxMemory = p.MaxMemory
	problems.InputPath = p.InputDst
	problems.ExpectedPath = p.ExpectedDst

	err := mysql.UpdateProblemWithFile(&problems)
	if err != nil {
		zap.L().Error("services-UpdateProblemWithFile-UpdateProblemWithFile ", zap.Error(err))
		response.Code = consts.InternalServerError
		return
	}
	response.Code = consts.Success
	return
}
