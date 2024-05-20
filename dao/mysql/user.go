package mysql

import (
	"online-judge/pkg/utils"
)

// CheckEmail 检查是否有这个邮箱
func CheckEmail(email string, countEmail *int64) error {
	return DB.Model(&User{}).Where("email=?", email).Count(countEmail).Error
}

// CheckUsername 检查是否有这个用户名
func CheckUsername(username string, countUsername *int64) error {
	return DB.Model(&User{}).Where("username=?", username).Count(countUsername).Error
}

// CheckUserID 检查是否有这个用户ID
func CheckUserID(UID int64, countUserID *int64) error {
	return DB.Model(&User{}).Where("user_id=?", UID).Count(countUserID).Error
}

// InsertNewUser 插入用户
func InsertNewUser(uID int64, Username, password, email string) error {
	newUser := User{
		UserID:   uID,
		UserName: Username,
		Password: password,
		Email:    email,
	}

	return DB.Create(&newUser).Error
}

// CheckPwd 检查密码是否正确
func CheckPwd(username, plainText string) error {
	var checkTmp User
	//fmt.Println("in mysql package", username)
	err := DB.Model(&User{}).Where("username=?", username).First(&checkTmp).Error
	// 数据库搜索不到
	if err != nil {
		return err
	}
	err = utils.DecryptPwd(checkTmp.Password, plainText)
	// 密码错误
	if err != nil {
		return err
	}
	// 密码正确
	return nil
}

// GetUserDetail 获取用户详细信息
func GetUserDetail(UID int64) (data User, err error) {
	err = DB.Model(&User{}).Select("user_id, username, email").
		Where("user_id=?", UID).First(&data).Error
	return
}

// DeleteUser 删除用户
func DeleteUser(UID int64) error {
	return DB.Delete(&User{}, UID).Error // 根据主键删除
}

// UpdateUserDetail 更新用户信息
func UpdateUserDetail(UID int64, email, pwd string) (err error) {
	updateData := make(map[string]interface{})
	if email != "" {
		updateData["email"] = email
	}
	if pwd != "" {
		updateData["password"] = pwd
	}

	err = DB.Model(&User{}).Where("user_id=?", UID).Updates(updateData).Error
	return
}

// GetUserID 根据 username 获取用户ID
func GetUserID(username string) (uid int64, err error) {
	var user User
	err = DB.Take(&user, "username=?", username).Error
	if err != nil {
		return 0, err
	}
	return user.UserID, nil
}

// CheckUserIsAdmin 根据uid判断用户是不是管理员
func CheckUserIsAdmin(uid int64) (err error) {
	var admin Admin
	err = DB.Model(&Admin{}).Where("user_id=?", uid).First(&admin).Error

	return err
}

// AddAdminUser 添加管理员用户
func AddAdminUser(uid int64) (err error) {
	admin := Admin{
		UserID: uid,
	}
	err = DB.Model(&Admin{}).Create(&admin).Error
	return
}

func CheckAdminUserID(uid int64, countUserID *int64) error {
	return DB.Model(&Admin{}).Where("user_id=?", uid).Count(countUserID).Error
}
