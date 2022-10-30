# Hello World

## Go version

```
go version
```

If you use [gvm](https://github.com/moovweb/gvm) to install Go, you can change Go versions easily.

1. Check the available versions

    ```
    gvm list
    ```
1. Swtich version
    ```
    gvm use go1.18
    ```
1. Install new version
    ```
    gvm install go1.17
    ```
1. Uninstall a version
    ```
    gvm uninstall go1.17
    ```

## Run `main.go`

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World")
}
```

1. Run with `go run`

    ```
    go run main.go
    Hello, World
    ```

1. Run with compiled file:

    ```
    go build main.go
    ```
    ```
    ./main
    Hello, World
    ```

## Function

```go
func greet(language, name string) (string, error) {
	if language == "Spanish" {
		return fmt.Sprintf("Ola, %s", name), nil
	}
	return fmt.Sprintf("Hello, %s", name), nil
}
```

```go
func main() {
	fmt.Println("Hello, World")
	greetText := greet("Spanish", "Naka")
	fmt.Println(greetText)
}
```

## Other frequently used commands

1. Format:
    ```
    go fmt main.go
    ```
1. Validate: report likely mistakes in packages
    ```
    go vet main.go
    ```
1. Install
    ```
    go install <package>
    ```
1. Import necessary libraries and remove unnecessary libraries.
    ```
    go mod tidy
    ```
For more details, https://pkg.go.dev/cmd/go

## Practices
1. [Methods](https://go.dev/tour/methods/1)
1. [Methods are functions](https://go.dev/tour/methods/2)
1. [Methods continued](https://go.dev/tour/methods/3)
1. [Pointer receivers](https://go.dev/tour/methods/4)
1. [Pointers and functions](https://go.dev/tour/methods/5)
1. [Methods and pointer indirection](https://go.dev/tour/methods/6)
1. [Methods and pointer indirection (2)](https://go.dev/tour/methods/7)
1. [Choosing a value or pointer receiver](https://go.dev/tour/methods/8)
