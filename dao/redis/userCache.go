package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"time"
)

type UserCache struct {
	redisClient *redis.Client
}

func NewUserCache(client *redis.Client) *UserCache {
	return &UserCache{
		redisClient: client,
	}
}

type Model struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// User 用户基本信息
type User struct {
	Model
	UserID           int64     `gorm:"type:bigint;primaryKey;column:userID" json:"user_id"`
	UserName         string    `gorm:"type:varchar(255);not null;column:userName" json:"user_name"`
	Password         string    `gorm:"type:varchar(255);not null;column:password" json:"password"`
	Email            string    `gorm:"type:varchar(255);unique;not null;column:email" json:"email"`
	RegistrationDate time.Time `gorm:"type:timestamp;not null;column:registrationDate" json:"registration_date"`
	LastLoginData    time.Time `gorm:"type:timestamp;column:lastLoginData" json:"last_login_data"`
	//Role             bool      `gorm:"type:boolean;not null;column:role" json:"role"`
	// true is Admin, false is user
}

func (u *UserCache) SetUser(user *User) error {
	redisKey := fmt.Sprintf("user:%d", user.UserID)
	err := u.redisClient.HMSet(Ctx, redisKey, map[string]interface{}{
		"UserID":           user.UserID,
		"UserName":         user.UserName,
		"Password":         user.Password,
		"Email":            user.Email,
		"RegistrationDate": user.RegistrationDate,
		"LastLoginData":    user.LastLoginData,
	}).Err()
	if err != nil {
		return err
	}
	return nil
}

func (u *UserCache) GetUser(userID int64) (*User, error) {
	redisKey := fmt.Sprintf("user:%d", userID)
	userInfo, err := u.redisClient.HGetAll(Ctx, redisKey).Result()
	if err != nil {
		return nil, err
	}
	registrationDate, _ := time.Parse("2006-01-02T15:04:05Z07:00", userInfo["RegistrationDate"])
	lastLoginData, _ := time.Parse("2006-01-02T15:04:05Z07:00", userInfo["LastLoginData"])
	// 将从 Redis 中获取到的数据组装成 User 结构

	user := &User{
		UserID:   userID,
		UserName: userInfo["UserName"],
		Password: userInfo["Password"],
		Email:    userInfo["Email"],
		// 报错：Cannot use 'userInfo["RegistrationDate"]' (type string) as the type time.Time
		RegistrationDate: registrationDate,
		LastLoginData:    lastLoginData,
	}
	return user, nil
}
func (u *UserCache) UpdateUser(user *User) error {
	// 更新用户信息的缓存逻辑
	redisKey := fmt.Sprintf("user:%d", user.UserID)
	userInfo := map[string]interface{}{
		"UserID":           user.UserID,
		"UserName":         user.UserName,
		"Password":         user.Password,
		"Email":            user.Email,
		"RegistrationDate": user.RegistrationDate.Format("2006-01-02T15:04:05Z07:00"),
		"LastLoginData":    user.LastLoginData.Format("2006-01-02T15:04:05Z07:00"),
		// 其他用户信息字段
	}
	err := u.redisClient.HMSet(Ctx, redisKey, userInfo).Err()
	if err != nil {
		return err
	}
	return nil
}

func (u *UserCache) DeleteUser(userID int64) error {
	// 删除用户信息的缓存逻辑
	redisKey := fmt.Sprintf("user:%d", userID)
	err := u.redisClient.Del(Ctx, redisKey).Err()
	if err != nil {
		return err
	}
	return nil
}
