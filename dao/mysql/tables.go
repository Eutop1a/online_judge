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

// `gorm:"foreignKey:关联表的结构体字段;references:当前表的结构体字段;`

// User 用户基本信息
type User struct {
	Model
	UserID           int64  `gorm:"type:bigint;primaryKey;column:user_id" json:"user_id"`
	FinishProblemNum int64  `gorm:"type:int(11);default:0;column:finish_num" json:"finish_num"`
	UserName         string `gorm:"type:varchar(255);not null;column:username;uniqueIndex" json:"username"`
	Password         string `gorm:"type:varchar(255);not null;column:password" json:"password"`
	Email            string `gorm:"type:varchar(255);not null;column:email;uniqueIndex" json:"email"`
	Role             bool   `gorm:"type:boolean;not null;column:role" json:"role"`
	// true is Admin, false is user
}

//// Admin 管理员表 从用户表定位到ID，再来这里找
//type Admin struct {
//	Model
//	UserName string `gorm:"type:varchar(255);primaryKey;column:username" json:"username"`
//	//UserID int64 `gorm:"type:bigint;primaryKey;column:user_id" json:"user_id"`
//}

// ProblemWithFile 题目信息
type ProblemWithFile struct {
	Model
	ID           int    `gorm:"autoIncrement;column:id" json:"id"`                                    // primary key
	ProblemID    string `gorm:"type:char(36);uniqueIndex;column:problem_id" json:"problem_id"`        // unique key
	Title        string `gorm:"type:varchar(255);not null;column:title" json:"title"`                 // problem title
	Content      string `gorm:"type:text;not null;column:content" json:"content"`                     // problem description
	Difficulty   string `gorm:"type:char(4);not null;column:difficulty" json:"difficulty"`            // easy mid hard
	MaxRuntime   int    `gorm:"type:bigint;not null;column:max_runtime" json:"max_runtime"`           // 时间限制
	MaxMemory    int    `gorm:"type:bigint;not null;column:max_memory" json:"max_memory"`             // 内存限制
	InputPath    string `gorm:"type:varchar(255);not null;column:input_path" json:"input_path"`       // 输入文件路径
	ExpectedPath string `gorm:"type:varchar(255);not null;column:expected_path" json:"expected_path"` // 期望输出文件路径

	TestCases []*TestCaseWithFile `gorm:"foreignKey:PID;references:ProblemID" json:"test_cases"` // 测试样例集
}

// TestCaseWithFile 输入输出为文件格式的测试样例
type TestCaseWithFile struct {
	Model
	TID          string `gorm:"type:char(36);column:tid" json:"tid"`
	PID          string `gorm:"type:char(36);not null;column:pid" json:"pid"`                 // 对应的题目ID
	InputPath    string `gorm:"type:text;not null;column:input_path" json:"input_path"`       // 输入文件
	ExpectedPath string `gorm:"type:text;not null;column:expected_path" json:"expected_path"` // 期望输出文件名
}

// Problems 题目信息
type Problems struct {
	Model
	ID                int                `gorm:"autoIncrement;column:id" json:"id"`                                         // primary key
	ProblemID         string             `gorm:"type:char(36);uniqueIndex;column:problem_id" json:"problem_id"`             // unique key
	ProblemCategories []*ProblemCategory `gorm:"foreignKey:ProblemIdentity;references:ProblemID" json:"problem_categories"` // 分类表
	Title             string             `gorm:"type:varchar(255);not null;column:title;uniqueIndex" json:"title"`          // problem title
	Content           string             `gorm:"type:text;not null;column:content" json:"content"`                          // problem description
	Difficulty        string             `gorm:"type:char(4);not null;column:difficulty" json:"difficulty"`                 // easy mid hard
	MaxRuntime        int                `gorm:"type:bigint;not null;column:max_runtime" json:"max_runtime"`                // 时间限制
	MaxMemory         int                `gorm:"type:bigint;not null;column:max_memory" json:"max_memory"`                  // 内存限制

	TestCases []*TestCase `gorm:"foreignKey:PID;references:ProblemID" json:"test_cases"` // 测试样例集
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
	UserID         int64     `gorm:"type:bigint;foreignKey:UserID;references:UserID;column:user_id" json:"user_id"`               //用户ID
	SubmissionID   string    `gorm:"type:char(36);primaryKey;column:submission_id" json:"submission_id"`                          // 提交ID
	ProblemID      string    `gorm:"type:char(36);foreignKey:ProblemID;references:ProblemID;column:problem_id" json:"problem_id"` //题目ID
	Language       string    `gorm:"type:varchar(16);column:language" json:"language"`                                            //编程语言
	Code           string    `gorm:"type:text;column:code" json:"code"`                                                           // 代码
	SubmissionTime time.Time `gorm:"type:timestamp;column:submission_time" json:"submission_time"`                                // 提交时间
}

// Judgement 评测结果
type Judgement struct {
	Model
	UID          int64  `gorm:"type:bigint;foreignKey:references:UID;references:UserID;column:user_id" json:"user_id"`
	JudgementID  string `gorm:"type:char(36);primaryKey;column:judgement_id" json:"judgement_id"`                                        // 评测ID
	SubmissionID string `gorm:"type:char(36);foreignKey:SubmissionID;references:SubmissionID;column:submission_id" json:"submission_id"` //提交记录
	ProblemID    string `gorm:"type:char(36);foreignKey:ProblemID;references:ProblemID;column:problem_id" json:"problem_id"`
	Verdict      string `gorm:"type:varchar(20);column:verdict" json:"verdict"`      // 评测结果
	MemoryUsage  int    `gorm:"type:bigint;column:memory_usage" json:"memory_usage"` // 内存用量
	Runtime      int    `gorm:"type:bigint;not null;column:runtime" json:"runtime"`  // 运行时间
	Output       string `gorm:"type:text;column:output" json:"output"`               // 错误信息比对输出
}

type ProblemCategory struct {
	Model
	ProblemIdentity  string    `gorm:"type:char(36);column:problem_id;not null" json:"problem_id"`
	CategoryIdentity string    `gorm:"type:char(36);column:category_id;not null" json:"category_id"`
	Category         *Category `gorm:"foreignKey:CategoryIdentity;references:CategoryID"`
	//关联分类表 Category:ProblemCategory.CategoryIdentity->Category.CategoryID
}

type Category struct {
	Model
	Name       string `gorm:"type:varchar(255);not null;uniqueIndex;column:name" json:"name"`
	ParentName string `gorm:"type:varchar(100);column:parent_name" json:"parent_name" ` //父级ID
	CategoryID string `gorm:"type:varchar(36);primaryKey;column:category_id" json:"category_id"`
}

// TableName gorm自动改表名
func (u *User) TableName() string {
	return "user"
}

func (p *Problems) TableName() string {
	return "problems"
}

func (p *ProblemWithFile) TableName() string {
	return "problem_with_file"
}

func (p *TestCaseWithFile) TableName() string {
	return "test_case_with_file"
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

func (p *ProblemCategory) TableName() string {
	return "problem_category"
}

func (s *Category) TableName() string {
	return "category"
}
