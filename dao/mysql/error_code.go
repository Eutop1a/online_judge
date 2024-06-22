package mysql

import (
	"errors"
	"github.com/go-sql-driver/mysql"
)

const (
	// ErrMySQLDuplicateEntry uniqueIndex冲突
	ErrMySQLDuplicateEntry = 1062
	// ErrMySQLForeignKeyConstraint 外键约束冲突
	ErrMySQLForeignKeyConstraint = 1452
	ErrMySQLDupEntryWithKeyName  = 1586
)

// ErrCategoryAlreadyExist 分类已经存在
var ErrCategoryAlreadyExist = errors.New("category already exist")

// ErrUserNotFound 用户不存在
var ErrUserNotFound = errors.New("user not found")

// ErrUseAlreadyRoot 用户已是root
var ErrUseAlreadyRoot = errors.New("user already root")

// ErrCategoryNotFound 分类不存在
var ErrCategoryNotFound = errors.New("category not found")

// ErrTitleAlreadyExist 题目标题已经存在
var ErrTitleAlreadyExist = errors.New("title already exist")

// ErrProblemIDNotExist 题目ID不存在
var ErrProblemIDNotExist = errors.New("problem id does not exist")

// IsUniqueConstraintError 检查是否为唯一约束错误
func IsUniqueConstraintError(err error) bool {
	mysqlErr, ok := err.(*mysql.MySQLError)
	return ok && mysqlErr.Number == ErrMySQLDuplicateEntry
}

// IsForeignKeyConstraintError 检查是否为外键约束错误
func IsForeignKeyConstraintError(err error) bool {
	mysqlErr, ok := err.(*mysql.MySQLError)
	return ok && mysqlErr.Number == ErrMySQLForeignKeyConstraint // MySQL 外键约束错误代码
}
