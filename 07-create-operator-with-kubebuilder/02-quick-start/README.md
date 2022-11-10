# 2. [Quick Start](https://book.kubebuilder.io/quick-start.html)

In this tutorial, we'll learn the followings:
1. Initialize an kubebuilder project with `kubebuilder` command
1. Create a custom resource `GuestBook` with `kubebuilder` command
1. Run the operator in your local
1. Run the operator in a Kubernetes cluster

## Versions

Checked version pairs:

|Docker|kind|kubernetes|kubebuilder|
|---|-----|---|---|
|[4.7.0 (77141)](https://docs.docker.com/desktop/mac/release-notes/#docker-desktop-471)|[v0.12.0](https://github.com/kubernetes-sigs/kind/releases/tag/v0.12.0)|v1.23.4|[v3.3.0](https://github.com/kubernetes-sigs/kubebuilder/releases/tag/v3.3.0)|
|[4.7.0 (77141)](https://docs.docker.com/desktop/mac/release-notes/#docker-desktop-471)|[v0.12.0](https://github.com/kubernetes-sigs/kind/releases/tag/v0.12.0)|v1.23.4|[v3.4.1](https://github.com/kubernetes-sigs/kubebuilder/releases/tag/v3.4.1)|
|[4.8.0 (78933)](https://docs.docker.com/desktop/release-notes/#docker-desktop-480)|[v0.17.0](https://github.com/kubernetes-sigs/kind/releases/tag/v0.17.0)|v1.24.0|[v3.5.0](https://github.com/kubernetes-sigs/kubebuilder/releases/tag/v3.5.0)|
|[4.13.1 (90346)](https://docs.docker.com/desktop/release-notes/#docker-desktop-4131)|[v0.17.0](https://github.com/kubernetes-sigs/kind/releases/tag/v0.17.0)|v1.25.3|[v3.6.0](https://github.com/kubernetes-sigs/kubebuilder/releases/tag/v3.6.0)|)
|[4.13.1 (90346)](https://docs.docker.com/desktop/release-notes/#docker-desktop-4131)|[v0.17.0](https://github.com/kubernetes-sigs/kind/releases/tag/v0.17.0)|v1.25.3|[v3.7.0](https://github.com/kubernetes-sigs/kubebuilder/releases/tag/v3.7.0)|)

## 2.1. Start a project

```
mkdir -p ~/projects/guestbook
cd ~/projects/guestbook
kubebuilder init --domain my.domain --repo my.domain/guestbook
```

If you are using [kubebuilder#v3.7.0](https://github.com/kubernetes-sigs/kubebuilder/releases/tag/v3.7.0) on M1 Mac, you might need to use `kubebuilder init --domain my.domain --repo my.domain/guestbook --plugins=go/v4-alpha`. (ref: [create a project](https://book.kubebuilder.io/quick-start.html#create-a-project))

<details><summary>Results</summary>

```
Writing kustomize manifests for you to edit...
Writing scaffold for you to edit...
Get controller runtime:
$ go get sigs.k8s.io/controller-runtime@v0.11.0
go: downloading sigs.k8s.io/controller-runtime v0.11.0
go: downloading k8s.io/apimachinery v0.23.0
go: downloading k8s.io/client-go v0.23.0
go: downloading k8s.io/utils v0.0.0-20210930125809-cb0fa318a74b
go: downloading sigs.k8s.io/structured-merge-diff/v4 v4.2.0
go: downloading k8s.io/component-base v0.23.0
go: downloading github.com/evanphx/json-patch v4.12.0+incompatible
go: downloading k8s.io/api v0.23.0
go: downloading golang.org/x/net v0.0.0-20210825183410-e898025ed96a
go: downloading github.com/prometheus/common v0.28.0
go: downloading k8s.io/apiextensions-apiserver v0.23.0
go: downloading golang.org/x/sys v0.0.0-20211029165221-6e7872819dc8
go: downloading github.com/fsnotify/fsnotify v1.5.1
Update dependencies:
$ go mod tidy
go: downloading go.uber.org/goleak v1.1.12
go: downloading golang.org/x/crypto v0.0.0-20210817164053-32db794688a5
go: downloading cloud.google.com/go v0.81.0
Next: define a resource with:
$ kubebuilder create api
```

</details>

Check generated files

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
│   ├── prometheus
│   │   ├── kustomization.yaml
│   │   └── monitor.yaml
│   └── rbac
│       ├── auth_proxy_client_clusterrole.yaml
│       ├── auth_proxy_role.yaml
│       ├── auth_proxy_role_binding.yaml
│       ├── auth_proxy_service.yaml
│       ├── kustomization.yaml
│       ├── leader_election_role.yaml
│       ├── leader_election_role_binding.yaml
│       ├── role_binding.yaml
│       └── service_account.yaml
├── go.mod
├── go.sum
├── hack
│   └── boilerplate.go.txt
└── main.go

6 directories, 24 files
```

1. Go
    1. go.mod
    1. go.sum
    1. main.go
1. Dockerfile
1. Makefile
1. config

## 2.2. Create an API

```
kubebuilder create api --group webapp --version v1 --kind Guestbook
```

```
Create Resource [y/n]
y
Create Controller [y/n]
y
Writing kustomize manifests for you to edit...
Writing scaffold for you to edit...
api/v1/guestbook_types.go
controllers/guestbook_controller.go
Update dependencies:
$ go mod tidy
Running make:
$ make generate
GOBIN=/Users/nakamasato/projects/guestbook/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.8.0
/Users/nakamasato/projects/guestbook/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
Next: implement your new API and generate the manifests (e.g. CRDs,CRs) with:
$ make manifests
```

<details><summary>If fails</summary>

```
Error: failed to create API: unable to run post-scaffold tasks of "base.go.kubebuilder.io/v3": exit status 2
```

```
kubebuilder create api --group webapp --version v1 --kind Guestbook

Create Resource [y/n]
y
Create Controller [y/n]
y
Writing kustomize manifests for you to edit...
Writing scaffold for you to edit...
api/v1/guestbook_types.go
controllers/guestbook_controller.go
Update dependencies:
$ go mod tidy
Running make:
$ make generate
go: creating new go.mod: module tmp
Downloading sigs.k8s.io/controller-tools/cmd/controller-gen@v0.8.0
go: downloading sigs.k8s.io/controller-tools v0.8.0
go: downloading golang.org/x/tools v0.1.6-0.20210820212750-d4cc65f0b2ff
go: added github.com/fatih/color v1.12.0
go: added github.com/go-logr/logr v1.2.0
go: added github.com/gobuffalo/flect v0.2.3
go: added github.com/gogo/protobuf v1.3.2
go: added github.com/google/go-cmp v0.5.6
go: added github.com/google/gofuzz v1.1.0
go: added github.com/inconshreveable/mousetrap v1.0.0
go: added github.com/json-iterator/go v1.1.12
go: added github.com/mattn/go-colorable v0.1.8
go: added github.com/mattn/go-isatty v0.0.12
go: added github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
go: added github.com/modern-go/reflect2 v1.0.2
go: added github.com/spf13/cobra v1.2.1
go: added github.com/spf13/pflag v1.0.5
go: added golang.org/x/mod v0.4.2
go: added golang.org/x/net v0.0.0-20210825183410-e898025ed96a
go: added golang.org/x/sys v0.0.0-20210831042530-f4d43177bf5e
go: added golang.org/x/text v0.3.7
go: added golang.org/x/tools v0.1.6-0.20210820212750-d4cc65f0b2ff
go: added golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
go: added gopkg.in/inf.v0 v0.9.1
go: added gopkg.in/yaml.v2 v2.4.0
go: added gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
go: added k8s.io/api v0.23.0
go: added k8s.io/apiextensions-apiserver v0.23.0
go: added k8s.io/apimachinery v0.23.0
go: added k8s.io/klog/v2 v2.30.0
go: added k8s.io/utils v0.0.0-20210930125809-cb0fa318a74b
go: added sigs.k8s.io/controller-tools v0.8.0
go: added sigs.k8s.io/json v0.0.0-20211020170558-c049b76a60c6
go: added sigs.k8s.io/structured-merge-diff/v4 v4.1.2
go: added sigs.k8s.io/yaml v1.3.0
/Users/nakamasato/projects/guestbook/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
bash: /Users/nakamasato/projects/guestbook/bin/controller-gen: No such file or directory
make: *** [generate] Error 127
Error: failed to create API: unable to run post-scaffold tasks of "base.go.kubebuilder.io/v3": exit status 2
Usage:
  kubebuilder create api [flags]

Examples:
  # Create a frigates API with Group: ship, Version: v1beta1 and Kind: Frigate
  kubebuilder create api --group ship --version v1beta1 --kind Frigate

  # Edit the API Scheme
  nano api/v1beta1/frigate_types.go

  # Edit the Controller
  nano controllers/frigate/frigate_controller.go

  # Edit the Controller Test
  nano controllers/frigate/frigate_controller_test.go

  # Generate the manifests
  make manifests

  # Install CRDs into the Kubernetes cluster using kubectl apply
  make install

  # Regenerate code and run against the Kubernetes cluster configured by ~/.kube/config
  make run


Flags:
      --controller           if set, generate the controller without prompting the user (default true)
      --force                attempt to create resource even if it already exists
      --group string         resource Group
  -h, --help                 help for api
      --kind string          resource Kind
      --make make generate   if true, run make generate after generating files (default true)
      --namespaced           resource is namespaced (default true)
      --plural string        resource irregular plural form
      --resource             if set, generate the resource without prompting the user (default true)
      --version string       resource Version

Global Flags:
      --plugins strings   plugin keys to be used for this subcommand execution

2022/04/27 05:56:30 failed to create API: unable to run post-scaffold tasks of "base.go.kubebuilder.io/v3": exit status 2
```

- Fixed in [✨ Remove deprecated go get from Makefile templates](https://github.com/kubernetes-sigs/kubebuilder/pull/2486)
- [v3.4.0](https://github.com/kubernetes-sigs/kubebuilder/releases/tag/v3.4.0) or later resolves the issue.

Quick fix for the older version (v3.3.0 or older) is updating Makefile:

```diff
- CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
- .PHONY: controller-gen
- controller-gen: ## Download controller-gen locally if necessary.
- 	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.8.0)
-
- KUSTOMIZE = $(shell pwd)/bin/kustomize
- .PHONY: kustomize
- kustomize: ## Download kustomize locally if necessary.
- 	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)
-
- ENVTEST = $(shell pwd)/bin/setup-envtest
- .PHONY: envtest
- envtest: ## Download envtest-setup locally if necessary.
- 	$(call go-get-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest@latest)
-
- # go-get-tool will 'go get' any package $2 and install it to $1.
- PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
- define go-get-tool
- @[ -f $(1) ] || { \
- set -e ;\
- TMP_DIR=$$(mktemp -d) ;\
- cd $$TMP_DIR ;\
- go mod init tmp ;\
- echo "Downloading $(2)" ;\
- GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
- rm -rf $$TMP_DIR ;\
- }
- endef
+ ##@ Build Dependencies
+
+ ## Location to install dependencies to
+ LOCALBIN ?= $(shell pwd)/bin
+ $(LOCALBIN): ## Ensure that the directory exists
+ 	mkdir -p $(LOCALBIN)
+
+ ## Tool Binaries
+ KUSTOMIZE ?= $(LOCALBIN)/kustomize
+ CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
+ ENVTEST ?= $(LOCALBIN)/setup-envtest
+
+ ## Tool Versions
+ KUSTOMIZE_VERSION ?= v3.8.7
+ CONTROLLER_TOOLS_VERSION ?= v0.8.0
+ ENVTEST_VERSION ?= latest
+
+ KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
+ .PHONY: kustomize
+ kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
+ $(KUSTOMIZE):
+ 	curl -s $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN)
+
+ .PHONY: controller-gen
+ controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
+ $(CONTROLLER_GEN):
+ 	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)
+
+ .PHONY: envtest
+ envtest: ## Download envtest-setup locally if necessary.
+ 	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@$(ENVTEST_VERSION)
```

</details>

Newly generated files:
- yaml
    - `config/crd/`:
    - `config/rbac/guestbook_editor_role.yaml`: `ClusterRole`
    - `config/rbac/guestbook_viewer_role.yaml`: `ClusterRole`
    - `config/samples/`: Sample yaml for `Guestbook`
- `api/`: Go types for the newly added custom resource `GuestBook`.
- `controllers/`:
    - `guestbook_controller.go`: controller for GuestBook
    - `suite_test.go`: Ginkgo test suite file.

Modified files:
- PROJECT: Added api resource for Guestbook.
- go.mod: Packages update. (e.g. Ginkgo, Gomega...)
- main.go:
    1. New resource added to scheme.
    1. Initialize `GuestbookReconciler` with Manager.

```
make manifests
```

With this command, controller-gen generates the following files:
- `config/crd/bases/webapp.my.domain_guestbooks.yaml`: `CRD` file for `Guestbook`
- `config/rbac/role.yaml`: `ClusterRole` for the operator (`manager-role`)

`Makefile` is a collection of frequently used commands. For more details, we'll see it in a later section.

## 2.3. Run the operator in your local

1. Start your Kubernetes cluster with kind. (You can use any Kubernetes cluster)
    ```
    kind create cluster
    ```

1. Ensure `kubeconfig` is correct.
    ```
    kubectl cluster-info
    ```

1. Install your CRD.
    ```bash
    make install # internally just run `kustomize build config/crd | kubectl apply -f -`
    ```

    <details><summary>error</summary>

    If you encounter the following error:

    ```
    make install
    /Users/m.naka/projects/guestbook/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
    curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash -s -- 3.8.7 /Users/m.naka/projects/guestbook/bin
    Version v3.8.7 does not exist or is not available for darwin/arm64.
    make: *** [/Users/m.naka/projects/guestbook/bin/kustomize] Error 1
    ```

    You can specify kustomize version or you can update the default KUSTOMIZE_VERSION in Makefile.

    ```
    KUSTOMIZE_VERSION=4.5.5 make install
    ```

    </details>

1. Run your controller in your local. (internally `go run main.go`)
    ```
    make run
    ```

    (Keep it running)

    <details><summary>stdout</summary>

    ```
    /Users/nakamasato/projects/guestbook/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
    /Users/nakamasato/projects/guestbook/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
    go fmt ./...
    go vet ./...
    go run ./main.go
    1.6510153771436281e+09  INFO    controller-runtime.metrics      Metrics server is starting to listen    {"addr": ":8080"}
    1.651015377144557e+09   INFO    setup   starting manager
    1.651015377145022e+09   INFO    Starting server {"path": "/metrics", "kind": "metrics", "addr": "[::]:8080"}
    1.651015377145098e+09   INFO    Starting server {"kind": "health probe", "addr": "[::]:8081"}
    1.651015377145259e+09   INFO    controller.guestbook    Starting EventSource    {"reconciler group": "webapp.my.domain", "reconciler kind": "Guestbook", "source": "kind source: *v1.Guestbook"}
    1.6510153771452842e+09  INFO    controller.guestbook    Starting Controller     {"reconciler group": "webapp.my.domain", "reconciler kind": "Guestbook"}
    1.651015377459231e+09   INFO    controller.guestbook    Starting workers        {"reconciler group": "webapp.my.domain", "reconciler kind": "Guestbook", "worker count": 1}
    ```

    </details>

1. Create an instance of Custom Resource `GuestBook`.

    (Note that we need to keep `make run` running when executing the following command.)

    ```
    kubectl apply -f config/samples/
    ```

    Nothing happens at this point as we haven't implemented anything yet. We'll implement a logic to capture create, update, and delete operation of the custom resource later.

1. Clean up the created insatnce.

    ```
    kubectl delete -f config/samples/
    ```

1. Stop `make run` by `ctrl-C`.
1. Remove CRD from the cluster (if you want to clean up completely)
    ```
    make uninstall
    ```

Let's wrap up what we've done:
1. Prepare a Kubernetes cluster.
1. Install CRD.
1. Run the operator.
1. Create an instance/object of the Custom Resource.
1. Delete the instance/object of the Custom Resource.
1. Stop the operator.
1. Uninstall CRD.

We use local run for development.

## 2.4. Run the operator in a Kubernetes cluster

1. Prepare a Kubernetes cluster (same as above)
    You can skip this step if you already set up your `kind` cluster above.

1. Install CRD (same as above).
    ```
    make install
    ```

    <details><summary>error</summary>

    If you encounter the following error:

    ```
    make install
    /Users/m.naka/projects/guestbook/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
    curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash -s -- 3.8.7 /Users/m.naka/projects/guestbook/bin
    Version v3.8.7 does not exist or is not available for darwin/arm64.
    make: *** [/Users/m.naka/projects/guestbook/bin/kustomize] Error 1
    ```

    You can specify kustomize version or you can update the default KUSTOMIZE_VERSION in Makefile.

    ```
    KUSTOMIZE_VERSION=4.5.5 make install
    ```

    </details>

1. Run the operator **on the Kubernetes cluster**.
    1. Set `IMG` variable.
        ```
        IMG=guestbook-controller:test
        ```

    1. To deploy to a Kubernetes cluster, we need a Docker image.
        ```
        make docker-build IMG=$IMG
        ```

        If you have a image registery you can just run `make docker-build docker-push IMG=<some-registry>/<project-name>:tag` as is written in the [book.kubebuilder.io](https://book.kubebuilder.io/quick-start.html#run-it-on-the-cluster).

    1. If you just want to deploy to `kind` cluster. There's a [trick to use the locally built image in a kind cluster](https://kind.sigs.k8s.io/docs/user/quick-start/#loading-an-image-into-your-cluster).

        ```
        kind load docker-image $IMG
        ```

        Note that the image needs to be loaded from the local node in kind cluster, to do so:
        - don't use a `:latest` tag or/and
        - specify `imagePullPolicy: IfNotPresent` or `imagePullPolicy: Never` on your container(s).

    1. Run the operator in the Kubernetes cluster.
        ```
        make deploy IMG=$IMG
        ```

        <details><summary>error</summary>

        If you encounter the following error:

        ```
        make install
        /Users/m.naka/projects/guestbook/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
        curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash -s -- 3.8.7 /Users/m.naka/projects/guestbook/bin
        /Users/m.naka/projects/guestbook/bin/kustomize exists. Remove it first.
        make: *** [/Users/m.naka/projects/guestbook/bin/kustomize] Error 1
        ```

        Delete the kustomize by `rm bin/kustomize` and rerun the command

        ```
        KUSTOMIZE_VERSION=4.5.5 make deploy IMG=$IMG
        ```

        </details>

    1. Check Pods in `guestbook-system` namespace.
        ```
        kubectl get po -n guestbook-system
        NAME                                            READY   STATUS    RESTARTS   AGE
        guestbook-controller-manager-6897c6457c-b8dtz   2/2     Running   0          2m26s
        ```

1. Create an instance/object of the Custom Resource. (same as above)
    ```
    kubectl apply -f config/samples/
    ```

1. Delete the instance/object of the Custom Resource. (same as above)
    ```
    kubectl delete -f config/samples/
    ```

1. Stop the operator.
    ```
    make undeploy
    ```

## 2.6. Wrap-up

### `kubebuilder` commands
1. Create a project: `kubebuilder init --domain my.domain --repo my.domain/guestbook`
1. Create an API: `kubebuilder create api --group webapp --version v1 --kind Guestbook`

### Commands for Operator Deployment
1. Prepare a Kubernetes cluster. (`kind create cluster` for local)
1. Install CRD. `make install`
1. Build and push Docker image.
    - with docker registry: `make docker-build docker-push IMG=<IMG>`
    - with local docker image in kind cluster: `make docker-build IMG=<IMG>` + `kind load docker-image <IMG>`
1. Run the operator.
    - run in local: `make run`
    - run in Kubernetes cluster: `make deploy IMG=<IMG>`
1. Create an instance/object of the Custom Resource. `kubectl apply -f config/samples/`
1. Delete the instance/object of the Custom Resource. `kubectl delete -f config/samples/`
1. Stop the operator.
    - running in local: `ctrl + C`
    - running in Kubernetes: `make undeploy`
1. Uninstall CRD. `make uninstall`

Cheatsheet:
1. [kind](../../99-cheatsheet/kind)
1. [kubebuilder](../../99-cheatsheet/kubebuilder): `kubebuilder` and `make` commands
