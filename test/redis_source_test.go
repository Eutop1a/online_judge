package test

import (
	"context"
	"fmt"
	"testing"
)

// 定义 Commander 接口
type Commander interface {
	Name() string
}

// 定义 ScanCmd 结构体，实现 Commander 接口
type ScanCmd struct {
	args []interface{}
}

func (cmd *ScanCmd) Name() string {
	return "scan"
}

// 定义 cmdable 类型，为一个函数类型
type cmdable func(ctx context.Context, cmd Commander) error

// NewScanCmd 函数，创建一个新的 ScanCmd 并调用 cmdable 函数
func NewScanCmd(ctx context.Context, c cmdable, args ...interface{}) *ScanCmd {
	cmd := &ScanCmd{args: args}
	_ = c(ctx, cmd)
	return cmd
}

// Scan 方法，使用 cmdable 类型来执行命令
func (c cmdable) Scan(ctx context.Context, cursor uint64, match string, count int64) *ScanCmd {
	args := []interface{}{"scan", cursor}
	if match != "" {
		args = append(args, "match", match)
	}
	if count > 0 {
		args = append(args, "count", count)
	}
	return NewScanCmd(ctx, c, args...)
}

// 实现一个 cmdable 函数，用于实际执行命令
func executeCommand(ctx context.Context, cmd Commander) error {
	fmt.Printf("Executing command: %s with args: %v\n", cmd.Name(), cmd.(*ScanCmd).args)
	return nil
}

func TestRedis(t *testing.T) {
	ctx := context.Background()
	client := cmdable(executeCommand)

	// 调用 Scan 方法
	cmd := client.Scan(ctx, 0, "pattern*", 10)
	fmt.Printf("Created ScanCmd: %v\n", cmd)
}
