package responses

const (
	Accepted = 1000 + iota
	WrongAnswer
	CompilerError
	TimeLimited
	RuntimeError
	MemoryLimited
	SystemError
)

// Path 位置
const (
	//Path = "D:/online_judge/submit"
	Path = "./submit"
)
