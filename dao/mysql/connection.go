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

	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Problems{}) //
	DB.AutoMigrate(&models.Submission{})
	DB.AutoMigrate(&models.Judgement{})
	DB.AutoMigrate(&models.ProgrammingLanguage{})
	DB.AutoMigrate(&models.TestCase{}) //
	DB.AutoMigrate(&models.SubmissionResult{})

	return
}
