package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

func TopicPublish() {
	conn, err1 := amqp.Dial("amqp://guest:guest@miaosha.peadx.live:5672/miaosha")
	if err1 != nil {
		log.Fatalln(err1)
	}
	c, err2 := conn.Channel()
	if err2 != nil {
		log.Fatalln(err2)
	}
	//exchange
	_ = c.ExchangeDeclare("topic_logs", "topic", false, false, false, false, nil)
	//bind
	//_ = c.QueueBind("b", "a.*", "topic_logs", false, nil)
	//_ = c.QueueBind("d", "b.*", "topic_logs", false, nil)
	//_ = c.QueueBind("r", "*", "topic_logs", false, nil)
	//publish
	_ = c.Publish("topic_logs", "a.d", false, false, amqp.Publishing{ContentType: "text/plain", Body: []byte("exchange topic_logs with key a.d")})
	_ = c.Publish("topic_logs", "d.a", false, false, amqp.Publishing{ContentType: "text/plain", Body: []byte("exchange topic_logs with key d.a")})
	_ = c.Publish("topic_logs", "a", false, false, amqp.Publishing{ContentType: "text/plain", Body: []byte("exchange topic_logs with key a")})
}

func TopicConsumer() {
	conn, err1 := amqp.Dial("amqp://guest:guest@miaosha.peadx.live:5672/miaosha")
	if err1 != nil {
		log.Fatalln(err1)
	}
	c, err2 := conn.Channel()
	if err2 != nil {
		log.Fatalln(err2)
	}

	go consumer(c, "b")
	go consumer(c, "r")
	go consumer(c, "d")
	time.Sleep(5 * time.Second)
}
