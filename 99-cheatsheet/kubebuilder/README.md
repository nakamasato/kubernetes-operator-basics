# Cheatsheet - kubebuilder

## `kubebuilder`
|Command|Explanation|
|---|---|
|`kubebuilder init --domain <domain> --repo <repo>` |Initialize a project|
|`kubebuilder create api --group <group> --version <version> --kind <kind> --controller --resource`|Create new API resource|
|`kubebuilder create webhook --group <group> secret --version <version> --kind <kind> <option>`|Create new webhook. option: `--conversion`, `defaulting`, `--programmatic-validation`|

## `make`
|Command|Explanation|
|---|---|
|`make install` and `make uninstall`|Install and uninstall CRD|
|`make run`|Run controller|
|`make docker-build docker-push IMG=<IMG>`|docker build and push|
|`make deploy IMG=<IMG>` and `make undeploy`| Deploy and undeploy the operator (CRD & controller) to the Kubernetes cluster|
|`make fmt` | Format Go files|
|`make test`| Run tests with `envtest`|
|`make manifests`| Generate Kubernetes yaml files (CRD, rbac, etc.) |
|`make generate`| Generate Go code for API resource with DeepCopy implementations|
