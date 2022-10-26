# Channel

**Channels** are a typed conduit through which you can send and receive values with the channel operator, `<-`.

conduit: a channel for conveying water or other fluid.

### Example: send/receive a message to/from channel

```go
func main() {
	ch := make(chan int) // Create a channel with type

	go func() { ch <- 3 }() // Send a message to the channel
	go func() { ch <- 5 }() // Send a message to the channel

	fmt.Println(<-ch, <-ch) // Receive messages and print them
}
```

### Example: check if a channel is closed

```go
close(ch)
v, ok := <-ch
fmt.Println(v, ok)
```

ok is `false` if there are no more values to receive and the channel is closed.

### Example: wait until process completed

```go
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
```

### Example in Kubernetes

[WaitForCacheSync](https://github.com/kubernetes/client-go/blob/v0.25.0/tools/cache/shared_informer.go#L266-L287): Function to wait until cache is synced with Kubernetes cluster.

```go
// WaitForCacheSync waits for caches to populate.  It returns true if it was successful, false
// if the controller should shutdown
// callers should prefer WaitForNamedCacheSync()
func WaitForCacheSync(stopCh <-chan struct{}, cacheSyncs ...InformerSynced) bool {
    err := wait.PollImmediateUntil(syncedPollPeriod,
        func() (bool, error) {
            for _, syncFunc := range cacheSyncs {
                if !syncFunc() {
                    return false, nil
                }
            }
            return true, nil
        },
        stopCh)
    if err != nil {
        klog.V(2).Infof("stop requested")
        return false
    }

    klog.V(4).Infof("caches populated")
    return true
}
```

[wait.PollImmediateUntil](https://github.com/kubernetes/apimachinery/blob/v0.25.0/pkg/util/wait/wait.go#L299)

```go
// ContextForChannel derives a child context from a parent channel.
//
// The derived context's Done channel is closed when the returned cancel function
// is called or when the parent channel is closed, whichever happens first.
//
// Note the caller must *always* call the CancelFunc, otherwise resources may be leaked.
func ContextForChannel(parentCh <-chan struct{}) (context.Context, context.CancelFunc) {
    ctx, cancel := context.WithCancel(context.Background())

    go func() {
        select {
        case <-parentCh: // Call cancel when parentCh
            cancel()
        case <-ctx.Done():
        }
    }()
    return ctx, cancel
}
```

Call `cancel` when `parentCh` (originally passed from `WaitForCacheSync` as `stopCh`) receives a message. `cancel` function closes the Done channel of the context.
For more details about `context`, we'll study in the next section [context](../06-context/)

## Practice
1. [Channels](https://go.dev/tour/concurrency/2)
1. [Buffered Channels](https://go.dev/tour/concurrency/3)
1. [Range and Close](https://go.dev/tour/concurrency/4)
1. [Select](https://go.dev/tour/concurrency/5)
1. [Default Selection](https://go.dev/tour/concurrency/6)

## References
1. https://gobyexample.com/channels
