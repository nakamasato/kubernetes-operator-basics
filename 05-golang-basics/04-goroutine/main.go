package main

import (
	"fmt"
	"time"
)

func main() {
	go process()
	fmt.Println("next")
	time.Sleep(2 * time.Second)
}

func process() {
	time.Sleep(1 * time.Second) // heavy process
	fmt.Println("goroutine completed")
}
