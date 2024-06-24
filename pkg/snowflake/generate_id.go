package snowflake

import (
	"github.com/yitter/idgenerator-go/idgen"
	"online_judge/pkg/common_define"
)

func Init() {
	var options = idgen.NewIdGeneratorOptions(1)
	options.BaseTime = common_define.BaseTime
	idgen.SetIdGenerator(options)
	return
}

func GetID() int64 {
	return idgen.NextId()
}
