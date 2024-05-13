package main

import (
	"fmt"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"online-judge/app/judgement/service"
	pb "online-judge/proto"
	"online-judge/setting"
)

func main() {
	//  loading config files
	if err := setting.Init(); err != nil {
		fmt.Printf("init setting failed, err: %v\n", err)
		return
	}

	// etcd 注册
	etcdReg := etcd.NewRegistry(
		registry.Addrs("127.0.0.1:2379"),
	)
	// new 一个微服务实例，使用gin暴露http接口并注册到etcd
	microService := micro.NewService(
		micro.Name("rpcSubmissionService"),
		micro.Address("127.0.0.1:8082"),
		micro.Registry(etcdReg),
	)
	// 结构命令行参数，初始化
	microService.Init()
	// 服务注册
	_ = pb.RegisterSubmissionHandler(microService.Server(), service.GetSubmitSrv())
	// 启动微服务
	_ = microService.Run()
}
