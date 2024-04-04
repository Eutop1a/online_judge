package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

var (
	DefaultPage       = "1"
	DefaultSize       = "20"
	Expired     int64 = 60000 //过期时间。单位：秒
)

// StoreVerificationCode 存储验证码到Redis并设置过期时间
func StoreVerificationCode(email, code string, timestamp int64) error {
	// 将时间戳转换为字符串
	tsStr := strconv.FormatInt(timestamp, 10)

	// 使用事务进行操作
	_, err := Client.TxPipelined(Ctx, func(pipe redis.Pipeliner) error {
		// 将验证码数据存储到哈希表中
		if err := pipe.HSet(Ctx, "VerificationDataMap", email, code+"_"+tsStr).Err(); err != nil {
			return err
		}

		// 设置过期时间
		if err := pipe.Expire(Ctx, "VerificationDataMap", time.Duration(Expired)*time.Second).Err(); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		zap.L().Error(fmt.Sprintf("store VerificationDataMap to redis error: %v", err))
		return err
	}
	// RDB持久化
	Client.BgSave(Ctx)

	return nil
}

// GetVerificationCode 从Redis获取验证码
func GetVerificationCode(email string) (string, error) {
	// 从哈希表中获取验证码数据
	result, err := Client.HGet(Ctx, "VerificationDataMap", email).Result()
	if err != nil {
		return "", err
	}

	// 解析验证码数据
	parts := strings.Split(result, "_")
	if len(parts) != 2 {
		zap.L().Error(fmt.Sprint("GetVerificationCode in redis error: invalid data format"))
		return "", fmt.Errorf("invalid data format")
	}
	code := parts[0]
	ts, _ := strconv.ParseInt(parts[1], 10, 64)

	// 检查验证码是否过期
	if time.Now().Unix() > ts+Expired {
		zap.L().Error(fmt.Sprintf("Verification code for email %s has expired", email))
		return "", fmt.Errorf("verification code expired")
	}

	return code, nil
}
