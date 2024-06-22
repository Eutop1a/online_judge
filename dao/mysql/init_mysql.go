package mysql

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"online_judge/setting"
	"time"
)

var DB *gorm.DB

func Init(cfg *setting.MySQLConfig) (err error) {
	// 创建 onlinejudge 数据库
	if err = CreateDatabase(cfg); err != nil {
		return err
	}
	// 创建 onlinejudge 数据库中所有需要使用的表
	if err = CreateTables(); err != nil {
		return err
	}
	// 开启后台删除功能，删除 deleted 的记录
	go BgDeleteMysql()
	return
}

func CreateDatabase(cfg *setting.MySQLConfig) (err error) {
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

	//DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})

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
	models := []interface{}{
		&User{},
		//&Admin{},
		&Problems{},
		&TestCase{},
		&ProblemWithFile{},
		&TestCaseWithFile{},
		&Category{},
		&ProblemCategory{},
		&Submission{},
		&Judgement{},
	}

	if err = DB.AutoMigrate(models...); err != nil {
		zap.L().Error("mysql-CreateTables-AutoMigrate", zap.Error(err))
		return
	}

	// 检查是否存在索引
	var count int64
	DB.Raw("SELECT COUNT(1) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'problems' AND index_name = 'idx_title'").Scan(&count)
	if count == 0 {
		// 添加FULLTEXT索引
		err := DB.Exec("ALTER TABLE problems ADD FULLTEXT INDEX idx_title (title)").Error
		if err != nil {
			zap.L().Error("mysql-CreateTables-create FULLTEXT error ", zap.Error(err))
			return err
		}
	}

	// 设置innodb事务行锁等待时间为10s，默认50s
	if err = DB.Exec("SET innodb_lock_wait_timeout = 10").Error; err != nil {
		zap.L().Error("mysql-CreateTables-Exec-SET-innodb_lock_wait_timeout ", zap.Error(err))
		return
	}
	return
}
