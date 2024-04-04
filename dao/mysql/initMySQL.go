package mysql

import (
	"OnlineJudge/setting"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(cfg *setting.MySQLConfig) (err error) {
	// 创建数据库
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Protocal,
		cfg.Host,
		cfg.Port,
		"mysql",
	)
	createDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Error("connect DB with root failed, err: %v\n", zap.Error(err))
		return
	}

	checkDBSQL := `SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = 'OnlineJudge'`
	var dbName string
	if err = createDB.Raw(checkDBSQL).Scan(&dbName).Error; err != nil {
		zap.L().Error("scan DB failed, err: %v\n", zap.Error(err))
		return
	}
	// 如果数据库不存在，就创建
	if dbName == "" {
		createDBSQL := `CREATE DATABASE IF NOT EXISTS OnlineJudge`
		if err = createDB.Exec(createDBSQL).Error; err != nil {
			zap.L().Error("create DB with root failed, err: %v\n", zap.Error(err))
			return
		}
	}

	dsn = fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8mb4&parseTime=True",
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
	// 不存在就创建
	if !DB.Migrator().HasTable(&User{}) {
		if DB.Debug().AutoMigrate(&User{}) != nil {
			fmt.Println("err in AutoMigrate(&models.User{}", err)
		}
	}
	if !DB.Migrator().HasTable(&Admin{}) {
		if DB.Debug().AutoMigrate(&Admin{}) != nil {
			fmt.Println("err in AutoMigrate(&models.User{}", err)
		}
	}

	if !DB.Migrator().HasTable(&Problems{}) {
		if DB.Debug().AutoMigrate(&Problems{}) != nil {
			fmt.Println("err in AutoMigrate(&models.Problems{}", err)
		}
	}

	if !DB.Migrator().HasTable(&Judgement{}) {
		if DB.Debug().AutoMigrate(&Judgement{}) != nil {
			fmt.Println("err in AutoMigrate(&models.Judgement{}", err)
		}
	}

	if !DB.Migrator().HasTable(&ProgrammingLanguage{}) {
		if DB.Debug().AutoMigrate(&ProgrammingLanguage{}) != nil {
			fmt.Println("err in AutoMigrate(&models.ProgrammingLanguage{}", err)
		}
	}

	if !DB.Migrator().HasTable(&Submission{}) {
		if DB.Debug().AutoMigrate(&Submission{}) != nil {
			fmt.Println("err in AutoMigrate(&models.Submission{}", err)
		}
	}
	if !DB.Migrator().HasTable(&TestCase{}) {
		if DB.Debug().AutoMigrate(&TestCase{}) != nil {
			fmt.Println("err in AutoMigrate(&models.TestCase{}", err)
		}
	}

	if !DB.Migrator().HasTable(&SubmissionResult{}) {
		if DB.Debug().AutoMigrate(&SubmissionResult{}) != nil {
			fmt.Println("err in AutoMigrate(&models.SubmissionResult{}", err)
		}
	}
	// 设置innodb事务行锁等待时间为10s，默认50s
	if err = DB.Exec("SET innodb_lock_wait_timeout = 10").Error; err != nil {
		fmt.Println("Failed to set innodb_lock_wait_timeout", err)
		return
	}

	return
}
