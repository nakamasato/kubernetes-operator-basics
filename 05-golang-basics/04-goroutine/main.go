package main

import (
	"fmt"
	"time"
)

func main() {
	go func() {
		time.Sleep(1 * time.Second) // heavy process
		fmt.Println("goroutine completed")
	}()
	fmt.Println("next")
	time.Sleep(2 * time.Second)
}

func process() {
	time.Sleep(1 * time.Second) // heavy process
	fmt.Println("goroutine completed")
}
