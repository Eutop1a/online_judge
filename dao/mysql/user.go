package mysql

import (
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"online_judge/consts/resp_code"
	"online_judge/pkg/utils"
)

// CheckEmail 检查是否有这个邮箱
func CheckEmail(email string) (bool, error) {
	var count int64
	err := DB.Model(&User{}).Where("email=?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckUsername 检查是否有这个用户名
func CheckUsername(username string) (bool, error) {
	var count int64
	err := DB.Model(&User{}).Where("username=?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckEmailAndUsername 检查是否有这个邮箱和用户名
func CheckEmailAndUsername(email, username string, countEmail, countUsername *int64) error {
	err := DB.Model(&User{}).Where("email = ?", email).Count(countEmail).Error
	if err != nil {
		return err
	}

	err = DB.Model(&User{}).Where("username = ?", username).Count(countUsername).Error
	return err
}

// CheckUserID 检查是否有这个用户ID
func CheckUserID(userID int64) (bool, error) {
	var count int64
	err := DB.Model(&User{}).Where("user_id=?", userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// InsertNewUser 插入用户
func InsertNewUser(user *User) error {
	return DB.Create(&user).Error
}

// IsDuplicateEntryError 检查错误是否是唯一约束冲突错误
func IsDuplicateEntryError(err error) bool {
	var mysqlError *mysql.MySQLError
	if errors.As(err, &mysqlError) {
		// 主键冲突
		return mysqlError.Number == ErrMySQLDuplicateEntry
	}
	return false
}

func CheckUserCredentials(username, password string) (*User, error) {
	// 提取用户名和密码
	var user User
	err := DB.Model(&User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, resp_code.ErrInvalidCredentials
		}
		return nil, err
	}
	// 检查密码是否正确
	if !utils.CheckPwd(password, user.Password) {
		return nil, resp_code.ErrInvalidCredentials
	}

	return &user, nil
}

//// CheckPwd 检查密码是否正确
//func CheckPwd(username, plainText string) error {
//	var checkTmp User
//	//fmt.Println("in mysql package", username)
//	err := DB.Model(&User{}).Where("username=?", username).First(&checkTmp).Error
//	// 数据库搜索不到
//	if err != nil {
//		return err
//	}
//	if !utils.CheckPwd(checkTmp.Password, plainText) {
//		return false
//	}
//	// 密码错误
//	if err != nil {
//		return err
//	}
//	// 密码正确
//	return nil
//}

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
func UpdateUserDetail(UID int64, username, email, pwd string) (err error) {
	updateData := make(map[string]interface{})
	if username != "" {
		updateData["username"] = username
	}
	if email != "" {
		updateData["email"] = email
	}
	if pwd != "" {
		updateData["password"] = pwd
	}

	return DB.Model(&User{}).Where("user_id=?", UID).Updates(updateData).Error
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
	return DB.Model(&Admin{}).Where("user_id=?", uid).First(&admin).Error
}

// CheckUserIsAdminByUsername 根据username判断用户是不是管理员
func CheckUserIsAdminByUsername(username string) (err error) {
	var admin Admin
	return DB.Model(&Admin{}).Where("username=?", username).First(&admin).Error
}

//// AddAdminUserByUserId 添加管理员用户
//func AddAdminUserByUserId(uid int64) (err error) {
//	admin := Admin{
//		UserID: uid,
//	}
//	return DB.Model(&Admin{}).Create(&admin).Error
//}

// AddAdminUserByUsername 添加管理员用户
func AddAdminUserByUsername(username string) (err error) {
	admin := Admin{
		UserName: username,
	}
	return DB.Model(&Admin{}).Create(&admin).Error
}

func CheckAdminUserID(uid int64, countUserID *int64) error {
	return DB.Model(&Admin{}).Where("user_id=?", uid).Count(countUserID).Error
}

func CheckAdminUsername(username string, countUsername *int64) error {
	return DB.Model(&Admin{}).Where("username=?", username).Count(countUsername).Error
}

func CheckUsernameAndAdminExists(username string) (userExists, adminExists bool, err error) {
	var usernameCount int64

	// 检查用户名是否在user表中
	err = DB.Model(&User{}).Where("username=?", username).Count(&usernameCount).Error
	if err != nil {
		return false, false, err
	}
	userExists = usernameCount > 0

	// 检查用户名是否在admin表中
	usernameCount = 0
	err = DB.Model(&Admin{}).Where("username=?", username).Count(&usernameCount).Error
	if err != nil {
		return false, false, err
	}
	adminExists = usernameCount > 0

	return userExists, adminExists, nil
}

// SetAdminByUsername 根据用户名设置为管理员
func SetAdminByUsername(username string) (err error) {
	result := DB.Model(&User{}).Where("username = ?", username).
		Update("role", true)
	// 有错误
	if result.Error != nil {
		return result.Error
	}
	// 没有被影响的行，说明username不存在
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

// SetAdminByUsername 根据用户名设置为管理员
func _SetAdminByUsername(username string) (err error) {
	// 开启事务
	tx := DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Printf("Transaction rollback due to panic: %v\n", r)
		}
	}()

	// 锁定行
	var user User
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("username = ?", username).First(&user).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// 更新角色
	result := tx.Model(&user).Update("role", true)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return ErrUseAlreadyRoot
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
