package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		err := process(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func process(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second) // time out with 3 seconds
	defer cancel()

	// randomly decide the processing time
	sec := rand.Intn(6)
	fmt.Printf("wait %d sec: ", sec)

	// pseudo process that takes <sec> seconds
	done := make(chan error, 1)
	go func(sec int) {
		time.Sleep(time.Duration(sec) * time.Second)
		done <- nil
	}(sec)

	// if context is done before receiving message from done channel, consider it as timeout.
	select {
	case <-done:
		fmt.Println("complete")
		return nil
	case <-ctx.Done():
		fmt.Println("timeout")
		return ctx.Err()
	}
}
