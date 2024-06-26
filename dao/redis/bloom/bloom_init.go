package bloom

import (
	"fmt"
	"github.com/bits-and-blooms/bloom/v3"
	"go.uber.org/zap"
	"online_judge/dao/mysql"
	"online_judge/pkg/define"
	"strconv"
)

// ProblemDetailBloomFilter 针对于题目详细信息的布隆过滤器
var ProblemDetailBloomFilter *bloom.BloomFilter

// ProblemListBloomFilter 针对于题目列表的布隆过滤器
var ProblemListBloomFilter *bloom.BloomFilter

// InitBloomFilters 初始化布隆过滤器
func InitBloomFilters() {
	ReBuildBloomFilters()
}

func ReBuildBloomFilters() {
	// 初始化题目详细信息的布隆过滤器
	ProblemDetailBloomFilter = bloom.NewWithEstimates(nDetail, fpDetail)

	// 初始化题目列表的布隆过滤器
	ProblemListBloomFilter = bloom.NewWithEstimates(nList, fpList)

	// 从数据库加载所有题目到题目详细信息布隆过滤器
	problemList, err := mysql.GetAllProblem()
	if err != nil {
		zap.L().Error("init-GetAllProblemIDs", zap.Error(err))
		return
	}

	for _, problem := range problemList {
		cacheKeyDetail := fmt.Sprintf("%s:%s",
			define.GlobalCacheKeyMap.ProblemDetailPrefix,
			problem.ProblemID)
		//cacheKeyDetail := fmt.Sprintf("%s", problem.ProblemID)

		ProblemDetailBloomFilter.AddString(cacheKeyDetail)
	}

	// 从数据库加载所有分页信息并生成缓存键到题目列表布隆过滤器
	totalProblems := len(problemList) // 获取题目总数

	pageSize, _ := strconv.Atoi(define.DefaultSize) // 假设每页显示10条记录
	totalPages := (totalProblems + pageSize - 1) / pageSize

	for page := 1; page <= totalPages; page++ {
		cacheKey := fmt.Sprintf("%s:page-%d:size-%d",
			define.GlobalCacheKeyMap.ProblemListPrefix,
			page, pageSize)
		ProblemListBloomFilter.AddString(cacheKey)
	}
}

//
//// ProblemDetailBloomFilter 针对于题目详细信息的布隆过滤器
//var ProblemDetailBloomFilter *bloom.BloomFilter
//
//// ProblemListBloomFilter 针对于题目列表的布隆过滤器
//var ProblemListBloomFilter *bloom.BloomFilter
//
//// ProblemTitleBloomFilter 针对题目标题的布隆过滤器
//var ProblemTitleBloomFilter *bloom.BloomFilter
//
//// InitBloomFilters 初始化布隆过滤器
//func InitBloomFilters() {
//	// 初始化题目详细信息的布隆过滤器
//	ProblemDetailBloomFilter = bloom.NewWithEstimates(nDetail, fpDetail)
//
//	// 初始化题目列表的布隆过滤器
//	ProblemListBloomFilter = bloom.NewWithEstimates(nList, fpList)
//
//	// 初始化题目标题的布隆过滤器
//	ProblemTitleBloomFilter = bloom.NewWithEstimates(nTitle, fpTitle)
//
//	// 从数据库加载所有题目到题目详细信息布隆过滤器
//	problemList, err := mysql.GetAllProblem()
//	if err != nil {
//		zap.L().Error("init-GetAllProblemIDs", zap.Error(err))
//		return
//	}
//
//	for _, problem := range problemList {
//		cacheKeyDetail := fmt.Sprintf("%s:%s",
//			define.GlobalCacheKeyMap.ProblemDetailPrefix,
//			problem.ProblemID)
//
//		ProblemDetailBloomFilter.Add([]byte(cacheKeyDetail))
//	}
//
//	// 从数据库加载所有分页信息并生成缓存键到题目列表布隆过滤器
//	totalProblems := len(problemList) // 获取题目总数
//
//	pageSize, _ := strconv.Atoi(define.DefaultSize) // 假设每页显示10条记录
//	totalPages := (totalProblems + pageSize - 1) / pageSize
//
//	for page := 1; page <= totalPages; page++ {
//		cacheKey := fmt.Sprintf("%s:%d:%d",
//			define.GlobalCacheKeyMap.ProblemListPrefix,
//			page, pageSize)
//		ProblemListBloomFilter.Add([]byte(cacheKey))
//	}
//}
