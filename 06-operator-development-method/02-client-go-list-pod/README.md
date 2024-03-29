# List Pods with client-go

1. Install Go from https://go.dev/doc/install.
    ```
    go version
    ```
1. Create `main.go`.
    ```go
    pods, _ := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
    ```

    [List function for Pods](https://github.com/kubernetes/client-go/blob/34d1b64e67dbec9cbfa9628109bdd7c913897f49/kubernetes/typed/core/v1/pod.go#L88-L102):
    ```go
    func (c *pods) List(ctx context.Context, opts metav1.ListOptions) (result *v1.PodList, err error) {
        var timeout time.Duration
        if opts.TimeoutSeconds != nil {
            timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
        }
        result = &v1.PodList{}
        err = c.client.Get().
            Namespace(c.ns).
            Resource("pods").
            VersionedParams(&opts, scheme.ParameterCodec).
            Timeout(timeout).
            Do(ctx).
            Into(result)
        return
    }
    ```

1. Init go module (ref: https://go.dev/ref/mod#go-mod-init)

    ```
    go mod init listpod
    go mod tidy
    ```
1. Make sure k8s cluster is running. (or start a kind cluster with `kind create cluster`)
1. Run.
    ```
    go run main.go
    ```

    (Optional) You can also build and run.
    ```
    go build main.go
    ./main
    ```

    <details><summary>Result</summary>

    ```
    INDEX   NAMESPACE       NAME
    0       kube-system     coredns-64897985d-dgjpv
    1       kube-system     coredns-64897985d-l4qdf
    2       kube-system     etcd-kind-control-plane
    3       kube-system     kindnet-588g9
    4       kube-system     kube-apiserver-kind-control-plane
    5       kube-system     kube-controller-manager-kind-control-plane
    6       kube-system     kube-proxy-pzmnt
    7       kube-system     kube-scheduler-kind-control-plane
    8       local-path-storage      local-path-provisioner-5ddd94ff66-628dq
    ```

    </details>

## Go libraries
- https://pkg.go.dev/k8s.io/client-go: Go clients for talking to a kubernetes cluster.
- https://pkg.go.dev/k8s.io/apimachinery: Scheme, typing, encoding, decoding, and conversion packages for Kubernetes and Kubernetes-like API objects.
