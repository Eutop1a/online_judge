package mysql

import (
	"online-judge/pkg"
)

// CheckEmail 检查是否有这个邮箱
func CheckEmail(email string, countEmail *int64) error {
	err := DB.Model(&User{}).Where("email=?", email).Count(countEmail).Error
	return err
}

// CheckUsername 检查是否有这个用户名
func CheckUsername(username string, countUsername *int64) error {
	err := DB.Model(&User{}).Where("userName=?", username).Count(countUsername).Error
	return err
}

// CheckUserID 检查是否有这个用户ID
func CheckUserID(UID int64, countUsername *int64) error {
	err := DB.Model(&User{}).Where("userID=?", UID).Count(countUsername).Error
	return err
}

// InsertNewUser 插入用户
func InsertNewUser(uID int64, Username, password, email string) error {
	newUser := User{
		UserID:   uID,
		UserName: Username,
		Password: password,
		Email:    email,
	}

	err := DB.Create(&newUser).Error
	return err
}

// CheckPwd 检查密码是否正确
func CheckPwd(username, pwd string) (bool, error) {
	var checkTmp User
	err := DB.Model(&User{}).First(&checkTmp).Where("userName=?", username).Error
	if err != nil {
		return false, err
	}
	ok := pkg.DecryptPwd(checkTmp.Password, pwd)
	if !ok {
		return false, nil
	} else {
		return true, nil
	}
}

// GetUserDetail 获取用户详细信息
func GetUserDetail(UID int64) (data User, err error) {
	err = DB.Model(&User{}).Select("userID, userName, email, registrationDate, lastLoginData").
		First(&data).Where("userID=?", UID).Error
	return
}

// DeleteUser 删除用户
func DeleteUser(UID int64) error {
	err := DB.Delete(&User{}, UID).Error // 根据主键删除
	return err
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

	err = DB.Model(&User{}).Where("userID=?", UID).Updates(updateData).Error
	return
}

func GetUserID(username string) (uid int64, err error) {
	var user User
	err = DB.Take(&user, "userName=?", username).Error
	if err != nil {
		return 0, err
	}
	return user.UserID, nil
}
