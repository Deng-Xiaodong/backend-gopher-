package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

func DirectPublish() {
	conn, err1 := amqp.Dial("amqp://guest:guest@miaosha.peadx.live:5672/miaosha")
	if err1 != nil {
		log.Fatalln(err1)
	}
	c, err2 := conn.Channel()
	if err2 != nil {
		log.Fatalln(err2)
	}
	cc, err3 := conn.Channel()
	if err3 != nil {
		log.Fatalln(err3)
	}

	//exchange
	//if err := c.ExchangeDeclare(
	//	"logs",
	//	"direct",
	//	false,
	//	false,
	//	false,
	//	false,
	//	nil,
	//); err != nil {
	//	log.Fatalln(err)
	//}
	//if err := c.ExchangeDeclare(
	//	"data",
	//	"direct",
	//	false,
	//	false,
	//	false,
	//	false,
	//	nil,
	//); err != nil {
	//	log.Fatalln(err)
	//}
	//queue
	c.QueueDeclare("a", false, false, false, false, nil)
	//c.QueueDeclare("b", false, false, false, false, nil)
	//c.QueueDeclare("c", false, false, false, false, nil)
	//c.QueueDeclare("g", false, false, false, false, nil)
	//c.QueueDelete("d", false, false, false)

	//bind
	_ = c.QueueBind("a", "key1", "logs", false, nil)
	//_ = c.QueueBind("b", "key2", "logs", false, nil)
	//_ = c.QueueBind("c", "key3", "logs", false, nil)
	//_ = c.QueueBind("d", "key4", "logs", false, nil)
	//_ = c.QueueBind("r", "key5", "data", false, nil)
	//_ = c.QueueBind("g", "key6", "data", false, nil)
	if err := c.Publish(
		"logs",
		"key1",
		false,
		false,
		amqp.Publishing{ContentType: "text/plain", Body: []byte("exchange logs with routing key key1")},
	); err != nil {
		log.Fatalln(err)
	}

	if err := c.Publish(
		"logs",
		"key2",
		false,
		false,
		amqp.Publishing{ContentType: "text/plain", Body: []byte("exchange logs with routing key key2")},
	); err != nil {
		log.Fatalln(err)
	}
	if err := c.Publish(
		"logs",
		"key3",
		false,
		false,
		amqp.Publishing{ContentType: "text/plain", Body: []byte("exchange logs with routing key key3")},
	); err != nil {
		log.Fatalln(err)
	}
	if err := cc.Publish(
		"data",
		"key5",
		false,
		false,
		amqp.Publishing{ContentType: "text/plain", Body: []byte("exchange data with routing key key5")},
	); err != nil {
		log.Fatalln(err)
	}
	if err := cc.Publish(
		"data",
		"key6",
		false,
		false,
		amqp.Publishing{ContentType: "text/plain", Body: []byte("exchange data with routing key key6")},
	); err != nil {
		log.Fatalln(err)
	}

}

func DirectConsumer() {
	conn, err1 := amqp.Dial("amqp://guest:guest@miaosha.peadx.live:5672/miaosha")
	if err1 != nil {
		log.Fatalln(err1)
	}
	c, err2 := conn.Channel()
	if err2 != nil {
		log.Fatalln("channel没打开", err2)
	}

	go consumer(c, "a")
	go consumer(c, "b")
	go consumer(c, "c")
	go consumer(c, "r")
	go consumer(c, "g")
	time.Sleep(5 * time.Second)

}
func consumer(c *amqp.Channel, name string) {
	csm, err := c.Consume(name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalln(err)
	}
	for msg := range csm {
		log.Printf("队列%s 接收到信息：%s\n", name, msg.Body)
		msg.Ack(false)
	}
}
