# Goroutine

A goroutine is a lightweight thread managed by the Go runtime.

```go
go f(x)
```

## Example: with anonymous function

```go
func main() {
	go func() {
		time.Sleep(1 * time.Second) // heavy process
		fmt.Println("goroutine completed")
	}()
	fmt.Println("next")
	time.Sleep(2 * time.Second)
}
```

## Example: named function

```go
func main() {
	go process()
	fmt.Println("next")
	time.Sleep(2 * time.Second)
}

func process() {
	time.Sleep(1 * time.Second) // heavy process
	fmt.Println("goroutine completed")
}
```



## Practice
1. [Goroutines](https://go.dev/tour/concurrency/1)
