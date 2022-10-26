# Interface and Struct

## 1. Struct

1. Define a struct and its method

    ```go
    type Parent struct {
        Id int
        Name string
    }

    func (p Parent) GetName() string {
        return p.Name
    }

    func (p *Parent) SetName(name string) {
        p.Name = name
    }
    ```

    1. `p` for `GetName` is a receiver.
    1. `p` for `SetName` is a pointer receiver, which is used to modify value of the receiver.

    Exercise: [Pointer Receiver](https://go.dev/tour/methods/4)

1. Use the new struct.

    ```go
    func main() {
        p := Parent{Id: 1, Name: "John"}
        fmt.Println(p)
        p.SetName("Bob")
        fmt.Println(p.GetName())
    }
    ```
1. Run.
    ```
    go run main.go
    {1 John}
    Bob
    ```

## 2. Embedding

1. Embed `Parent` struct in `Child` struct.

    ```go
    type Child struct {
        Parent // Embed Parent struct
        OtherField string
    }
    ```

1. `Child` can use the `Parent`'s methods.

    ```go
    func main() {
        p := Parent{Id: 1, Name: "John"}
        ch := Child{Parent: p}
        fmt.Printf("id: %d, name: %s\n", ch.Id, ch.Name)

        ch = Child{Parent{Id: 1, Name: "child"}, "Other"}
        fmt.Println(ch.GetName())
    }
    ```
1. Run.
    ```
    go run main.go
    id: 1, name: John
    child
    ```

### Example

**Embedding** is often used in Kubernetes operator too!

Example: [memcached_controller.go](https://github.com/operator-framework/operator-sdk/blob/de6a14d03de3c36dcc9de3891af788b49d15f0f3/testdata/go/v3/memcached-operator/controllers/memcached_controller.go#L57-L61)

```go
type MemcachedReconciler struct {
	client.Client // Embeded
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}
```

In this example, `MemcachedReconciler` has all the methods of [client.Client](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/client/interfaces.go) such as `Get`, `Create`, `Update`, etc.


## Interface

1. Define an interface `Product`

    ```go
    type Object interface {
        GetName() string
    }
    ```

1. Use it.
    ```go
    func main() {
        p := Parent{Id: 1, Name: "John"}
        printObjectName(&p)
        ch := Child{Parent: p}
        printObjectName(&ch)

        ch = Child{Parent{Id: 1, Name: "child"}, "Other"}
        printObjectName(&ch)
    }

    func printObjectName(obj Object) {
        fmt.Println(obj.GetName())
    }
    ```
1. Run.
    ```
    go run main.go
    John
    John
    child
    ```

## References
1. [Structs and Interfaces](https://www.golang-book.com/books/intro/9)
1. [Method Pointer Receivers in Interfaces](https://sentry.io/answers/interface-pointer-receiver/)


## Practices
1. [Structs](https://go.dev/tour/moretypes/2)
1. [Struct Fields](https://go.dev/tour/moretypes/3)
1. [Pointer Receiver](https://go.dev/tour/methods/4)
1. [Interfaces](https://go.dev/tour/methods/9)
1. [Interfaces are implemented implicitly](https://go.dev/tour/methods/10)
1. [Interface values](https://go.dev/tour/methods/11)
1. [Interface values with nil underlying values](https://go.dev/tour/methods/12)
1. [Nil interface values](https://go.dev/tour/methods/13)
1. [The empty interface](https://go.dev/tour/methods/14)
1. [Type assertions](https://go.dev/tour/methods/15)
1. [Type switches](https://go.dev/tour/methods/16)
1. [Stringers](https://go.dev/tour/methods/17)
1. [Exercise: Stringers](https://go.dev/tour/methods/18)
