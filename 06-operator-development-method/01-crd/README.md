# CustomResourceDefinition

1. Create `CustomResourceDefinition` for `Foo`.
    ```
    kubectl apply -f foo.crd.yaml
    ```
1. Create a custom resource `Foo` object named `test`.
    ```
    kubectl apply -f foo.yaml
    ```
1. Delete the created custom resource `Foo` object named `test`.
    ```
    kubectl delete -f foo.yaml
    ```
1. Delete `CustomResourceDefinition` for `Foo`.
    ```
    kubectl delete -f foo.crd.yaml
    ```
