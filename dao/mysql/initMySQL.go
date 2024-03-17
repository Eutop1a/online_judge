package mysql

import (
	"OnlineJudge/models"
	"OnlineJudge/setting"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(cfg *setting.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Protocal,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Error("connect DB failed, err: %v\n", zap.Error(err))
		return
	}

	// 创建表
	if DB.Debug().AutoMigrate(&models.User{}) != nil {
		fmt.Println("err in AutoMigrate(&models.User{}", err)
	}

	if DB.Debug().AutoMigrate(&models.Problems{}) != nil {
		fmt.Println("err in AutoMigrate(&models.Problems{}", err)
	}
	if DB.Debug().AutoMigrate(&models.Judgement{}) != nil {
		fmt.Println("err in AutoMigrate(&models.Judgement{}", err)
	}
	if DB.Debug().AutoMigrate(&models.ProgrammingLanguage{}) != nil {
		fmt.Println("err in AutoMigrate(&models.ProgrammingLanguage{}", err)
	}
	if DB.Debug().AutoMigrate(&models.Submission{}) != nil {
		fmt.Println("err in AutoMigrate(&models.Submission{}", err)
	}
	if DB.Debug().AutoMigrate(&models.TestCase{}) != nil {
		fmt.Println("err in AutoMigrate(&models.TestCase{}", err)
	}
	if DB.Debug().AutoMigrate(&models.SubmissionResult{}) != nil {
		fmt.Println("err in AutoMigrate(&models.SubmissionResult{}", err)
	}
	// 设置innodb事务行锁等待时间为10s，默认50s
	if err = DB.Exec("SET innodb_lock_wait_timeout = 10").Error; err != nil {
		fmt.Println("Failed to set innodb_lock_wait_timeout", err)
		return
	}

	return
}
