package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"online_judge/dao/mysql"
	"online_judge/dao/redis"
	"online_judge/logger"
	"online_judge/pkg/snowflake"
	"online_judge/router"
	"online_judge/setting"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title online_judge
// @version 1.0
// @description Refactoring
// @termsOfService http://swagger.io/terms/
// @contact.name Eutop1a
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 127.0.0.1:65533
// @BasePath /api/v1
func main() {
	// 1. loading config files
	if err := setting.Init(); err != nil {
		//zap.L().Error("main-setting-init error", zap.Error(err))
		fmt.Printf("init setting failed, err: %v\n", err)
		return
	}

	// 2. init logger
	if err := logger.Init(setting.Conf.LogConfig, setting.Conf.Mode); err != nil {
		//zap.L().Error("main-logger-init error", zap.Error(err))
		fmt.Printf("init logger failed, err: %v\n", err)
		return
	}
	defer zap.L().Sync()

	// 3. init mysql connection
	if err := mysql.Init(setting.Conf.MySQLConfig); err != nil {
		//zap.L().Error("main-mysql-init error", zap.Error(err))
		fmt.Printf("init mysql failed, err: %v\n", err)
		return
	}

	// 4. init redis connection
	if err := redis.Init(setting.Conf.RedisConfig); err != nil {
		//zap.L().Error("main-redis-init error", zap.Error(err))
		fmt.Printf("init redis failed, err: %v\n", err)
		return
	}
	defer redis.Close()

	// 雪花算法生成分布式ID
	snowflake.Init()

	//// 5. init rabbitmq connection
	//if err := mq.InitRabbitMQ(setting.Conf.RabbitMQConfig); err != nil {
	//	fmt.Printf("init rabbitmq failed, err: %v\n", err)
	//	return
	//}
	// 6. register route
	r := router.SetUpRouter(setting.Conf.Mode)

	err := r.Run(fmt.Sprintf(":%d", setting.Conf.Port))
	if err != nil {
		//zap.L().Error("main-Run error", zap.Error(err))
		fmt.Printf("run server failed, err: %v\n", err)
		return
	}

	// 7. start services
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("server Shutdown:", zap.Error(err))
	}
	zap.L().Info("server exiting")
}
