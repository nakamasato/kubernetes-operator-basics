package main

import (
	"fmt"
	"time"
)

func main() {
    waitCh := make(chan int)

    go process(waitCh)

    fmt.Println("waiting")

    <-waitCh // receive a message (wait until process is completed)

    fmt.Println("finished")
}

func process(ch chan int) {
    fmt.Println("process start")
    time.Sleep(1 * time.Second)

    ch <- 1 // send a message to the channel

    fmt.Println("process finished") // this might not be visible as main() finishes earlier
}
