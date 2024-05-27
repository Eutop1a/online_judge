package mysql

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"online-judge/setting"
	"time"
)

var DB *gorm.DB

func Init(cfg *setting.MySQLConfig) (err error) {
	// 创建 onlinejudge 数据库
	CreateDatabase(cfg)
	// 创建 onlinejudge 数据库中所有需要使用的表
	err = CreateTables()
	// 开启后台删除功能，删除 deleted 的记录
	go BgDeleteMysql()
	return
}

func CreateDatabase(cfg *setting.MySQLConfig) {
	// 创建数据库
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Protocol,
		cfg.Host,
		cfg.Port,
		"mysql",
	)
	createDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Error("mysql-CreateDatabase-Open ", zap.Error(err))
		return
	}

	checkDBSQL := `SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = 'onlinejudge'`
	var dbName string
	if err = createDB.Raw(checkDBSQL).Scan(&dbName).Error; err != nil {
		zap.L().Error("mysql-CreateDatabase-Scan ", zap.Error(err))
		return
	}
	// 如果数据库不存在，就创建
	if dbName == "" {
		createDBSQL := `CREATE DATABASE IF NOT EXISTS onlinejudge`
		if err = createDB.Exec(createDBSQL).Error; err != nil {
			zap.L().Error("mysql-CreateDatabase-Exec create DB with root failed ", zap.Error(err))
			return
		}
	}

	dsn = fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Protocol,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Error("mysql-CreateDatabase-Open connect DB failed ", zap.Error(err))
		return
	}
	// 设置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		zap.L().Error("mysql-CreateDatabase-DB ", zap.Error(err))
		return
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdelConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return
}

func CreateTables() (err error) {
	// 创建表
	// 不存在就创建
	if !DB.Migrator().HasTable(&User{}) {
		if DB.Debug().AutoMigrate(&User{}) != nil {
			zap.L().Error("mysql-CreateTables-AutoMigrate-User ", zap.Error(err))
			return
		}
	}

	if !DB.Migrator().HasTable(&Admin{}) {
		if DB.Debug().AutoMigrate(&Admin{}) != nil {
			zap.L().Error("mysql-CreateTables-AutoMigrate-Admin ", zap.Error(err))
			return
		}
	}

	if !DB.Migrator().HasTable(&Problems{}) {
		if DB.Debug().AutoMigrate(&Problems{}) != nil {
			zap.L().Error("mysql-CreateTables-AutoMigrate-Problems ", zap.Error(err))
			return
		}
	}

	if !DB.Migrator().HasTable(&TestCase{}) {
		if DB.Debug().AutoMigrate(&TestCase{}) != nil {
			zap.L().Error("mysql-CreateTables-AutoMigrate-TestCase ", zap.Error(err))
			return
		}
	}

	if !DB.Migrator().HasTable(&ProblemWithFile{}) {
		if DB.Debug().AutoMigrate(&ProblemWithFile{}) != nil {
			zap.L().Error("mysql-CreateTables-AutoMigrate-ProblemWithFile ", zap.Error(err))
			return
		}
	}
	if !DB.Migrator().HasTable(&TestCaseWithFile{}) {
		if DB.Debug().AutoMigrate(&TestCaseWithFile{}) != nil {
			zap.L().Error("mysql-CreateTables-AutoMigrate-ProblemWithFile ", zap.Error(err))
			return
		}
	}
	//if !DB.Migrator().HasTable(&ProgrammingLanguage{}) {
	//	if DB.Debug().AutoMigrate(&ProgrammingLanguage{}) != nil {
	//		zap.L().Error("mysql-CreateTables-AutoMigrate-ProgrammingLanguage ", zap.Error(err))
	//      return
	//	}
	//}

	if !DB.Migrator().HasTable(&Submission{}) {
		if DB.Debug().AutoMigrate(&Submission{}) != nil {
			zap.L().Error("mysql-CreateTables-AutoMigrate-Submission ", zap.Error(err))
			return
		}
	}

	//if !DB.Migrator().HasTable(&SubmissionResult{}) {
	//	if DB.Debug().AutoMigrate(&SubmissionResult{}) != nil {
	//		zap.L().Error("mysql-CreateTables-AutoMigrate-SubmissionResult ", zap.Error(err))
	//		return
	//	}
	//}

	if !DB.Migrator().HasTable(&Judgement{}) {
		if DB.Debug().AutoMigrate(&Judgement{}) != nil {
			zap.L().Error("mysql-CreateTables-AutoMigrate-Judgement ", zap.Error(err))
			return
		}
	}
	// 设置innodb事务行锁等待时间为10s，默认50s
	if err = DB.Exec("SET innodb_lock_wait_timeout = 10").Error; err != nil {
		zap.L().Error("mysql-CreateTables-Exec-SET-innodb_lock_wait_timeout ", zap.Error(err))
		return
	}
	return
}
