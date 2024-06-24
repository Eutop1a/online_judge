package verify

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	redis2 "online_judge/dao/redis"
	"strconv"
	"strings"
	"time"
)

type CacheVerify struct {
}

var (
	Expired int64 = 600 //过期时间。单位：秒
)

// StoreVerifyCode 存储验证码到Redis并设置过期时间
func (c *CacheVerify) StoreVerifyCode(email, code string, timestamp int64) error {
	// 将时间戳转换为字符串
	tsStr := strconv.FormatInt(timestamp, 10)

	// 使用事务进行操作
	_, err := redis2.Client.TxPipelined(redis2.Ctx, func(pipe redis.Pipeliner) error {
		// 将验证码数据存储到哈希表中
		if err := pipe.HSet(redis2.Ctx, "VerificationDataMap", email, code+"_"+tsStr).Err(); err != nil {
			return err
		}

		// 设置过期时间
		if err := pipe.Expire(redis2.Ctx, "VerificationDataMap", time.Duration(Expired)*time.Second).Err(); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		zap.L().Error("redis-StoreVerifyCode-TxPipelined ", zap.Error(err))
		return err
	}
	// RDB持久化
	redis2.Client.BgSave(redis2.Ctx)

	return nil
}

// GetVerifyCode 从Redis获取验证码
func (c *CacheVerify) GetVerifyCode(email string) (string, error) {
	// 从哈希表中获取验证码数据
	result, err := redis2.Client.HGet(redis2.Ctx, "VerificationDataMap", email).Result()
	if err != nil {
		return "", err
	}

	// 解析验证码数据
	parts := strings.Split(result, "_")
	if len(parts) != 2 {
		zap.L().Error("redis-GetVerifyCode-Split " +
			fmt.Sprint("invalid data format"))
		return "", fmt.Errorf("invalid data format")
	}
	code := parts[0]
	ts, _ := strconv.ParseInt(parts[1], 10, 64)

	// 检查验证码是否过期
	if time.Now().Unix() > ts+Expired {
		zap.L().Error("redis-GetVerifyCode-Split " +
			fmt.Sprintf("verify code for email %s has expired ", email))
		return "", fmt.Errorf("verify code expired")
	}
	// 使用完之后删除
	_, err = redis2.Client.HDel(redis2.Ctx, "VerificationDataMap", email).Result()

	if err != nil {
		zap.L().Error("redis-GetVerifyCode-HDel " +
			fmt.Sprintf("failed to delete a record %s", email))
		return "", err
	}
	return code, nil
}

func (c *CacheVerify) StorePictureCode(username, code string, timestamp int64) error {
	// 将时间戳转换为字符串
	tsStr := strconv.FormatInt(timestamp, 10)

	// 使用事务进行操作
	_, err := redis2.Client.TxPipelined(redis2.Ctx, func(pipe redis.Pipeliner) error {
		// 将验证码数据存储到哈希表中
		if err := pipe.HSet(redis2.Ctx, "PictureCodeMap", username, code+"_"+tsStr).Err(); err != nil {
			return err
		}

		// 设置过期时间
		if err := pipe.Expire(redis2.Ctx, "PictureCodeMap", time.Duration(Expired)*time.Second).Err(); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		zap.L().Error("redis-StorePictureCode-TxPipelined " +
			fmt.Sprintf("store PictureCode to redis error: %v", err))
		return err
	}

	return nil
}

// GetPictureCode 从Redis获取验证码
func (c *CacheVerify) GetPictureCode(username string) (string, error) {

	// 从哈希表中获取验证码数据
	result, err := redis2.Client.HGet(redis2.Ctx, "PictureCodeMap", username).Result()
	if err != nil {
		return "", err
	}

	// 解析验证码数据
	parts := strings.Split(result, "_")
	if len(parts) != 2 {
		zap.L().Error("redis-GetPictureCode-Split " +
			fmt.Sprint("invalid data format"))
		return "", fmt.Errorf("invalid data format")
	}
	code := parts[0]
	ts, _ := strconv.ParseInt(parts[1], 10, 64)

	// 检查验证码是否过期
	if time.Now().Unix() > ts+Expired {
		zap.L().Error("redis-GetPictureCode-Expired " +
			fmt.Sprintf("picture code code for email %s has expired", username))
		return "", fmt.Errorf("PictureCode code expired")
	}
	// 使用完之后删除
	_, err = redis2.Client.HDel(redis2.Ctx, "PictureCodeMap", username).Result()

	if err != nil {
		zap.L().Error("redis-GetPictureCode-HDel " +
			fmt.Sprintf("failed to delete a record %s", username))
		return "", err
	}
	return code, nil
}

// StoreVerifyCode 存储验证码到Redis并设置过期时间
func StoreVerifyCode(email, code string, timestamp int64) error {
	// 将时间戳转换为字符串
	tsStr := strconv.FormatInt(timestamp, 10)

	// 使用事务进行操作
	_, err := redis2.Client.TxPipelined(redis2.Ctx, func(pipe redis.Pipeliner) error {
		// 将验证码数据存储到哈希表中
		if err := pipe.HSet(redis2.Ctx, "VerificationDataMap", email, code+"_"+tsStr).Err(); err != nil {
			return err
		}

		// 设置过期时间
		if err := pipe.Expire(redis2.Ctx, "VerificationDataMap", time.Duration(Expired)*time.Second).Err(); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		zap.L().Error("redis-StoreVerifyCode-TxPipelined ", zap.Error(err))
		return err
	}

	return nil
}

// GetVerifyCode 从Redis获取验证码
func GetVerifyCode(email string) (string, error) {
	// 从哈希表中获取验证码数据
	result, err := redis2.Client.HGet(redis2.Ctx, "VerificationDataMap", email).Result()
	if err != nil {
		return "", err
	}

	// 解析验证码数据
	parts := strings.Split(result, "_")
	if len(parts) != 2 {
		zap.L().Error("redis-GetVerifyCode-Split " +
			fmt.Sprint("invalid data format"))
		return "", fmt.Errorf("invalid data format")
	}
	code := parts[0]
	ts, _ := strconv.ParseInt(parts[1], 10, 64)

	// 检查验证码是否过期
	if time.Now().Unix() > ts+Expired {
		zap.L().Error("redis-GetVerifyCode-Split " +
			fmt.Sprintf("verify code for email %s has expired ", email))
		return "", fmt.Errorf("verify code expired")
	}
	// 使用完之后删除
	_, err = redis2.Client.HDel(redis2.Ctx, "VerificationDataMap", email).Result()

	if err != nil {
		zap.L().Error("redis-GetVerifyCode-HDel " +
			fmt.Sprintf("failed to delete a record %s", email))
		return "", err
	}
	return code, nil
}

func StorePictureCode(username, code string, timestamp int64) error {
	// 将时间戳转换为字符串
	tsStr := strconv.FormatInt(timestamp, 10)

	// 使用事务进行操作
	_, err := redis2.Client.TxPipelined(redis2.Ctx, func(pipe redis.Pipeliner) error {
		// 将验证码数据存储到哈希表中
		if err := pipe.HSet(redis2.Ctx, "PictureCodeMap", username, code+"_"+tsStr).Err(); err != nil {
			return err
		}

		// 设置过期时间
		if err := pipe.Expire(redis2.Ctx, "PictureCodeMap", time.Duration(Expired)*time.Second).Err(); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		zap.L().Error("redis-StorePictureCode-TxPipelined " +
			fmt.Sprintf("store PictureCode to redis error: %v", err))
		return err
	}

	return nil
}

// GetPictureCode 从Redis获取验证码
func GetPictureCode(username string) (string, error) {

	// 从哈希表中获取验证码数据
	result, err := redis2.Client.HGet(redis2.Ctx, "PictureCodeMap", username).Result()
	if err != nil {
		return "", err
	}

	// 解析验证码数据
	parts := strings.Split(result, "_")
	if len(parts) != 2 {
		zap.L().Error("redis-GetPictureCode-Split " +
			fmt.Sprint("invalid data format"))
		return "", fmt.Errorf("invalid data format")
	}
	code := parts[0]
	ts, _ := strconv.ParseInt(parts[1], 10, 64)

	// 检查验证码是否过期
	if time.Now().Unix() > ts+Expired {
		zap.L().Error("redis-GetPictureCode-Expired " +
			fmt.Sprintf("picture code code for email %s has expired", username))
		return "", fmt.Errorf("PictureCode code expired")
	}
	// 使用完之后删除
	_, err = redis2.Client.HDel(redis2.Ctx, "PictureCodeMap", username).Result()

	if err != nil {
		zap.L().Error("redis-GetPictureCode-HDel " +
			fmt.Sprintf("failed to delete a record %s", username))
		return "", err
	}
	return code, nil
}
