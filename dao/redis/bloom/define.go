package bloom

var (
	// 预计插入的题目详细信息数量
	nDetail uint = 100000
	// ProblemDetailBloomFilter 的误判率
	fpDetail = 0.01

	// 预计插入的题目列表数量
	nList uint = 100000
	// ProblemListBloomFilter 的误判率
	fpList = 0.01

	// 预计插入的题目标题数量
	nTitle uint = 100000
	// ProblemTitleBloomFilter 的误判率
	fpTitle = 0.01
)
