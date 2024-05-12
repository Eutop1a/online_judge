package mq

import (
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"online-judge/setting"
)

var RabbitMq *amqp.Connection

func InitRabbitMQ(cfg *setting.RabbitMQConfig) error {
	connString := fmt.Sprintf("%s://%s:%s@%s:%d/",
		cfg.RabbitMQ,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)
	fmt.Println(connString)
	conn, err := amqp.Dial(connString)
	if err != nil {
		zap.L().Error("mq Dial", zap.Error(err))
		return err
	}
	RabbitMq = conn
	return nil
}
