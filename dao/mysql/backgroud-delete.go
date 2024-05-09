package mysql

import (
	"time"
)

// BgDeleteMysql 后台删除 mysql 中 deleted_at 为 not null 的记录
func BgDeleteMysql() {
	ticker := time.NewTicker(5 * time.Hour)
	for {
		select {
		case <-ticker.C:
			BgDeleteUserTable()
			BgDeleteTestCaseTable()
			BgDeleteProblemsTable()
			BgDeleteSubmissionTable()
			BgDeleteJudgementTable()
		}
	}
}

// BgDeleteUserTable 删除mysql-onlinejudge-user表
func BgDeleteUserTable() {
	DB.Unscoped().Where("deleted_at is not null").Delete(&User{})
}

// BgDeleteProblemsTable 删除mysql-onlinejudge-problem表
func BgDeleteProblemsTable() {
	DB.Unscoped().Where("deleted_at is not null").Delete(&Problems{})
}

// BgDeleteTestCaseTable 删除mysql-onlinejudge-test_case表
func BgDeleteTestCaseTable() {
	DB.Unscoped().Where("deleted_at is not null").Delete(&TestCase{})
}

// BgDeleteJudgementTable 删除mysql-onlinejudge-judgement表
func BgDeleteJudgementTable() {
	DB.Unscoped().Where("deleted_at is not null").Delete(&Judgement{})
}

// BgDeleteSubmissionTable 删除mysql-onlinejudge-submission表
func BgDeleteSubmissionTable() {
	DB.Unscoped().Where("deleted_at is not null").Delete(&Submission{})
}
