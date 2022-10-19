package main

import (
	"fmt"
)

func main() {
	ch := make(chan int) // Create a channel with type

	go func() { ch <- 3 }() // Send a message to the channel
	go func() { ch <- 5 }() // Send a message to the channel

	fmt.Println(<-ch, <-ch) // Receive messages and print them
}
