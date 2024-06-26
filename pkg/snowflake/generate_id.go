package snowflake

import (
	"github.com/yitter/idgenerator-go/idgen"
	"online_judge/pkg/define"
)

func Init() {
	var options = idgen.NewIdGeneratorOptions(1)
	options.BaseTime = define.BaseTime
	idgen.SetIdGenerator(options)
	return
}

func GetID() int64 {
	return idgen.NextId()
}
