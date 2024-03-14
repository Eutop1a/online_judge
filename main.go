package main

import (
	"OnlineJudge/dao/mysql"
	"OnlineJudge/logger"
	"OnlineJudge/routes"
	"OnlineJudge/setting"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 1. loading config files
	if err := setting.Init(); err != nil {
		fmt.Printf("Init setting failed, err: %v\n", err)
		return
	}

	// 2. init logger
	if err := logger.Init(setting.Conf.LogConfig, setting.Conf.Mode); err != nil {
		fmt.Printf("Init logger failed, err: %v\n", err)
		return
	}
	defer zap.L().Sync()

	// 3. init MYSQL connection
	if err := mysql.Init(setting.Conf.MySQLConfig); err != nil {
		fmt.Printf("Init mysql failed, err: %v\n", err)
		return
	}

	// 5. register route
	r := routes.SetUp(setting.Conf.Mode)

	err := r.Run(fmt.Sprintf(":%d", setting.Conf.Port))
	if err != nil {
		fmt.Printf("Run server failed, err: %v\n", err)
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
