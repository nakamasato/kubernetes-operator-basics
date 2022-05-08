# [Memcached Operator](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/)

## Versions
- Go: go1.18.1
- operator-sdk: v1.19.1

## 1. Initialize a project with `operator-sdk init`

```
mkdir -p ~/projects/memcached-operator && cd ~/projects/memcached-operator
git init
```

```
operator-sdk init --domain example.com --repo github.com/example/memcached-operator
```

<details>

```
tree .
.
├── Dockerfile
├── Makefile
├── PROJECT
├── config
│   ├── default
│   │   ├── kustomization.yaml
│   │   ├── manager_auth_proxy_patch.yaml
│   │   └── manager_config_patch.yaml
│   ├── manager
│   │   ├── controller_manager_config.yaml
│   │   ├── kustomization.yaml
│   │   └── manager.yaml
│   ├── manifests
│   │   └── kustomization.yaml
│   ├── prometheus
│   │   ├── kustomization.yaml
│   │   └── monitor.yaml
│   ├── rbac
│   │   ├── auth_proxy_client_clusterrole.yaml
│   │   ├── auth_proxy_role.yaml
│   │   ├── auth_proxy_role_binding.yaml
│   │   ├── auth_proxy_service.yaml
│   │   ├── kustomization.yaml
│   │   ├── leader_election_role.yaml
│   │   ├── leader_election_role_binding.yaml
│   │   ├── role_binding.yaml
│   │   └── service_account.yaml
│   └── scorecard
│       ├── bases
│       │   └── config.yaml
│       ├── kustomization.yaml
│       └── patches
│           ├── basic.config.yaml
│           └── olm.config.yaml
├── go.mod
├── go.sum
├── hack
│   └── boilerplate.go.txt
└── main.go

10 directories, 29 files
```

</details>

`operator-sdk` generates exactly the same files as `kubebuilder`.
As we saw in the previous section with `kubebuilder`, the `main.go` initializes and runs a `Manager`.
When you initialize a Manager, you can specify namespace to restrict which namespace to monitor by the operator. You can also set `""` to monitor all namespaces.

## 2. Create API resource and controller with `operator-sdk create api`

```
operator-sdk create api --group cache --version v1alpha1 --kind Memcached --resource --controller
```

<details><summary>If fails</summary>

If you encountered the following error

```
/Users/nakamasato/projects/memcached-operator/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
bash: line 1: /Users/nakamasato/projects/memcached-operator/bin/controller-gen: No such file or directory
make: *** [generate] Error 127
Error: failed to create API: unable to run post-scaffold tasks of "base.go.kubebuilder.io/v3": exit status 2
Usage:
```

You can fix it by replacing installation steps of controller-gen, kustomize, and envtest with the following codes in `Makefile`:

```makefile
##@ Build Dependencies
## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN): ## Ensure that the directory exists
    mkdir -p $(LOCALBIN)
## Tool Binaries
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest
## Tool Versions
KUSTOMIZE_VERSION ?= v3.8.7
CONTROLLER_TOOLS_VERSION ?= v0.8.0
ENVTEST_VERSION ?= latest
KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE):
    curl -s $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN)
.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN):
    GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)
.PHONY: envtest
envtest: ## Download envtest-setup locally if necessary.
    GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@$(ENVTEST_VERSION)
```

The diff is something like this:

```diff
        $(KUSTOMIZE) build config/default | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

-CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
-.PHONY: controller-gen
-controller-gen: ## Download controller-gen locally if necessary.
-       $(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.8.0)
-
-KUSTOMIZE = $(shell pwd)/bin/kustomize
+##@ Build Dependencies
+## Location to install dependencies to
+LOCALBIN ?= $(shell pwd)/bin
+$(LOCALBIN): ## Ensure that the directory exists
+       mkdir -p $(LOCALBIN)
+## Tool Binaries
+KUSTOMIZE ?= $(LOCALBIN)/kustomize
+CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
+ENVTEST ?= $(LOCALBIN)/setup-envtest
+## Tool Versions
+KUSTOMIZE_VERSION ?= v3.8.7
+CONTROLLER_TOOLS_VERSION ?= v0.8.0
+ENVTEST_VERSION ?= latest
+KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
 .PHONY: kustomize
-kustomize: ## Download kustomize locally if necessary.
-       $(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)
-
-ENVTEST = $(shell pwd)/bin/setup-envtest
+kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
+$(KUSTOMIZE):
+       curl -s $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN)
+.PHONY: controller-gen
+controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
+$(CONTROLLER_GEN):
+       GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)
 .PHONY: envtest
 envtest: ## Download envtest-setup locally if necessary.
-       $(call go-get-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest@latest)
-
-# go-get-tool will 'go get' any package $2 and install it to $1.
-PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
-define go-get-tool
-@[ -f $(1) ] || { \
-set -e ;\
-TMP_DIR=$$(mktemp -d) ;\
-cd $$TMP_DIR ;\
-go mod init tmp ;\
-echo "Downloading $(2)" ;\
-GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
-rm -rf $$TMP_DIR ;\
-}
-endef
+       GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@$(ENVTEST_VERSION)
```

You can commit it with `git add Makefile && git commit -m "Fix Makefile"`.

Then try the command again with `--force`

```
operator-sdk create api --group cache --version v1alpha1 --kind Memcached --resource --controller --force
```

</details>

```
make manifests
```

this commands generates:
- `config/rbac/role.yaml`: `Role` for the controller to access to the new resource `Memcached`
- `config/crd/bases`: CRD for the new resource `Memcached`

Let's make a commit for the newly generated resource and controller. `git add . && git commit -m "Create API resource and controller"`

In general, it is recommended to have one controller for one API resource. (e.g. memcached controller & memcached resource)

## 3. Define API Memcached

- `MemcachedSpec.Size`: replicas of memcached Deployment
- `MemcachedStatus.Nodes`: the names of the memcached Pods

Update `MemcachedSpec` and `MemcachedStatus` in `api/v1alpha1/memcached_types.go`

```go
// MemcachedSpec defines the desired state of Memcached
type MemcachedSpec struct {
	//+kubebuilder:validation:Minimum=0
	// Size is the size of the memcached deployment
	Size int32 `json:"size"`
}

// MemcachedStatus defines the observed state of Memcached
type MemcachedStatus struct {
	// Nodes are the names of the memcached pods
	Nodes []string `json:"nodes"`
}
```

Let's update the codes automatically generated based on the Go types.

```bash
make fmt generate manifests
```

- `fmt`: format go codes
- `generate`: go types -> zz_generated.deepcopy.go
- `manifests`: go types & marker -> yaml (crd, rbac...)

Let's commit the changes: `git add . && git commit -m "Design API"`

## 4. Implement Controller

### 4.1. Implement Controller - Fetch the Memcached instance

### 4.2. Implement Controller - Check if the deployment already exists, and create one if not exists

### 4.3. Implement Controller - Ensure the deployment size is the same as the spec

### 4.4. Implement Controller - Update the Memcached status with the pod names

## 5. Deploy with Deployment

## 6. Deploy with OLM

## 7. Write a test


## Versions

Checked version pairs:

|Docker|kind|kubernetes|operator-sdk|
|---|-----|---|---|
|[4.7.0 (77141)](https://docs.docker.com/desktop/mac/release-notes/#docker-desktop-471)|[v0.12.0](https://github.com/kubernetes-sigs/kind/releases/tag/v0.12.0)|v1.23.4|[v1.19.1](https://github.com/operator-framework/operator-sdk/releases/tag/v1.19.1)|
