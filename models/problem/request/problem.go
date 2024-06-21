package request

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type GetProblemListReq struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`

	RedisClient *redis.Client
	Ctx         context.Context
}

type GetProblemDetailReq struct {
	ProblemID string `json:"problem_id" form:"problem_id"`

	RedisClient *redis.Client
	Ctx         context.Context
}

type GetProblemIDReq struct {
	Title string `json:"title" form:"title"`

	RedisClient *redis.Client
	Ctx         context.Context
}

type GetProblemRandomReq struct {
	RedisClient *redis.Client
	Ctx         context.Context
}

type SearchProblemReq struct {
	Msg         string `json:"msg" form:"msg"`
	RedisClient *redis.Client
	Ctx         context.Context
}
