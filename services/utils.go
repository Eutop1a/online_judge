package services

import (
	"context"
	"github.com/go-redis/redis/v8"
	"online-judge/dao/mysql"
)

func convertTestCases(testCases []*TestCase) []*mysql.TestCase {
	// 提前转换类型
	var convertedTestCases []*mysql.TestCase
	for _, tc := range testCases {
		// 进行类型转换
		convertedTestCases = append(convertedTestCases, &mysql.TestCase{
			TID:      tc.TID,
			PID:      tc.PID,
			Input:    tc.Input,
			Expected: tc.Expected,
		})
	}
	return convertedTestCases
}

func deleteCacheByPrefix(redisClient *redis.Client, prefix string) error {
	ctx := context.Background()
	iter := redisClient.Scan(ctx, 0, prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		if err := redisClient.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}
