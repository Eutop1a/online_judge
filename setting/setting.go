package setting

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(AppConfig)

type AppConfig struct {
	Name            string `mapstructure:"name"`
	Mode            string `mapstructure:"mode"`
	Version         string `mapstructure:"version"`
	StartTime       string `mapstructure:"start_time"`
	MachineID       int64  `mapstructure:"machine_id"`
	Port            int    `mapstructure:"port"`
	*LogConfig      `mapstructure:"log"`
	*MySQLConfig    `mapstructure:"mysql"`
	*RedisConfig    `mapstructure:"redis"`
	*RabbitMQConfig `mapstructure:"rabbitmq"`
	*EtcdConfig     `mapstructure:"etcd"`
}

type LogConfig struct {
	Level     string `mapstructure:"level"`
	Filename  string `mapstructure:"filename"`
	MaxSize   int    `mapstructure:"max_size"`
	MaxAge    int    `mapstructure:"max_age"`
	MaxBackup int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Protocol     string `mapstructure:"protocol"`
	DbName       string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conn"`
	MaxIdelConns int    `mapstructure:"max_idle_conn"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type RabbitMQConfig struct {
	RabbitMQ string `mapstructure:"rabbit_mq"`
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
}

type EtcdConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func Init() (err error) {
	viper.SetConfigFile("./conf/config_docker.yaml")

	err = viper.ReadInConfig() // 读取配置信息
	if err != nil {
		fmt.Printf("viper.ReadInConfig() failed, err: %v\n", err)
		return
	}
	// 将配置文件读取到Conf变量中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal() failed, err: %v\n", err)
	}

	// 检测到配置文件变动就读取
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal() failed, err: %v\n", err)
		}
	})
	return
}
