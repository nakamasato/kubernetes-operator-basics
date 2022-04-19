# CustomResourceDefinition

1. `Foo` is not available before creating CRD for it.
    ```
    kubectl get foo
    error: the server doesn't have a resource type "foo"
    ```
1. Create `CustomResourceDefinition` for `Foo`.
    ```
    kubectl apply -f foo.crd.yaml
    ```
1. Create a custom resource `Foo` object named `test`.
    ```
    kubectl apply -f foo.yaml
    ```
    1. `apiVersion` needs to specify one of the defined versions.
1. Get the `Foo`.
    ```
    kubectl get foo test
    ```
    Result:
    ```
    NAME   AGE
    test   5m51s
    ```
1. Print `TestString` when getting `Foo` objects.

    Add the following piece of codes to the version:
    ```yaml
    additionalPrinterColumns:
      - name: Test String
        jsonPath: .testString
        type: string
    ```
1. Delete the created custom resource `Foo` object named `test`.
    ```
    kubectl delete -f foo.yaml
    ```
1. Delete `CustomResourceDefinition` for `Foo`.
    ```
    kubectl delete -f foo.crd.yaml
    ```
