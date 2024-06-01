package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"online-judge/server/extra/profiling"
	"online-judge/server/utils/prom"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"online-judge/server/consts/config"
	"online-judge/server/extra/tracing"
	"online-judge/server/utils/logging"
)

func main() {
	// 分布式链路追踪
	tp, err := tracing.SetTraceProvider(config.AuthRpcServerName)

	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panicf("Error to set the trace")
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("Error to set the trace")
		}
	}()

	// Configure Pyroscope
	profiling.InitPyroscope("online-judge.AuthService")

	log := logging.LogService(config.AuthRpcServerName)
	//lis, err := net.Listen("tcp", config.EnvCfg.PodIpAddr+config.AuthRpcServerPort)

	if err != nil {
		log.Panicf("Rpc %s listen happens error: %v", config.AuthRpcServerName, err)
	}

	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{
				0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)

	reg := prom.Client
	reg.MustRegister(srvMetrics)

}
