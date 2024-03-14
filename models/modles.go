package models

import (
	"gorm.io/gorm"
	"time"
)

// User 用户基本信息
type User struct {
	gorm.Model
	UserID           int64     `gorm:"type:bigint;primaryKey" json:"user_id"`
	UserName         string    `gorm:"type:varchar(255);not null" json:"user_name"`
	Password         string    `gorm:"type:varchar(255);not null" json:"password"`
	Email            string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	RegistrationDate time.Time `gorm:"type:timestamp;not null" json:"registration_date"`
	LastLoginData    time.Time `gorm:"type:timestamp" json:"last_login_data"`
	Role             bool      `gorm:"type:boolean;not null" json:"role"`
	// true is Admin, false is user
}

// Problems 题目信息
type Problems struct {
	gorm.Model
	ProblemID  int64       `gorm:"type:bigint;primaryKey" json:"problem_id"`                                            // primary key
	MaxRuntime int64       `gorm:"type:bigint;not null" json:"max_runtime"`                                             // 时间限制
	MaxMemory  int64       `gorm:"type:bigint;not null" json:"max_memory"`                                              // 内存限制
	Title      string      `gorm:"type:varchar(255);not null" json:"title"`                                             // problem title
	Content    string      `gorm:"type:varchar(65535);not null" json:"content"`                                         // problem description
	Difficulty string      `gorm:"type:char(4);not null；check:difficulty IN ('easy', 'mid', 'hard')" json:"difficulty"` // easy mid hard
	PassNum    int         `gorm:"type:int;default:0" json:"pass_num"`                                                  // 通过测试样例数
	SubmitNum  int         `gorm:"type:int;default:0" json:"submit_num"`                                                // 提交数
	TestCases  []*TestCase `gorm:"foreignKey:ProblemID" json:"test_cases"`                                              // 测试样例集
}

// Submission 提交记录
type Submission struct {
	gorm.Model
	SubmissionID   int64
	UserID         int64
	ProblemID      int64
	Language       string
	Code           string
	SubmissionTime time.Duration
}

// Judgement 评测结果
type Judgement struct {
	gorm.Model
	SubmissionID int64
	Verdict      string
	Runtime      time.Duration
	MemoryUsage  int64
}

// ProgrammingLanguage 支持的编程语言信息
type ProgrammingLanguage struct {
	gorm.Model
	LanguageID int
	Name       string
	Version    string
}

// TestCase 测试样例
type TestCase struct {
	ID        int64
	ProblemID int64
	Input     string
	Expected  string
}

// SubmissionResult 测试样例集
type SubmissionResult struct {
	ID           int64
	SubmissionID int64
	TestCaseID   int64
	UserOutput   string
	Verdict      string
	Runtime      time.Duration
	MemoryUsage  int64
}
