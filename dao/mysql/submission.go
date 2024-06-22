package mysql

import "gorm.io/gorm"

// SaveSubmitCode 将提交记录保存在数据库
func SaveSubmitCode(submission *Submission) error {
	return DB.Create(submission).Error
}

// AddPassNum 题目AC，增加通过题目的数量
func AddPassNum(uid int64) error {
	return DB.Model(&User{}).Where("user_id = ?", uid).
		UpdateColumn("finish_num", gorm.Expr("finish_num + ?", 1)).Error
}
