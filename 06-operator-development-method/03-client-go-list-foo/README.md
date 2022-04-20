# List Foos with client-go

1. Create `CRD` for `Foo` in [01-crd](../01-crd)
    ```
    kubectl apply -f ../01-crd/foo.crd.yaml
    kubectl apply -f ../01-crd/foo.yaml
    ```
1. List `Foos` with client-go.
    1. Define `Foo` and `FooList` with `struct`.
        ```go

        ```
    1. Make `listFoos` func to list `Foos` with `dynamic.Interface`.
        ```go
        list, err := client.Resource(gvr).Namespace(namespace).List(context.Background(), metav1.ListOptions{}) // returns (*unstructured.UnstructuredList, error)
        ```

        [List for dynamicResourceClient](https://github.com/kubernetes/client-go/blob/28ccde769fc5519dd84e5512ebf303ac86ef9d7c/dynamic/simple.go#L272-L294):

        ```go
        func (c *dynamicResourceClient) List(ctx context.Context, opts metav1.ListOptions) (*unstructured.UnstructuredList, error) {
        	result := c.client.client.Get().AbsPath(c.makeURLSegments("")...).SpecificallyVersionedParams(&opts, dynamicParameterCodec,         versionV1).Do(ctx)
        	if err := result.Error(); err != nil {
        		return nil, err
        	}
        	retBytes, err := result.Raw()
        	if err != nil {
        		return nil, err
        	}
        	uncastObj, err := runtime.Decode(unstructured.UnstructuredJSONScheme, retBytes)
        	if err != nil {
        		return nil, err
        	}
        	if list, ok := uncastObj.(*unstructured.UnstructuredList); ok {
        		return list, nil
        	}

        	list, err := uncastObj.(*unstructured.Unstructured).ToList()
        	if err != nil {
        		return nil, err
        	}
        	return list, nil
        }
        ```

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
