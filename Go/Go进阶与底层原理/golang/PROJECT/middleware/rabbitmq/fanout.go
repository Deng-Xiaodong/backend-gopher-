package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

func FanoutPublish() {
	conn, err1 := amqp.Dial("amqp://guest:guest@miaosha.peadx.live:5672/miaosha")
	if err1 != nil {
		log.Fatalln(err1)
	}
	c, err2 := conn.Channel()
	if err2 != nil {
		log.Fatalln(err2)
	}
	//exchange
	_ = c.ExchangeDeclare("fanout_logs", "fanout", false, false, false, false, nil)
	//queue 用已有的
	//bind
	_ = c.QueueBind("a", "", "fanout_logs", false, nil)
	_ = c.QueueBind("c", "", "fanout_logs", false, nil)
	_ = c.QueueBind("g", "", "fanout_logs", false, nil)

	//publish
	_ = c.Publish("fanout_logs", "", false, false, amqp.Publishing{ContentType: "text/plain", Body: []byte("exchange fanout_logs with key nil")})
}

func FanoutConsumer() {
	conn, err1 := amqp.Dial("amqp://guest:guest@miaosha.peadx.live:5672/miaosha")
	if err1 != nil {
		log.Fatalln(err1)
	}
	c, err2 := conn.Channel()
	if err2 != nil {
		log.Fatalln(err2)
	}
	go consumer(c, "a")
	go consumer(c, "b")
	go consumer(c, "c")
	go consumer(c, "r")
	go consumer(c, "g")
	time.Sleep(5 * time.Second)
}
