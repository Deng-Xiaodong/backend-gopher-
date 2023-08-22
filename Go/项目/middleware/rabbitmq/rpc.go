package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

func toUpper(s string) string {
	return strings.ToUpper(s)
}
func startServer() {
	url := "amqp://guest:guest@miaosha.peadx.live:5672/miaosha"
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalln(err)
	}
	ch, _ := conn.Channel()

	q, _ := ch.QueueDeclare("rpc", false, false, false, false, nil)

	msgs, _ := ch.Consume(q.Name, "", true, false, false, false, nil)

	for msg := range msgs {
		log.Printf("接收到请求号为%s的请求\n", msg.CorrelationId)
		_ = ch.Publish("", msg.ReplyTo, false, false, amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: msg.CorrelationId,
			Body:          []byte(toUpper(string(msg.Body))),
		})
	}
}
func startClient() {
	url := "amqp://guest:guest@miaosha.peadx.live:5672/miaosha"
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalln(err)
	}
	ch, _ := conn.Channel()
	rq, _ := ch.QueueDeclare("", false, false, true, false, nil)

	q, _ := ch.QueueDeclare("rpc", false, false, false, false, nil)
	//发送请求
	for i := 0; i < 5; i++ {
		ch.Publish("", q.Name, false, false, amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: strconv.Itoa(i),
			Body:          []byte(randomString(10)),
			ReplyTo:       rq.Name,
		})
	}
	//接收响应
	msgs, _ := ch.Consume(rq.Name, "", true, false, false, false, nil)
	for msg := range msgs {
		fmt.Printf("接收到序列号为 %s 的响应数据：%s\n", msg.CorrelationId, msg.Body)
	}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}
