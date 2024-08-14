package amqp

import (
	"fmt"
	"log"

	"golang.org/x/net/context"

	. "github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/common"

	rabbitmq "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	amqpURI      string
	exchange     string
	exchangeType string
	routingKey   string
	reliable     bool
}

type Consumer struct {
	amqpURI      string
	exchange     string
	exchangeType string
	queue        string
	routingKey   string
	consumerTag  string
}

func NewPublisher(cfg AmqpConfig) *Publisher {
	return &Publisher{
		amqpURI:      cfg.Uri,
		exchange:     cfg.ExchangeName,
		exchangeType: cfg.ExchangeType,
		routingKey:   cfg.RoutingKey,
		reliable:     cfg.Reliable,
	}
}

func NewConsumer(cfg AmqpConfig) *Consumer {
	return &Consumer{
		amqpURI:      cfg.Uri,
		exchange:     cfg.ExchangeName,
		exchangeType: cfg.ExchangeType,
		queue:        cfg.QueueName,
		routingKey:   cfg.RoutingKey,
		consumerTag:  cfg.ConsumerTag,
	}
}

func (r *Publisher) Publish(message string) error {
	log.Printf("dialing %q", r.amqpURI)
	connection, err := rabbitmq.Dial(r.amqpURI)
	if err != nil {
		return fmt.Errorf("dial: %s", err)
	}
	defer func(connection *rabbitmq.Connection) {
		err := connection.Close()
		if err != nil {
			log.Printf("Error closing connection: %s", err)
		}
	}(connection)

	log.Printf("got Connection, getting Channel")
	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("channel: %s", err)
	}

	log.Printf("got Channel, declaring %q Exchange (%q)", r.exchangeType, r.exchange)
	if err := channel.ExchangeDeclare(
		r.exchange,
		r.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("exchange Declare: %s", err)
	}

	if r.reliable {
		log.Printf("enabling publishing confirms.")
		if err := channel.Confirm(false); err != nil {
			return fmt.Errorf("channel could not be put into confirm mode: %s", err)
		}

		confirms := channel.NotifyPublish(make(chan rabbitmq.Confirmation, 1))

		defer confirmOne(confirms)
	}

	log.Printf("declared Exchange, publishing %dB body (%v)", len(message), message)
	if err = channel.PublishWithContext(
		context.Background(),
		r.exchange,
		r.routingKey,
		false,
		false,
		rabbitmq.Publishing{
			Headers:      rabbitmq.Table{},
			ContentType:  "text/plain",
			Body:         []byte(message),
			DeliveryMode: rabbitmq.Transient,
		},
	); err != nil {
		return fmt.Errorf("exchange Publish: %w", err)
	}

	return nil
}

func (c *Consumer) Consume() error {
	log.Printf("dialing %q", c.amqpURI)
	conn, err := rabbitmq.Dial(c.amqpURI)
	if err != nil {
		return fmt.Errorf("dial: %s", err)
	}

	go func() {
		fmt.Printf("closing: %s", <-conn.NotifyClose(make(chan *rabbitmq.Error)))
	}()

	log.Printf("got Connection, getting Channel")
	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %s", err)
	}

	log.Printf("got Channel, declaring Exchange (%q)", c.exchange)
	if err = channel.ExchangeDeclare(
		c.exchange,
		c.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("exchange Declare: %s", err)
	}

	log.Printf("declared Exchange, declaring Queue %q", c.queue)
	queue, err := channel.QueueDeclare(
		c.queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue Declare: %s", err)
	}

	log.Printf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, c.routingKey)

	if err = channel.QueueBind(
		queue.Name,
		c.routingKey,
		c.exchange,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("queue Bind: %s", err)
	}

	log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", c.consumerTag)
	deliveries, err := channel.Consume(
		queue.Name,
		c.consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue Consume: %s", err)
	}

	for d := range deliveries {
		log.Printf("got %dB delivery: [%v] %v", len(d.Body), d.DeliveryTag, string(d.Body))
		_ = d.Ack(false)
	}

	return nil
}

func confirmOne(confirms <-chan rabbitmq.Confirmation) {
	log.Printf("waiting for confirmation of one publishing")

	if confirmed := <-confirms; confirmed.Ack {
		log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}
