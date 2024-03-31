package mysql

import (
	"OnlineJudge/pkg"
	"time"
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

// InsertNewUser 插入用户
func InsertNewUser(uID int64, Username, password, email string, regData, loginData time.Time) error {
	// 将时间转换为 UTC 时间
	regDataUTC := regData.UTC()
	loginDataUTC := loginData.UTC()

	newUser := User{
		UserID:           uID,
		UserName:         Username,
		Password:         password,
		Email:            email,
		RegistrationDate: regDataUTC,
		LastLoginData:    loginDataUTC,
		Role:             false,
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

// UpdateLoginData 更新最后登录时间
func UpdateLoginData(username string, lastLoginTime time.Time) (T int, err error) {
	var updateTmp User
	if err = DB.Model(&User{}).First(&updateTmp).Where("userName=?", username).Error; err != nil {
		return 0, err
	}
	// 将时间转换为 UTC 时间
	formattedTime := lastLoginTime.UTC()

	if err = DB.Model(&User{}).Where("userName=?", username).
		Update("LastLoginData", formattedTime).Error; err != nil {
		return -1, err
	}

	return 1, nil
}
