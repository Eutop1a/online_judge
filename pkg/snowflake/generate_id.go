package snowflake

import (
	"fmt"
	"github.com/sony/sonyflake"
	"time"
)

var (
	sonyFlake     *sonyflake.Sonyflake // 实例
	sonyMachineID uint16               // 机器ID
)

func getMachineID() (uint16, error) {
	return sonyMachineID, nil
}

func Init(machineID uint16) (err error) {
	sonyMachineID = machineID
	// 初始化一个开始的时间
	t, _ := time.Parse("2006-01-02", "2024-06-03")
	// 生成全局配置
	settings := sonyflake.Settings{
		StartTime: t,
		MachineID: getMachineID, // 指定机器ID
	}
	// 用配置生成sonyflake节点
	sonyFlake = sonyflake.NewSonyflake(settings)
	return
}

func GetID() (id uint64, err error) {
	if sonyFlake == nil {
		err = fmt.Errorf("snoy flake not inited")
		return
	}
	id, err = sonyFlake.NextID()
	return
}
