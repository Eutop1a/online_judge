package resp

const (
	Accepted = 1000 + iota
	WrongAnswer
	ComplierError
	TimeLimited
	RuntimeError
	MemoryLimited
	SystemError
)
