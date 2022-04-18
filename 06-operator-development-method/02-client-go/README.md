# Client Go

1. Install go from https://go.dev/doc/install.
1. Init go module (ref: https://go.dev/ref/mod#go-mod-init)

    ```
    go mod init foo
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
