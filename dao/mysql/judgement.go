package mysql

import (
	"gorm.io/gorm"
)

// InsertNewSubmission 添加提交记录
func InsertNewSubmission(sub *Judgement) error {
	return DB.Model(sub).Create(sub).Error
}

// CheckIfAlreadyFinished 检查这个题目是否已经被解决
func CheckIfAlreadyFinished(uid int64, pid string) (finished bool, err error) {
	var tmp []Judgement
	err = DB.Model(&Judgement{}).Where("user_id = ? AND problem_id = ?", uid, pid).
		Find(&tmp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 没有找到匹配的记录
			return false, nil
		}
		// 处理其他错误
		return false, err
	}
	var count int
	for _, v := range tmp {
		if v.Verdict == "accepted" {
			count++
		}
	}
	if count == 1 {
		return false, nil
	} else {
		return true, nil
	}

}
