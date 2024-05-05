package mysql

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// User 用户基本信息
type User struct {
	Model
	UserID           int64  `gorm:"type:bigint;primaryKey;column:user_id" json:"user_id"`
	FinishProblemNum int64  `gorm:"type:int(11);default:0;column:finish_num" json:"finish_num"`
	UserName         string `gorm:"type:varchar(255);not null;column:username" json:"username"`
	Password         string `gorm:"type:varchar(255);not null;column:password" json:"password"`
	Email            string `gorm:"type:varchar(255);not null;column:email" json:"email"`

	//Role             bool      `gorm:"type:boolean;not null;column:role" json:"role"`
	// true is Admin, false is user
}

// Admin 管理员表 从用户表定位到ID，再来这里找
type Admin struct {
	Model
	UserID int64 `gorm:"type:bigint;primaryKey;column:user_id" json:"user_id"`
}

// Problems 题目信息
type Problems struct {
	Model
	ProblemID  string `gorm:"type:char(36);primaryKey;column:problem_id" json:"problem_id"` // primary key
	Title      string `gorm:"type:varchar(255);not null;column:title" json:"title"`         // problem title
	Content    string `gorm:"type:text;not null;column:content" json:"content"`             // problem description
	Difficulty string `gorm:"type:char(4);not null;column:difficulty" json:"difficulty"`    // easy mid hard
	MaxRuntime int    `gorm:"type:bigint;not null;column:max_runtime" json:"max_runtime"`   // 时间限制
	MaxMemory  int    `gorm:"type:bigint;not null;column:max_memory" json:"max_memory"`     // 内存限制

	TestCases []*TestCase `gorm:"foreignKey:pid" json:"test_cases"` // 测试样例集
}

// TestCase 测试样例
type TestCase struct {
	Model
	TID      string `gorm:"type:char(36);column:tid" json:"tid"`
	PID      string `gorm:"type:char(36);not null;column:pid" json:"pid"` // 对应的题目ID
	Input    string `gorm:"type:text;column:input" json:"input"`          // 输入
	Expected string `gorm:"type:text;column:expected" json:"expected"`    // 期望输出
}

// Submission 提交记录
type Submission struct {
	Model
	UserID         int64     `gorm:"type:bigint;foreignKey:user_id;references:user(user_id);column:user_id" json:"user_id"`                   //用户ID
	SubmissionID   string    `gorm:"type:char(36);primaryKey;column:submission_id" json:"submission_id"`                                      // 提交ID
	ProblemID      string    `gorm:"type:char(36);foreignKey:problem_id;references:problems(problem_id);column:problem_id" json:"problem_id"` //题目ID
	Language       string    `gorm:"type:varchar(16);column:language" json:"language"`                                                        //编程语言
	Code           string    `gorm:"type:text;column:code" json:"code"`                                                                       // 代码
	SubmissionTime time.Time `gorm:"type:timestamp;column:submission_time" json:"submission_time"`                                            // 提交时间
}

// Judgement 评测结果
type Judgement struct {
	Model
	JudgementID  int64  `gorm:"type:bigint;primaryKey;column:judgement_id" json:"judgement_id"`                                                      // 评测ID
	SubmissionID int64  `gorm:"type:bigint;foreignKey:submission_id;references:submission(submission_id);column:submission_id" json:"submission_id"` //提交记录
	MemoryUsage  int64  `gorm:"type:bigint;column:memory_usage" json:"memory_usage"`                                                                 // 内存用量
	Verdict      string `gorm:"type:varchar(20);column:verdict" json:"verdict"`                                                                      // 评测结果
	Runtime      int64  `gorm:"type:bigint;not null;column:runtime" json:"runtime"`                                                                  // 运行时间
}

//
//// ProgrammingLanguage 支持的编程语言信息
//type ProgrammingLanguage struct {
//	Model
//	Language string `gorm:"type:varchar(10);not null;primaryKey;column:language" json:"language"` // 语言名称
//	Version  string `gorm:"type:varchar(20);not null;column:version" json:"version"`              // 语言版本
//}

// SubmissionResult 测试样例集
type SubmissionResult struct {
	Model
	SubmissionID int64  `gorm:"type:bigint;not null;column:submission_id" json:"submission_id"` // 提交记录
	TestCaseID   int64  `gorm:"type:bigint;not null;column:test_case_id" json:"test_case_id"`   // 测试样例ID
	MemoryUsage  int64  `gorm:"type:bigint;not null;column:memory_usage" json:"memory_usage"`   // 内存用量
	UserOutput   string `gorm:"type:text;column:user_output" json:"user_output"`                // 用户输出
	Verdict      string `gorm:"type:varchar(20);not null;column:verdict" json:"verdict"`        // 评测结果
	Runtime      int64  `gorm:"type:bigint;not null;column:runtime" json:"runtime"`             // 运行时间
}

// TableName gorm自动改表名
func (u *User) TableName() string {
	return "user"
}

func (p *Problems) TableName() string {
	return "problems"
}

func (t *TestCase) TableName() string {
	return "test_case"
}

func (s *Submission) TableName() string {
	return "submission"
}

func (j *Judgement) TableName() string {
	return "judgement"
}

//func (p *ProgrammingLanguage) TableName() string {
//	return "programming_language"
//}

func (s *SubmissionResult) TableName() string {
	return "submission_result"
}
