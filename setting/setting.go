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
	Host           string `mapstructure:"host"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	Protocol       string `mapstructure:"protocol"`
	DbName         string `mapstructure:"dbname"`
	Port           int    `mapstructure:"port"`
	MaxOpenConns   int    `mapstructure:"max_open_conn"`
	MaxIdelConns   int    `mapstructure:"max_idle_conn"`
	DeleteInterval int    `mapstructure:"delete_interval"`
}

type RedisConfig struct {
	Host                string `mapstructure:"host"`
	Password            string `mapstructure:"password"`
	Port                int    `mapstructure:"port"`
	DB                  int    `mapstructure:"db"`
	PoolSize            int    `mapstructure:"pool_size"`
	PersistenceInterval int    `mapstructure:"persistence_interval"`
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
	// 读取配置文件
	viper.SetConfigFile("./conf/config.yaml")
	// 读取环境变量
	viper.WatchConfig()
	// 监听配置文件变化
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config file changed:", in.Name)
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal() failed, err: %v\n", err)
		}
	})
	// 查找并读取配置文件
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("viper.ReadInConfig() failed, err: %v", err))
	}
	// 把读取到的配置信息反序列化到Conf变量中
	if err = viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Errorf("viper.Unmarshal() failed, err: %v\n", err))
	}

	return
}
