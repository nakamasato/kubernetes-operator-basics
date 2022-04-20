# List Foos with client-go

1. Create `CRD` for `Foo` in [01-crd](../01-crd)
    ```
    kubectl apply -f ../01-crd/foo.crd.yaml
    kubectl apply -f ../01-crd/foo.yaml
    ```
1. List `Foos` with client-go.
    1. Define `Foo` and `FooList` with `struct`.
    1. Make `listFoos` func to list `Foos` with `dynamic.Interface`.
1. Run.
    ```
    go run main.go
    ```

    Result:
    ```
    INDEX   NAMESPACE       NAME
    0       default test
    ```

## Go libraries
- https://pkg.go.dev/k8s.io/client-go/dynamic
