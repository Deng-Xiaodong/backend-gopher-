package rabbitmq

import (
	"sync"
	"testing"
)

func TestDirect(t *testing.T) {
	DirectPublish()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		DirectConsumer()
	}()

	wg.Wait()
}

func TestFanout(t *testing.T) {
	FanoutPublish()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		FanoutConsumer()
	}()
	wg.Wait()
}

func TestTopic(t *testing.T) {
	TopicPublish()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		TopicConsumer()
	}()
	wg.Wait()
}
