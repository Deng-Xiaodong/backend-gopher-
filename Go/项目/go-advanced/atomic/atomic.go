package main

import (
	"fmt"
	"sync/atomic"
)

type Config struct {
	state int
}

func main() {

	var config atomic.Value
	config.Store(&Config{
		state: 10,
	})

	val := config.Load()
	fmt.Printf("val:%T,%v", val, val)
}
