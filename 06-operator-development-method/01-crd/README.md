# CustomResourceDefinition

## [What is custom resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources)

Before defining a *custom resource*, let's remember what a *resource* is:

> A *resource* is an endpoint in the Kubernetes API that stores a collection of API objects of a certain kind; for example, the built-in pods resource contains a collection of Pod objects.

Pod is a built-in resource.

> A *custom resource* is an extension of the Kubernetes API that is not necessarily available in a default Kubernetes installation.

You can define your own custom resource, and a *custom resource* is definied using the Kubernetes resource `CustomResourceDefinition` (ref: [Extend the Kubernetes API with CustomResourceDefinitions](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/)).

> custom resources let you store and retrieve structured data.

To make it fully work, we need a *custom controller* combined with your own custom resource.

## Example: Create Custom Resource `Foo`

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
