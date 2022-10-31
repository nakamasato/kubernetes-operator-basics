# Package and Module

A ***package*** is a collection of source files in the same directory that are compiled together.
A ***module*** is a collection of related Go packages that are released together.


## Steps

### 1. Create a Go module


```
go mod init <module_name>
```

1. If you want to use your module from other module. You can set `github.com/<github user/org name>/<github repo name>` to your module name.
1. This command will generate a go.mod file.

Example:

```
go mod init github.com/nakamasato/kubernetes-operator-basics/05-golang-basics/02-package-and-module
go: creating new go.mod: module github.com/nakamasato/kubernetes-operator-basics/05-golang-basics/02-package-and-module
```

I specified the path to the current directory.

`go.mod` file will be created.

```
module github.com/nakamasato/kubernetes-operator-basics/05-golang-basics/02-module

go 1.19
```

### 2. Hello World

1. Create a `main.go`.
    ```go
    package main

    import "fmt"

    func main() {
        fmt.Println("Hello, World")
    }
    ```

1. Run `main.go`
    ```
    go run main.go
    Hello, World
    ```

### 3. Use your own package

1. Add `mypackage`.
    Create a `mypackage/mypackage.go`
    ```go
    package mypackage

    func GetName() string {
        return "MyName"
    }
    ```
1. Use `mypackage` in the module.
    1. Import the pacakge in `main.go`
        ```go
        import (
            "fmt" // Packages in the standard library do not have a module path prefix
            "<module_name>/mypackage" // e.g. "github.com/nakamasato/kubernetes-operator-basics/05-golang-basics/02-module/mypackage"
        )
        ```
        - Standard library: https://pkg.go.dev/std
    1. Use `GetName` function in `mypackage` package.
        ```go
        func main() {
            name := mypackage.GetName()
            fmt.Printf("Hello, %s\n", name)
        }
        ```
1. Run `main.go`
    ```
    go run main.go
    Hello, MyName
    ```

### 4. Use an external package

1. Import a package
    ```go
    import "github.com/google/go-github/v47/github"
    ```

1. Use it.
    ```go
    client := github.NewClient(nil)
	// list all organizations for user "willnorris"
	orgs, _, err := client.Organizations.List(context.Background(), "willnorris", nil)
	if err != nil {
		fmt.Println("Error")
		os.Exit(1)
	}
	for i, org := range orgs {
		fmt.Println(i, *org.Login)
	}
    ```
1. Run `go mod tidy`
    `go.mod` is updated:
    ```
    require github.com/google/go-github/v47 v47.1.0

    require (
        github.com/google/go-querystring v1.1.0 // indirect
        golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
    )
    ```
1. Run main.go
    ```
    go run main.go
    ```

    ```
    Hello, MyName
    0 diso
    1 activitystreams
    2 indieweb
    3 webfinger
    4 todogroup
    5 maintainers
    6 perkeep
    7 tailscale
    ```

## Practice
1. Packages, variables, and functions.
    1. [Packages](https://go.dev/tour/basics/1)
    1. [Imports](https://go.dev/tour/basics/2)
    1. [Exported names](https://go.dev/tour/basics/3)
    1. [Functions](https://go.dev/tour/basics/4)
    1. [Functions continued](https://go.dev/tour/basics/5)n
    1. [Multiple results](https://go.dev/tour/basics/6)
    1. [Named return values](https://go.dev/tour/basics/7)
    1. [Variables](https://go.dev/tour/basics/8)
    1. [Variables with initializers](https://go.dev/tour/basics/9)
    1. [Short variable declarations](https://go.dev/tour/basics/10)
    1. [Basic types](https://go.dev/tour/basics/11)
    1. [Zero values](https://go.dev/tour/basics/12)
    1. [Type conversions](https://go.dev/tour/basics/13)
    1. [Type inference](https://go.dev/tour/basics/14)
    1. [Constants](https://go.dev/tour/basics/15)
    1. [Numeric Constants](https://go.dev/tour/basics/16)
1. Flow control statements: for, if, else, switch and defer
    1. [For](https://go.dev/tour/flowcontrol/1)
    1. [For continued](https://go.dev/tour/flowcontrol/2)
    1. [For is Go's "while"](https://go.dev/tour/flowcontrol/3)
    1. [Forever](https://go.dev/tour/flowcontrol/4)
    1. [If](https://go.dev/tour/flowcontrol/5)
    1. [If with a short statement](https://go.dev/tour/flowcontrol/6)
    1. [If and else](https://go.dev/tour/flowcontrol/7)
    1. [Exercise: Loops and Functions](https://go.dev/tour/flowcontrol/8)
    1. [Switch](https://go.dev/tour/flowcontrol/9)
    1. [Switch evaluation order](https://go.dev/tour/flowcontrol/10)
    1. [Switch with no condition](https://go.dev/tour/flowcontrol/11)
    1. [Defer](https://go.dev/tour/flowcontrol/12)
    1. [Stacking defers](https://go.dev/tour/flowcontrol/13)
