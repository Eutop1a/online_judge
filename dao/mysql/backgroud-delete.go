package mysql

import (
	"time"
)

// BgDeleteMysql 后台删除 mysql 中 delete_at 为 not null 的记录
func BgDeleteMysql() {
	ticker := time.NewTicker(5 * time.Hour)
	for {
		select {
		case <-ticker.C:
			BgDeleteUserTable()
			BgDeleteProblemsTable()
			BgDeleteJudgementTable()
			BgDeleteSubmissionTable()
			BgDeleteTestCaseTable()
		}
	}
}

// BgDeleteRedis 后台删除 redis 中已经被使用过的 map 的记录
func BgDeleteRedis() {

}

// BgDeleteUserTable 删除mysql-onlinejudge-user表
func BgDeleteUserTable() {
	DB.Unscoped().Where("delete_at is not null").Delete(&User{})
}

// BgDeleteProblemsTable 删除mysql-onlinejudge-problem表
func BgDeleteProblemsTable() {
	DB.Unscoped().Where("delete_at is not null").Delete(&Problems{})
}

// BgDeleteJudgementTable 删除mysql-onlinejudge-judgement表
func BgDeleteJudgementTable() {
	DB.Unscoped().Where("delete_at is not null").Delete(&Judgement{})
}

// BgDeleteSubmissionTable 删除mysql-onlinejudge-submission表
func BgDeleteSubmissionTable() {
	DB.Unscoped().Where("delete_at is not null").Delete(&Submission{})
}

// BgDeleteTestCaseTable 删除mysql-onlinejudge-test_case表
func BgDeleteTestCaseTable() {
	DB.Unscoped().Where("delete_at is not null").Delete(&TestCase{})
}
