package main

import (
	"fmt"
	"sync"
	"time"
)

var done = false

func write(content string, c *sync.Cond) {
	println("write start")
	c.L.Lock()
	fmt.Printf("write:%s\n", content)
	done = true
	c.L.Unlock()
	println("wake all")
	c.Broadcast()
}
func read(content string, c *sync.Cond) {
	c.L.Lock()
	for !done {
		c.Wait()
	}
	fmt.Printf("read:%s\n", content)
	c.L.Unlock()
}
func main() {
	cd := sync.NewCond(&sync.Mutex{})
	go write("hello", cd)
	go read("hi i am one", cd)
	go read("hi i am two", cd)
	go read("hi i am three", cd)
	time.Sleep(1 * time.Second)
}
