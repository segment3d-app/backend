package rabbitmq

import (
	"context"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMq struct {
	Channel *amqp091.Channel
}

func NewRabbitMq(source string) (*RabbitMq, error) {
	con, err := amqp091.Dial(source)
	if err != nil {
		return nil, err
	}

	rcon, err := con.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMq{
		Channel: rcon,
	}, nil
}

func (rmq *RabbitMq) PublishEvent(queue string, msg []byte) error {
	// create channel
	q, err := rmq.Channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	publishedMsg := amqp091.Publishing{
		ContentType:  "application/json",
		DeliveryMode: 2,
		Body:         msg,
	}

	err = rmq.Channel.PublishWithContext(ctx, "", q.Name, false, false, publishedMsg)
	if err != nil {
		return err
	}

	return nil
}
