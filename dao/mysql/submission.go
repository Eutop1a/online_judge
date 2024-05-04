package mysql

// SubmitCode 将提交记录保存在数据库
func SubmitCode(submission *Submission) error {
	return DB.Create(submission).Error
}
