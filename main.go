package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"online-judge/dao/mysql"
	"online-judge/dao/redis"
	"online-judge/logger"
	"online-judge/router"
	"online-judge/setting"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title online-judge
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
		zap.L().Error("main-setting-Init error", zap.Error(err))
		//fmt.Printf("Init setting failed, err: %v\n", err)
		return
	}

	// 2. init logger
	if err := logger.Init(setting.Conf.LogConfig, setting.Conf.Mode); err != nil {
		zap.L().Error("main-logger-Init error", zap.Error(err))
		//fmt.Printf("Init logger failed, err: %v\n", err)
		return
	}
	defer zap.L().Sync()

	// 3. init MYSQL connection
	if err := mysql.Init(setting.Conf.MySQLConfig); err != nil {
		zap.L().Error("main-mysql-Init error", zap.Error(err))
		//fmt.Printf("Init mysql failed, err: %v\n", err)
		return
	}

	// 4. init Redis connection
	if err := redis.Init(setting.Conf.RedisConfig); err != nil {
		zap.L().Error("main-redis-Init error", zap.Error(err))
		//fmt.Printf("Init redis failed, err: %v\n", err)
		return
	}

	// 5. register route
	r := router.SetUp(setting.Conf.Mode)

	err := r.Run(fmt.Sprintf(":%d", setting.Conf.Port))
	if err != nil {
		zap.L().Error("main-Run error", zap.Error(err))
		//fmt.Printf("Run server failed, err: %v\n", err)
		return
	}

	// 6. start services
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
