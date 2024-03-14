package models

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
	UserID           int64     `gorm:"type:bigint;primaryKey;column:userID" json:"user_id"`
	UserName         string    `gorm:"type:varchar(255);not null;column:userNameID" json:"user_name"`
	Password         string    `gorm:"type:varchar(255);not null;column:password" json:"password"`
	Email            string    `gorm:"type:varchar(255);unique;not null;column:email" json:"email"`
	RegistrationDate time.Time `gorm:"type:timestamp;not null;column:registrationDate" json:"registration_date"`
	LastLoginData    time.Time `gorm:"type:timestamp;column:lastLoginData" json:"last_login_data"`
	Role             bool      `gorm:"type:boolean;not null;column:role" json:"role"`
	// true is Admin, false is user
}

// Problems 题目信息
type Problems struct {
	Model
	ProblemID  int64       `gorm:"type:bigint;primaryKey;column:problemID" json:"problem_id"`                        // primary key
	MaxRuntime int64       `gorm:"type:bigint;not null;column:maxRuntime" json:"max_runtime"`                        // 时间限制
	MaxMemory  int64       `gorm:"type:bigint;not null;column:maxMemory" json:"max_memory"`                          // 内存限制
	Title      string      `gorm:"type:varchar(255);not null;column:title" json:"title"`                             // problem title
	Content    string      `gorm:"type:text;not null;column:content" json:"content"`                                 // problem description
	Difficulty string      `gorm:"type:char(4);not null；column:difficulty" json:"difficulty"`                        // easy mid hard
	PassNum    int         `gorm:"type:int;default:0;column:passNum" json:"pass_num"`                                // 通过测试样例数
	SubmitNum  int         `gorm:"type:int;default:0;column:submitNum" json:"submit_num"`                            // 提交数
	TestCases  []*TestCase `gorm:"foreignKey:PID;references:Problems(ProblemID);column:testCases" json:"test_cases"` // 测试样例集
}

// Submission 提交记录
type Submission struct {
	Model
	SubmissionID   int64     `gorm:"type:bigint;primaryKey;column:submissionID" json:"submission_id"`                                               // 提交ID
	UserID         int64     `gorm:"type:bigint;foreignKey:UserID;references:User(UserID);column:userID" json:"user_id"`                            //用户ID
	ProblemID      int64     `gorm:"type:bigint;foreignKey:ProblemID;references:Problems(ProblemID);column:problemID" json:"problem_id"`            //题目ID
	Language       string    `gorm:"type:varchar(16);foreignKey:Language;references:ProgrammingLanguage(Language);column:language" json:"language"` //编程语言
	Code           string    `gorm:"type:text;column:code" json:"code"`                                                                             // 代码
	SubmissionTime time.Time `gorm:"type:timestamp;column:submissionTime" json:"submission_time"`                                                   // 提交时间
}

// Judgement 评测结果
type Judgement struct {
	Model
	JudgementID  int64  `gorm:"type:bigint;primaryKey;column:judgementID" json:"judgement_id"`                                                    // 评测ID
	SubmissionID int64  `gorm:"type:bigint;foreignKey:SubmissionID;references:Submission(SubmissionID);column:submissionID" json:"submission_id"` //提交记录
	MemoryUsage  int64  `gorm:"type:bigint;column:memoryUsage" json:"memory_usage"`                                                               // 内存用量
	Verdict      string `gorm:"type:varchar(20);column:verdict" json:"verdict"`                                                                   // 评测结果
	Runtime      int64  `gorm:"type:bigint;not null" json:"runtime"`                                                                              // 运行时间
}

// ProgrammingLanguage 支持的编程语言信息
type ProgrammingLanguage struct {
	Model
	Language string `gorm:"type:varchar(10);not null;primaryKey;column:language" json:"language"` // 语言名称
	Version  string `gorm:"type:varchar(20);not null;column:version" json:"version"`              // 语言版本
}

// TestCase 测试样例
type TestCase struct {
	Model
	PID      int64  `gorm:"type:bigint;not null;column:PID" json:"PID"` // 对应的题目ID
	Input    string `gorm:"type:text;column:input" json:"input"`        // 输入
	Expected string `gorm:"type:text;column:expected" json:"expected"`  // 期望输出
}

// SubmissionResult 测试样例集
type SubmissionResult struct {
	Model
	SubmissionID int64  `gorm:"type:bigint;not null;column:submissionID" json:"submission_id"` // 提交记录
	TestCaseID   int64  `gorm:"type:bigint;not null;column:testCaseID" json:"test_case_id"`    // 测试样例ID
	MemoryUsage  int64  `gorm:"type:bigint;not null;column:memoryUsage" json:"memory_usage"`   // 内存用量
	UserOutput   string `gorm:"type:text;column:userOutput" json:"user_output"`                // 用户输出
	Verdict      string `gorm:"type:varchar(20);not null;column:verdict" json:"verdict"`       // 评测结果
	Runtime      int64  `gorm:"type:bigint;not null;column:runtime" json:"runtime"`            // 运行时间
}

// TableName gorm自动改表名
func (u *User) TableName() string {
	return "User"
}

func (p *Problems) TableName() string {
	return "Problems"
}

func (s *Submission) TableName() string {
	return "Submission"
}

func (j *Judgement) TableName() string {
	return "Judgement"
}

func (p *ProgrammingLanguage) TableName() string {
	return "ProgrammingLanguage"
}

func (t *TestCase) TableName() string {
	return "TestCase"
}

func (s *SubmissionResult) TableName() string {
	return "SubmissionResult"
}
