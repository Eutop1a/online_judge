package test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"online_judge/pkg/snowflake"
	"strconv"
	"testing"
)

func init() {
	// 雪花算法生成分布式ID
	snowflake.Init()
}

func TestSnowflakeIDLen(t *testing.T) {
	// 生成唯一的ID
	for i := 0; i < 100; i++ {
		Id := snowflake.GetID()
		fmt.Println(Id)
		require.Equal(t, 14, len(strconv.FormatInt(Id, 10)))
	}

}
