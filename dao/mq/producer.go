package mq

import (
	"go.uber.org/zap"
	"online-judge/consts/mq_name"

	"github.com/streadway/amqp"
)

func SendMessage2MQ(body []byte) (err error) {
	ch, err := RabbitMq.Channel()
	if err != nil {
		return
	}
	q, _ := ch.QueueDeclare(mq_name.RabbitMQProblemQueueName, true, false, false, false, nil)
	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})
	if err != nil {
		return
	}
	zap.L().Debug("mq-producer-Publish send msg to MQ successfully")
	return
}
