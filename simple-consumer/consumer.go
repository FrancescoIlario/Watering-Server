// This example declares a durable Exchange, an ephemeral (auto-delete) Queue,
// binds the Queue to the Exchange with a binding key, and consumes every
// message published to that Exchange with that routing key.
//
package simple_consumer

import (
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

var (
	uriC          = flag.String("uriC", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	exchangeC     = flag.String("exchangeC", "watering-exchange", "Durable, non-auto-deleted AMQP exchange name")
	exchangeTypeC = flag.String("exchange-typeC", "direct", "Exchange type - direct|fanout|topic|x-custom")
	queueC        = flag.String("queueC", "watering-queue", "Ephemeral AMQP queue name")
	bindingKeyC   = flag.String("keyC", "watering-key", "AMQP binding key")
	consumerTagC  = flag.String("consumer-tagC", "watering-consumer", "AMQP consumer tag (should not be blank)")
)

func init() {
	flag.Parse()
}

func StartConsumer(lifetime time.Duration) {
	c, err := NewConsumer(*uriC, *exchangeC, *exchangeTypeC, *queueC, *bindingKeyC, *consumerTagC)
	if err != nil {
		log.Fatalf("%s", err)
	}

	if lifetime > 0 {
		log.Printf("running for %s", lifetime)
		time.Sleep(lifetime)
	} else {

		log.Printf("running forever")
		select {}
	}

	log.Printf("shutting down")

	if err := c.Shutdown(); err != nil {
		log.Fatalf("error during shutdown: %s", err)
	}
}

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

func NewConsumer(amqpURI, exchange, exchangeType, queueName, key, ctag string) (*Consumer, error) {
	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
		done:    make(chan error),
	}

	var err error

	log.Printf("dialing %q", amqpURI)
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("dial: %s", err)
	}

	go func() {
		fmt.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	log.Printf("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("channel: %s", err)
	}

	log.Printf("got Channel, declaring Exchange (%q)", exchange)
	if err = c.channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("exchange Declare: %s", err)
	}

	log.Printf("declared Exchange, declaring Queue %q", queueName)
	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("queue declare: %s", err)
	}

	log.Printf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, key)

	if err = c.channel.QueueBind(
		queue.Name, // name of the queue
		key,        // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return nil, fmt.Errorf("queue bind: %s", err)
	}

	log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)
	deliveries, err := c.channel.Consume(
		queue.Name, // name
		c.tag,      // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("queue consume: %s", err)
	}

	go handle(deliveries, c.done)

	return c, nil
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		processDelivery(&d)
		_ = d.Ack(false)
	}

	log.Printf("handle: deliveries channel closed")
	done <- nil
}

func processDelivery(delivery *amqp.Delivery) {
	logDelivery(delivery)
}

func logDelivery(delivery *amqp.Delivery) {
	log.Printf(
		"got %dB delivery: [%v] %q",
		len(delivery.Body),
		delivery.DeliveryTag,
		delivery.Body,
	)
}