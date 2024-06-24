package cache

import (
	"online_judge/dao/redis/cache/admin"
	"online_judge/dao/redis/cache/auth"
	"online_judge/dao/redis/cache/leaderboard"
	"online_judge/dao/redis/cache/problem"
	"online_judge/dao/redis/cache/submission"
	"online_judge/dao/redis/cache/user"
	"online_judge/dao/redis/cache/verify"
)

type CacheGroup struct {
	CacheAdmin       admin.CacheGroup
	CacheAuth        auth.CacheGroup
	CacheLeaderboard leaderboard.CacheGroup
	CacheProblem     problem.CacheGroup
	CacheSubmission  submission.CacheGroup
	CacheUser        user.CacheGroup
	CacheVerify      verify.CacheGroup
}

var CacheGroupApp = new(CacheGroup)
