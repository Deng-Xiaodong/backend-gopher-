package rabbitmq

import (
	"fmt"
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
	//queue
	qa, _ := c.QueueDeclare("a", false, false, false, false, nil)
	qb, _ := c.QueueDeclare("b", false, false, false, false, nil)
	//bind
	_ = c.QueueBind(qa.Name, "", "fanout_logs", false, nil)
	_ = c.QueueBind(qb.Name, "", "fanout_logs", false, nil)
	//publish
	for i := 0; i < 5; i++ {
		_ = c.Publish("fanout_logs", "", false, false,
			amqp.Publishing{ContentType: "text/plain", Body: []byte(fmt.Sprintf("logs_%d", i))})
	}

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
	time.Sleep(5 * time.Second)
}
