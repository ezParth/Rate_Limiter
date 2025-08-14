package main

import (
	"context"
	"fmt"
	"time"
)

var ctx, cancel = context.WithCancel(context.Background())

func worker(id int) {
	if id == 5 {
		fmt.Println("Cancelled!")
		cancel()
	}

	fmt.Println("id: ", id)
}

func main() {
	defer cancel()
	for id := 1; id < 15; id++ {
		go worker(id)
	}

	select {
	case <-ctx.Done():
		fmt.Println("ctx done: ", ctx.Err())
	case <-time.After(time.Second * 15):
		fmt.Println("All work done!")
	}
}
