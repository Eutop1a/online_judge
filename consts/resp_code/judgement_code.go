package resp_code

const (
	Accepted = 1000 + iota
	WrongAnswer
	ComplierError
	TimeLimited
	RuntimeError
	MemoryLimited
	SystemError
)

// 位置
