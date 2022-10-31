# [Memcached Operator](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/)

## Versions
- Go: go1.17.9
- operator-sdk: v1.20.1

## 0. memcached-operator overview

1. When custom resource `Memcached` is created, the controller creates a `Deployment` if it doesn't exist.
1. Ensure that the Deployment's `replicas` field is same as `Memcached`(CR)'s `size` field.
1. Update the `Memcached` CR status with the names of the Pods created by the corresponding `Deployment`

![](memcached-operator.drawio.svg)

## 1. [operator-sdk] Init project with `operator-sdk init`

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
When you initialize a Manager, you can specify namespace to restrict which namespace to monitor by the operator. e.g. `Namespace: "some-ns"`. The default is all namespaces. For more details: [manager#Options](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/manager#Options)

Commit:

```
git add . && git commit -m "1. [operator-sdk] Init project with \`operator-sdk init\`"
```

## 2. [operator-sdk] Create API Memcached (Controller & Resource) with `operator-sdk create api`

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

Commit:

```
git add Makefile && git commit -m "1. Fix Makefile"
```

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

Let's make a commit for the newly generated resource and controller.

```
git add . && git commit -m "2. [operator-sdk] Create API Memcached (Controller & Resource) with \`operator-sdk create api\`"
```

In general, it is recommended to have one controller for one API resource. (e.g. memcached controller & memcached resource)

## 3. [API] Define API Memcached

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

Let's commit the changes:

```
git add . && git commit -m "3. [API] Define API Memcached"
```

## 4. Implement Controller

### 4.1. [Controller] Fetch the Memcached instance

1. Add necessary package.
    ```go
    import (
        "context" // already imported

        "k8s.io/apimachinery/pkg/api/errors"
        ///...
    )
    ```

1. Write the following lines in `Reconcile` function in [controllers/memcached_controller.go]().

    ```go
    func (r *MemcachedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result,     error) {
        log := log.FromContext(ctx)

        // 1. Fetch the Memcached instance
        memcached := &cachev1alpha1.Memcached{}
        err := r.Get(ctx, req.NamespacedName, memcached)
        if err != nil {
            if errors.IsNotFound(err) {
                log.Info("1. Fetch the Memcached instance. Memcached resource not found.     Ignoring since object must be deleted")
                return ctrl.Result{}, nil
            }
            // Error reading the object - requeue the request.
            log.Error(err, "1. Fetch the Memcached instance. Failed to get Mmecached")
            return ctrl.Result{}, err
        }
        log.Info("1. Fetch the Memcached instance. Memchached resource found", "memcached.Name",     memcached.Name, "memcached.Namespace", memcached.Namespace)
        return ctrl.Result{}, nil
    }
    ```

1. Check
    1. Install CRD and run the controller.
        ```bash
        make install run
        ```
    1. Apply a `Memcached` (CR).
        ```bash
        kubectl apply -f config/samples/cache_v1alpha1_memcached.yaml
        ```
    1. Check logs.

        ```bash
        2021-12-10T12:14:10.123+0900    INFO    controller.memcached    1. Fetch the Memcached instance. Memchached resource found {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Name": "memcached-sample", "memcached.Namespace": "default"}
        ```

    1. Delete the CR.
        ```bash
        kubectl delete -f config/samples/cache_v1alpha1_memcached.yaml
        ```

    1. Check logs.
        ```bash
        2021-12-10T12:15:37.234+0900    INFO    controller.memcached    1. Fetch the Memcached instance. Memcached resource not found. Ignoring since object must be deleted       {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default"}
        ```
    1. Stop the controller.

Commit:
```
git add . && git commit -m "4.1. [Controller] Fetch the Memcached instance"
```
### 4.2. [Controller] Check if the deployment already exists, and create one if not exists

1. Add necessary packages to `import`.
    ```go
    import (
        ...
        "k8s.io/apimachinery/pkg/types"
        ...

        appsv1 "k8s.io/api/apps/v1"
        corev1 "k8s.io/api/core/v1"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

        ...
    )
    ```

1. Add the following logics to `Reconcile` function.

    ```go
    // 2. Check if the deployment already exists, if not create a new one
    found := &appsv1.Deployment{}
    err = r.Get(ctx, types.NamespacedName{Name: memcached.Name, Namespace: memcached.Namespace}, found)
    if err != nil && errors.IsNotFound(err) {
            // Define a new deployment
            dep := r.deploymentForMemcached(memcached)
            log.Info("2. Check if the deployment already exists, if not create a new one. Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
            err = r.Create(ctx, dep)
            if err != nil {
                    log.Error(err, "2. Check if the deployment already exists, if not create a new one. Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
                    return ctrl.Result{}, err
            }
            // Deployment created successfully - return and requeue
            return ctrl.Result{Requeue: true}, nil
    } else if err != nil {
            log.Error(err, "2. Check if the deployment already exists, if not create a new one. Failed to get Deployment")
            return ctrl.Result{}, err
    }
    ```
1. Create `deploymentForMemcached` and `labelsForMemcached` functions.

    <details><summary>deploymentForMemcached</summary>

    ```go
    // deploymentForMemcached returns a memcached Deployment object
    func (r *MemcachedReconciler) deploymentForMemcached(m *cachev1alpha1.Memcached) *appsv1.Deployment {
        ls := labelsForMemcached(m.Name)
        replicas := m.Spec.Size

        dep := &appsv1.Deployment{
                ObjectMeta: metav1.ObjectMeta{
                        Name:      m.Name,
                        Namespace: m.Namespace,
                },
                Spec: appsv1.DeploymentSpec{
                        Replicas: &replicas,
                        Selector: &metav1.LabelSelector{
                                MatchLabels: ls,
                        },
                        Template: corev1.PodTemplateSpec{
                                ObjectMeta: metav1.ObjectMeta{
                                        Labels: ls,
                                },
                                Spec: corev1.PodSpec{
                                        Containers: []corev1.Container{{
                                                Image:   "memcached:1.4.36-alpine",
                                                Name:    "memcached",
                                                Command: []string{"memcached", "-m=64", "-o", "modern", "-v"},
                                                Ports: []corev1.ContainerPort{{
                                                        ContainerPort: 11211,
                                                        Name:          "memcached",
                                                }},
                                        }},
                                },
                        },
                },
        }
        // Set Memcached instance as the owner and controller
        ctrl.SetControllerReference(m, dep, r.Scheme)
        return dep
    }
    ```

    </details>

    <details><summary>labelsForMemcached</summary>

    ```go
    // labelsForMemcached returns the labels for selecting the resources
    // belonging to the given memcached CR name.
    func labelsForMemcached(name string) map[string]string {
        return map[string]string{"app": "memcached", "memcached_cr": name}
    }
    ```

    </details>
1. Add necessary `RBAC` to the reconciler.

    ```diff
    //+kubebuilder:rbac:groups=cache.example.com,resources=memcacheds,verbs=get;list;watch;create;update;patch;delete
    //+kubebuilder:rbac:groups=cache.example.com,resources=memcacheds/status,verbs=get;update;patch
    //+kubebuilder:rbac:groups=cache.example.com,resources=memcacheds/finalizers,verbs=update
    + //+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
    ```

1. Add `Owns(&appsv1.Deployment{})` to the controller manager.

    ```go
    // SetupWithManager sets up the controller with the Manager.
    func (r *MemcachedReconciler) SetupWithManager(mgr ctrl.Manager) error {
        return ctrl.NewControllerManagedBy(mgr).
            For(&cachev1alpha1.Memcached{}).
            Owns(&appsv1.Deployment{}).
            Complete(r)
    }
    ```

1. Check
    1. Run the controller.
        ```bash
        make run
        ```
    1. Apply a `Memcached` (CR).
        ```bash
        kubectl apply -f config/samples/cache_v1alpha1_memcached.yaml
        ```
    1. Check logs.

        ```bash
        2021-12-10T12:34:34.587+0900    INFO    controller.memcached    1. Fetch the Memcached instance. Memchached resource found      {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Name": "memcached-sample", "memcached.Namespace": "default"}
        2021-12-10T12:34:34.587+0900    INFO    controller.memcached    2. Check if the deployment already exists, if not create a new one. Creating a new Deployment       {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "Deployment.Namespace": "default", "Deployment.Name": "memcached-sample"}
        2021-12-10T12:34:34.599+0900    INFO    controller.memcached    1. Fetch the Memcached instance. Memchached resource found      {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Name": "memcached-sample", "memcached.Namespace": "default"}
        2021-12-10T12:34:34.604+0900    INFO    controller.memcached    1. Fetch the Memcached instance. Memchached resource found      {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Name": "memcached-sample", "memcached.Namespace": "default"}
        2021-12-10T12:34:34.648+0900    INFO    controller.memcached    1. Fetch the Memcached instance. Memchached resource found      {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Name": "memcached-sample", "memcached.Namespace": "default"}
        2021-12-10T12:34:34.662+0900    INFO    controller.memcached    1. Fetch the Memcached instance. Memchached resource found      {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Name": "memcached-sample", "memcached.Namespace": "default"}
        2021-12-10T12:34:34.724+0900    INFO    controller.memcached    1. Fetch the Memcached instance. Memchached resource found      {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Name": "memcached-sample", "memcached.Namespace": "default"}
        2021-12-10T12:34:43.285+0900    INFO    controller.memcached    1. Fetch the Memcached instance. Memchached resource found      {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Name": "memcached-sample", "memcached.Namespace": "default"}
        2021-12-10T12:34:46.333+0900    INFO    controller.memcached    1. Fetch the Memcached instance. Memchached resource found      {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Name": "memcached-sample", "memcached.Namespace": "default"}
        2021-12-10T12:34:48.363+0900    INFO    controller.memcached    1. Fetch the Memcached instance. Memchached resource found      {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Name": "memcached-sample", "memcached.Namespace": "default"}
        ```

        There are ten lines of logs:
        1. When `Memcached` object is created.
        1. Create `Deployment`.
        1. When `Deployment` is created.
        1. 8 more events are created accordingly.


    1. Check `Deployment`.

        ```
        kubectl get deploy memcached-sample
        NAME               READY   UP-TO-DATE   AVAILABLE   AGE
        memcached-sample   3/3     3            3           19s
        ```

    1. Delete the CR.
        ```bash
        kubectl delete -f config/samples/cache_v1alpha1_memcached.yaml
        ```

    1. Check logs.
        ```bash
        2021-12-10T12:38:50.473+0900    INFO    controller.memcached 1. Fetch the Memcached instance. Memcached resource not found. Ignoring since object must be deleted      {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default"}
        2021-12-10T12:38:50.512+0900    INFO    controller.memcached 1. Fetch the Memcached instance. Memcached resource not found. Ignoring since object must be deleted      {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default"}
        ```
    1. Check `Deployment`.
        ```
        kubectl get deploy
        No resources found in default namespace.
        ```
    1. Stop the controller.

Commit:
```
git add . && git commit -m "4.2. [Controller] Check if the deployment already exists, and create one if not exists"
```

### 4.3. [Controller] Ensure the deployment size is the same as the spec

1. Add the following lines to `Reconcile` function.

    ```go
    // 3. Ensure the deployment size is the same as the spec
    size := memcached.Spec.Size
    if *found.Spec.Replicas != size {
            found.Spec.Replicas = &size
            err = r.Update(ctx, found)
            if err != nil {
                    log.Error(err, "3. Ensure the deployment size is the same as the spec. Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
                    return ctrl.Result{}, err
            }
            // Spec updated - return and requeue
            log.Info("3. Ensure the deployment size is the same as the spec. Update deployment size", "Deployment.Spec.Replicas", size)
            return ctrl.Result{Requeue: true}, nil
    }
    ```
1. Check
    1. Run the controller.
        ```bash
        make run
        ```
    1. Apply a `Memcached` (CR).
        ```bash
        kubectl apply -f config/samples/cache_v1alpha1_memcached.yaml
        ```
    1. Check `Deployment`.

        ```
        kubectl get deploy memcached-sample
        NAME               READY   UP-TO-DATE   AVAILABLE   AGE
        memcached-sample   3/3     3            3           19s
        ```

    1. Change the size to 2 in [config/samples/cache_v1alpha1_memcached.yaml]()

        ```
        kubectl apply -f config/samples/cache_v1alpha1_memcached.yaml
        ```

    1. Check logs.

        ```bash
        2021-12-10T12:59:09.880+0900    INFO    controller.memcached    1. Fetch the Memcached instance. Memchached resource found      {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Name": "memcached-sample", "memcached.Namespace": "default"}
        2021-12-10T12:59:09.888+0900    INFO    controller.memcached    3. Ensure the deployment size is the same as the spec. Update deployment size{"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "Deployment.Spec.Replicas": 2}
        2021-12-10T12:59:09.888+0900    INFO    controller.memcached    1. Fetch the Memcached instance. Memchached resource found      {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Name": "memcached-sample", "memcached.Namespace": "default"}
        2021-12-10T12:59:09.894+0900    INFO    controller.memcached    1. Fetch the Memcached instance. Memchached resource found      {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Name": "memcached-sample", "memcached.Namespace": "default"}
        2021-12-10T12:59:09.911+0900    INFO    controller.memcached    1. Fetch the Memcached instance. Memchached resource found      {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Name": "memcached-sample", "memcached.Namespace": "default"}
        2021-12-10T12:59:09.951+0900    INFO    controller.memcached    1. Fetch the Memcached instance. Memchached resource found      {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Name": "memcached-sample", "memcached.Namespace": "default"}
        ```

    1. Check `Deployment`.

        ```
        kubectl get deploy
        NAME               READY   UP-TO-DATE   AVAILABLE   AGE
        memcached-sample   2/2     2            2           115s
        ```

    1. Delete the CR.
        ```bash
        kubectl delete -f config/samples/cache_v1alpha1_memcached.yaml
        ```

    1. Check logs.
        ```bash
        2021-12-10T13:00:50.149+0900    INFO    controller.memcached    1. Fetch the Memcached instance. Memcached resource not found. Ignoring since object must be deleted {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default"}
        2021-12-10T13:00:50.185+0900    INFO    controller.memcached    1. Fetch the Memcached instance. Memcached resource not found. Ignoring since object must be deleted {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default"}
        ```
    1. Check `Deployment`.
        ```
        kubectl get deploy
        No resources found in default namespace.
        ```
    1. Stop the controller.

Commit:

```
git commit -a -m "4.3. [Controller] Ensure the deployment size is the same as the spec"
```

### 4.4. [Controller] Update the Memcached status with the pod names

1. Add `"reflect"` to `import`.
1. Add the following logic to `Reconcile` functioin.

    ```go
    // 4. Update the Memcached status with the pod names
    // List the pods for this memcached's deployment
    podList := &corev1.PodList{}
    listOpts := []client.ListOption{
            client.InNamespace(memcached.Namespace),
            client.MatchingLabels(labelsForMemcached(memcached.Name)),
    }
    if err = r.List(ctx, podList, listOpts...); err != nil {
            log.Error(err, "4. Update the Memcached status with the pod names. Failed to list pods", "Memcached.Namespace", memcached.Namespace, "Memcached.Name", memcached.Name)
            return ctrl.Result{}, err
    }
    podNames := getPodNames(podList.Items)
    log.Info("4. Update the Memcached status with the pod names. Pod list", "podNames", podNames)
    // Update status.Nodes if needed
    if !reflect.DeepEqual(podNames, memcached.Status.Nodes) {
            memcached.Status.Nodes = podNames
            err := r.Status().Update(ctx, memcached)
            if err != nil {
                    log.Error(err, "4. Update the Memcached status with the pod names. Failed to update Memcached status")
                    return ctrl.Result{}, err
            }
    }
    log.Info("4. Update the Memcached status with the pod names. Update memcached.Status", "memcached.Status.Nodes", memcached.Status.Nodes)
    ```
1. Add `getPodNames` function.

    ```go
    // getPodNames returns the pod names of the array of pods passed in
    func getPodNames(pods []corev1.Pod) []string {
        var podNames []string
        for _, pod := range pods {
                podNames = append(podNames, pod.Name)
        }
        return podNames
    }
    ```
1. Add necessary `RBAC`.
    ```diff
      //+kubebuilder:rbac:groups=cache.example.com,resources=memcacheds,verbs=get;list;watch;create;update;patch;delete
      //+kubebuilder:rbac:groups=cache.example.com,resources=memcacheds/status,verbs=get;update;patch
      //+kubebuilder:rbac:groups=cache.example.com,resources=memcacheds/finalizers,verbs=update
      //+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
    + //+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;
    ```

1. Check
    1. Run the controller.
        ```bash
        make run
        ```
    1. Apply a `Memcached` (CR).
        ```bash
        kubectl apply -f config/samples/cache_v1alpha1_memcached.yaml
        ```

    1. Check logs.

        ```bash
        2021-12-10T13:09:03.716+0900    INFO    controller.memcached    1. Fetch the Memcached instance. Memchached resource found      {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Name": "memcached-sample", "memcached.Namespace": "default"}
        2021-12-10T13:09:03.716+0900    INFO    controller.memcached    2. Check if the deployment already exists, if not create a new one. Creating a new Deployment    {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "Deployment.Namespace": "default", "Deployment.Name": "memcached-sample"}
        2021-12-10T13:09:03.727+0900    INFO    controller.memcached    1. Fetch the Memcached instance. Memchached resource found      {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Name": "memcached-sample", "memcached.Namespace": "default"}
        2021-12-10T13:09:03.829+0900    INFO    controller.memcached    4. Update the Memcached status with the pod names. Pod list     {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "podNames": ["memcached-sample-6c765df685-f9jpl", "memcached-sample-6c765df685-cf725"]}
        2021-12-10T13:09:03.841+0900    INFO    controller.memcached    4. Update the Memcached status with the pod names. Update memcached.Status       {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Status.Nodes": ["memcached-sample-6c765df685-f9jpl", "memcached-sample-6c765df685-cf725"]}
        2021-12-10T13:09:03.841+0900    INFO    controller.memcached    1. Fetch the Memcached instance. Memchached resource found      {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Name": "memcached-sample", "memcached.Namespace": "default"}
        2021-12-10T13:09:03.841+0900    INFO    controller.memcached    4. Update the Memcached status with the pod names. Pod list     {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "podNames": ["memcached-sample-6c765df685-f9jpl", "memcached-sample-6c765df685-cf725"]}
        2021-12-10T13:09:03.841+0900    INFO    controller.memcached    4. Update the Memcached status with the pod names. Update memcached.Status       {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Status.Nodes": ["memcached-sample-6c765df685-f9jpl", "memcached-sample-6c765df685-cf725"]}
        2021-12-10T13:09:05.565+0900    INFO    controller.memcached    1. Fetch the Memcached instance. Memchached resource found      {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Name": "memcached-sample", "memcached.Namespace": "default"}
        2021-12-10T13:09:05.565+0900    INFO    controller.memcached    4. Update the Memcached status with the pod names. Pod list     {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "podNames": ["memcached-sample-6c765df685-f9jpl", "memcached-sample-6c765df685-cf725"]}
        2021-12-10T13:09:05.565+0900    INFO    controller.memcached    4. Update the Memcached status with the pod names. Update memcached.Status       {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Status.Nodes": ["memcached-sample-6c765df685-f9jpl", "memcached-sample-6c765df685-cf725"]}
        2021-12-10T13:09:05.587+0900    INFO    controller.memcached    1. Fetch the Memcached instance. Memchached resource found      {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Name": "memcached-sample", "memcached.Namespace": "default"}
        2021-12-10T13:09:05.587+0900    INFO    controller.memcached    4. Update the Memcached status with the pod names. Pod list     {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "podNames": ["memcached-sample-6c765df685-f9jpl", "memcached-sample-6c765df685-cf725"]}
        2021-12-10T13:09:05.588+0900    INFO    controller.memcached    4. Update the Memcached status with the pod names. Update memcached.Status       {"reconciler group": "cache.example.com", "reconciler kind": "Memcached", "name": "memcached-sample", "namespace": "default", "memcached.Status.Nodes": ["memcached-sample-6c765df685-f9jpl", "memcached-sample-6c765df685-cf725"]}
        ```

    1. Check `Deployment`.

        ```
        kubectl get deploy
        NAME               READY   UP-TO-DATE   AVAILABLE   AGE
        memcached-sample   2/2     2            2           115s
        ```

    1. Check `status` in `Memcached` object.

        ```bash
        kubectl get Memcached memcached-sample -o jsonpath='{.status}' | jq
        {
          "nodes": [
            "memcached-sample-6c765df685-9drvp",
            "memcached-sample-6c765df685-g7nl8"
          ]
        }
        ```

    1. Delete the CR.
        ```bash
        kubectl delete -f config/samples/cache_v1alpha1_memcached.yaml
        ```

    1. Stop the controller.

commit:

```
git commit -am "4.4. [Controller] Update the Memcached status with the pod names"
```

## 5. Write a test

### 5.1. Tools

1. [envtest](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/envtest): provides libraries for integration testing by starting a local control plane. (`etcd` and `kube-apiserver`)
1. [Ginkgo](https://pkg.go.dev/github.com/onsi/ginkgo): BDD framework.
1. [Gomega](https://pkg.go.dev/github.com/onsi/gomega): Matcher library for testing.
### 5.2. Prepare `suite_test.go`

1. Import necessary packages.
    ```diff
     import (
    +       "context"
            "path/filepath"
            "testing"
    +       ctrl "sigs.k8s.io/controller-runtime"
    +
            . "github.com/onsi/ginkgo"
            . "github.com/onsi/gomega"
            "k8s.io/client-go/kubernetes/scheme"
    -       "k8s.io/client-go/rest"
            "sigs.k8s.io/controller-runtime/pkg/client"
            "sigs.k8s.io/controller-runtime/pkg/envtest"
            "sigs.k8s.io/controller-runtime/pkg/envtest/ter"
            logf "sigs.k8s.io/controller-runtime/pkg/log"
            "sigs.k8s.io/controller-runtime/pkg/log/zap"
    +       "sigs.k8s.io/controller-runtime/pkg/manager"
    ```

1. Prepare global variables.
    ```diff
    -var cfg *rest.Config
    -var k8sClient client.Client
    -var testEnv *envtest.Environment
    +var (
    +       k8sClient  client.Client
    +       k8sManager manager.Manager
    +       testEnv    *envtest.Environment
    +       ctx        context.Context
    +       cancel     context.CancelFunc
    +)
    ```

1. Add the following lines at the end of `BeforeSuite` in `controllers/suite_test.go`.

    ```go
        // Create context with cancel.
        ctx, cancel = context.WithCancel(context.TODO())

        // Register the schema to manager.
        k8sManager, err = ctrl.NewManager(cfg, ctrl.Options{
            Scheme: scheme.Scheme,
        })

        // Initialize `MemcachedReconciler` with the manager client schema.
        err = (&MemcachedReconciler{
            Client: k8sManager.GetClient(),
            Scheme: k8sManager.GetScheme(),
        }).SetupWithManager(k8sManager)

        // Start the with a goroutine.
        go func() {
            defer GinkgoRecover()
            err = k8sManager.Start(ctx)
            Expect(err).ToNot(HaveOccurred(), "failed to run ger")
        }()
    ```

1. Add `cancel()` to AfterSuite.

    ```diff
     var _ = AfterSuite(func() {
    +       cancel()
            By("tearing down the test environment")
            err := testEnv.Stop()
            Expect(err).NotTo(HaveOccurred())
    ```

### 5.3. Write tests

Test cases in `controllers/memcached_controller_test.go`:

1. When `Memcached` is created
    1. `Deployment` should be created.
    1. `Memcached`'s nodes have pods' names.
1. When `Memcached`'s `size` is updated
    1. `Deployment`'s `replicas` should be updated.
    1. `Memcached`'s nodes have new pods' names.
1. When `Deployment` is updated
    1. Deleting `Deployment` -> `Deployment` is recreated.
    1. Updating `Deployment` with `replicas = 0` -> `Deployment`'s replicas is updated to the original number.

<details><summary>memcached_controller_test.go</summary>

```go
package controllers

import (
	"context"
	"fmt"
	"time"

	cachev1alpha1 "github.com/example/memcached-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	memcachedApiVersion = "cache.example.com/v1alphav1"
	memcachedKind       = "Memcached"
	memcachedName       = "sample"
	memcachedNamespace  = "default"
	memcachedStartSize  = int32(3)
	memcachedUpdateSize = int32(10)
	timeout             = time.Second * 10
	interval            = time.Millisecond * 250
)

var _ = Describe("Memcached controller", func() {

	lookUpKey := types.NamespacedName{Name: memcachedName, Namespace: memcachedNamespace}

	AfterEach(func() {
		// Delete Memcached
		deleteMemcached(ctx, lookUpKey)
		// Delete all Pods
		deleteAllPods(ctx)
	})

	Context("When creating Memcached", func() {
		BeforeEach(func() {
			// Create Memcached
			createMemcached(ctx, memcachedStartSize)
		})
		It("Should create Deployment with the specified size and memcached image", func() {
			// Deployment is created
			deployment := &appsv1.Deployment{}
			Eventually(func() error {
				return k8sClient.Get(ctx, lookUpKey, deployment)
			}, timeout, interval).Should(Succeed())
			Expect(*deployment.Spec.Replicas).Should(Equal(memcachedStartSize))
			Expect(deployment.Spec.Template.Spec.Containers[0].Image).Should(Equal("memcached:1.4.36-alpine"))
			// https://github.com/kubernetes-sigs/controller-runtime/blob/master/pkg/controller/controllerutil/controllerutil_test.go
			Expect(deployment.OwnerReferences).ShouldNot(BeEmpty())
		})
		It("Should have pods name in Memcached Node", func() {
			checkIfDeploymentExists(ctx, lookUpKey)

			By("By creating Pods with labels")
			podNames := createPods(ctx, int(memcachedStartSize))

			updateMemcacheSize(ctx, lookUpKey, memcachedUpdateSize) // just to trigger reconcile

			checkMemcachedStatusNodes(ctx, lookUpKey, podNames)
		})
	})

	Context("When updating Memcached", func() {
		BeforeEach(func() {
			// Create Memcached
			createMemcached(ctx, memcachedStartSize)
			// Deployment is ready
			checkDeploymentReplicas(ctx, lookUpKey, memcachedStartSize)
		})

		It("Should update Deployment replicas", func() {
			By("Changing Memcached size")
			updateMemcacheSize(ctx, lookUpKey, memcachedUpdateSize)

			checkDeploymentReplicas(ctx, lookUpKey, memcachedUpdateSize)
		})

		It("Should update the Memcached status with the pod names", func() {
			By("Changing Memcached size")
			updateMemcacheSize(ctx, lookUpKey, memcachedUpdateSize)

			podNames := createPods(ctx, int(memcachedUpdateSize))
			checkMemcachedStatusNodes(ctx, lookUpKey, podNames)
		})
	})
	Context("When changing Deployment", func() {
		BeforeEach(func() {
			// Create Memcached
			createMemcached(ctx, memcachedStartSize)
			// Deployment is ready
			checkDeploymentReplicas(ctx, lookUpKey, memcachedStartSize)
		})

		It("Should check if the deployment already exists, if not create a new one", func() {
			By("Deleting Deployment")
			deployment := &appsv1.Deployment{}
			Expect(k8sClient.Get(ctx, lookUpKey, deployment)).Should(Succeed())
			Expect(k8sClient.Delete(ctx, deployment)).Should(Succeed())

			// Deployment will be recreated by the controller
			checkIfDeploymentExists(ctx, lookUpKey)
		})

		It("Should ensure the deployment size is the same as the spec", func() {
			By("Changing Deployment replicas")
			deployment := &appsv1.Deployment{}
			Expect(k8sClient.Get(ctx, lookUpKey, deployment)).Should(Succeed())
			*deployment.Spec.Replicas = 0
			Expect(k8sClient.Update(ctx, deployment)).Should(Succeed())

			// replicas will be updated back to the original one by the controller
			checkDeploymentReplicas(ctx, lookUpKey, memcachedStartSize)
		})
	})

	// Deployment is expected to be deleted when Memcached is deleted.
	// As it's garbage collector's responsibility, which is not part of envtest, we don't test it here.
})

func checkIfDeploymentExists(ctx context.Context, lookUpKey types.NamespacedName) {
	deployment := &appsv1.Deployment{}
	Eventually(func() error {
		return k8sClient.Get(ctx, lookUpKey, deployment)
	}, timeout, interval).Should(Succeed())
}

func checkDeploymentReplicas(ctx context.Context, lookUpKey types.NamespacedName, expectedSize int32) {
	Eventually(func() (int32, error) {
		deployment := &appsv1.Deployment{}
		err := k8sClient.Get(ctx, lookUpKey, deployment)
		if err != nil {
			return int32(0), err
		}
		return *deployment.Spec.Replicas, nil
	}, timeout, interval).Should(Equal(expectedSize))
}

func newMemcached(size int32) *cachev1alpha1.Memcached {
	return &cachev1alpha1.Memcached{
		TypeMeta: metav1.TypeMeta{
			APIVersion: memcachedApiVersion,
			Kind:       memcachedKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      memcachedName,
			Namespace: memcachedNamespace,
		},
		Spec: cachev1alpha1.MemcachedSpec{
			Size: size,
		},
	}
}

func createMemcached(ctx context.Context, size int32) {
	memcached := newMemcached(size)
	Expect(k8sClient.Create(ctx, memcached)).Should(Succeed())
}

func updateMemcacheSize(ctx context.Context, lookUpKey types.NamespacedName, size int32) {
	memcached := &cachev1alpha1.Memcached{}
	Expect(k8sClient.Get(ctx, lookUpKey, memcached)).Should(Succeed())
	memcached.Spec.Size = size
	Expect(k8sClient.Update(ctx, memcached)).Should(Succeed())
}

func deleteMemcached(ctx context.Context, lookUpKey types.NamespacedName) {
	memcached := &cachev1alpha1.Memcached{}
	Expect(k8sClient.Get(ctx, lookUpKey, memcached)).Should(Succeed())
	Expect(k8sClient.Delete(ctx, memcached)).Should(Succeed())
}

func checkMemcachedStatusNodes(ctx context.Context, lookUpKey types.NamespacedName, podNames []string) {
	memcached := &cachev1alpha1.Memcached{}
	Eventually(func() ([]string, error) {
		err := k8sClient.Get(ctx, lookUpKey, memcached)
		if err != nil {
			return nil, err
		}
		return memcached.Status.Nodes, nil
	}, timeout, interval).Should(ConsistOf(podNames))
}

func createPods(ctx context.Context, num int) []string {
	podNames := []string{}
	for i := 0; i < num; i++ {
		podName := fmt.Sprintf("pod-%d", i)
		podNames = append(podNames, podName)
		pod := newPod(podName)
		Expect(k8sClient.Create(ctx, pod)).Should(Succeed())
	}
	return podNames
}

func deleteAllPods(ctx context.Context) {
	err := k8sClient.DeleteAllOf(ctx, &v1.Pod{}, client.InNamespace(memcachedNamespace))
	Expect(err).NotTo(HaveOccurred())
}

func newPod(name string) *v1.Pod {
	return &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: memcachedNamespace,
			Labels: map[string]string{
				"app":          "memcached",
				"memcached_cr": memcachedName,
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "memcached",
					Image: "memcached",
				},
			},
		},
		Status: v1.PodStatus{},
	}
}
```

</details>

### 5.4. Run the tests

```
make test
```

<details>

```
/Users/nakamasato/repos/nakamasato/memcached-operator/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
/Users/nakamasato/repos/nakamasato/memcached-operator/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
go vet ./...
GOBIN=/Users/nakamasato/repos/nakamasato/memcached-operator/bin go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
KUBEBUILDER_ASSETS="/Users/nakamasato/Library/Application Support/io.kubebuilder.envtest/k8s/1.23.3-darwin-amd64" go test ./... -coverprofile cover.out
?       github.com/example/memcached-operator   [no test files]
?       github.com/example/memcached-operator/api/v1alpha1      [no test files]
ok      github.com/example/memcached-operator/controllers       18.284s coverage: 79.3% of statements
```

</details>

## 6. Deployment

### 6.1. Deploy with Deployment

```bash
export IMG=memcached-operator:deploy
make docker-build IMG=$IMG
kind load docker-image $IMG # only necessary when deploying to kind cluster
make deploy
```

### 6.2. Deploy with OLM

## Versions

Checked version pairs:

|Docker|kind|kubernetes|operator-sdk|
|---|-----|---|---|
|[4.7.0 (77141)](https://docs.docker.com/desktop/mac/release-notes/#docker-desktop-471)|[v0.12.0](https://github.com/kubernetes-sigs/kind/releases/tag/v0.12.0)|v1.23.4|[v1.19.1](https://github.com/operator-framework/operator-sdk/releases/tag/v1.19.1)|
|[4.7.0 (77141)](https://docs.docker.com/desktop/mac/release-notes/#docker-desktop-471)|[v0.12.0](https://github.com/kubernetes-sigs/kind/releases/tag/v0.12.0)|v1.23.4|[v1.20.1](https://github.com/operator-framework/operator-sdk/releases/tag/v1.21.0)|
|[4.7.0 (77141)](https://docs.docker.com/desktop/mac/release-notes/#docker-desktop-471)|[v0.12.0](https://github.com/kubernetes-sigs/kind/releases/tag/v0.12.0)|v1.23.4|[v1.21.0](https://github.com/operator-framework/operator-sdk/releases/tag/v1.21.0)|
