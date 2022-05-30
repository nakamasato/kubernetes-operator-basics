# PasswordOperator

## Version
- Go 1.17.9
- Kubebuilder: [v3.4.0](https://github.com/kubernetes-sigs/kubebuilder/releases/tag/v3.4.0) or later

## 0. Design Operator `PasswordOperator`

1. Make only one Custom Resource `Password`.
1. When custom resource `Password` is created, the controller generates a password, create a `Secret` with the same name as the `Password` object, and store the password in `password` field.
1. Provide password generation options. e.g. password length, the number of digits and symbols, etc.

![](01-design-operator.drawio.svg)

## 1. [kubebuilder] Init project

Make a directory and initialize git repository

```
mkdir -p ~/projects/password-operator
cd ~/projects/password-operator
git init
```

Initialize a project

```
kubebuilder init --domain example.com --repo example.com/password-operator
```

Commit
```
git add . && git commit -m "[kubebuilder] Init project"
```

## 2. [kubebuilder] Create API `Password` (Controller & Resource)
Create an API `password` (and choose `y` for resource and controller)

```
kubebuilder create api --group secret --version v1alpha1 --kind Password
```

<details><summary>Check if failed</summary>

If you're using `kubebuilder` version less than 3.4.0 and go version 1.18, you'll encounter the following error.

```
bash: /path/to/your/guestbook/bin/controller-gen: No such file or directory
make: *** [generate] Error 127
```

Fix `Makefile` -> Commit `git commit -m "Fix Makefile"`

```Makefile
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

Try again.

```
kubebuilder create api --group secret --version v1alpha1 --kind Password --force --controller --resource
```

</details>

Update CRD yaml files (Go types → CRD)

```
make manifests
```

Commit

```
git add . && git commit -m "[kubebuilder] Create API Password (Controller & Resource)"
```

### Column: About `kubebuilder` project development

![](development.drawio.svg)

There are three types of changes:
1. **kubebuilder** command create files (`init`, `create api`, `create webhook`, et.c)
1. Implement **API** (schema `apis/<version>/xxx_types.go`, `apis/<version>/xxx_webhook.go`)
1. Implement **controller** (reconcile loop `controllers/xxx_controller.go`)

If you're aware of which kind of change you're making, it'll be helpful to understand what exactly you're doing. I use `[kubebuilder]`, `[Controller]`, `[API]` as a prefix for each step title.

## 3. [Controller] Add log in Reconcile function

Update `controllers/password_controller.go`

```go
func (r *PasswordReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    logger := log.FromContext(ctx)

    logger.Info("Reconcile is called.")

    return ctrl.Result{}, nil
}
```

Install CRD

```
make install
```

Run
```
make run
```

Result:

```
go run ./main.go
1.651026659120943e+09   INFO    controller-runtime.metrics     Metrics server is starting to listen    {"addr": ":8080"}
1.65102665912137e+09    INFO    setup   starting manager
1.65102665912194e+09    INFO    Starting server {"path": "/metrics", "kind": "metrics", "addr": "[::]:8080"}
1.65102665912202e+09    INFO    Starting server {"kind": "health probe", "addr": "[::]:8081"}
1.651026659122195e+09   INFO    controller.password     Starting EventSource    {"reconciler group": "secret.example.com", "reconciler kind": "Password", "source": "kind source: *v1alpha1.Password"}
1.65102665912222e+09    INFO    controller.password     Starting Controller     {"reconciler group": "secret.example.com", "reconciler kind": "Password"}
1.6510266592234979e+09  INFO    controller.password     Starting workers        {"reconciler group": "secret.example.com", "reconciler kind": "Password", "worker count": 1}
```

Create `Password` object
```
kubectl apply -f config/samples
```

See logs

Reconcile function is called with `Password` with name `password-sample`
```
1.651026742035841e+09   INFO    controller.password     Reconcile is called.       {"reconciler group": "secret.example.com", "reconciler kind": "Password", "name": "password-sample", "namespace": "default"}
```

Let's see the arguments of `Reconcile`: `Reconcile(ctx context.Context, req ctrl.Request)`
- `context.Context`:
- `ctrl.Request`: <- Request that calls the Reconcile function. Get the object from the request.


Delete the object with
```
kubectl delete -f config/samples
```

See logs.
Reconcile function is called with `Password` with name `password-sample`

```
1.65102695302618e+09    INFO    controller.password     Reconcile is called.       {"reconciler group": "secret.example.com", "reconciler kind": "Password", "name": "password-sample", "namespace": "default"}
```

From the logs, the two events (creation and deletion) triggered the `Reconcile` function exactly the same way. We cannot distinguish them in `Reconcile`. (**Important**)

Stop the controller with `Ctrl+C`.
Commit
```
git add . && git commit -m "[Controller] Add log in Reconcile function"
```

**Point**: Reconcile function is called when custom resource object is created, updated, or deleted. Inside the Reconcile function, the reconciliation logic should not be dependent on the triggering type (`created`, `updated`, `deleted`).

## 4. [API] Remove Foo field from custom resource Password

By default, `Password` has `PasswordSpec` (with `Foo` field) and `PasswordStatus` (without any field):
```go
type PasswordSpec struct {
    // INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
    // Important: Run "make" to regenerate code after modifying this file
    // Foo is an example field of Password. Edit password_types.go to remove/update
    Foo string `json:"foo,omitempty"`
}
type PasswordStatus struct {
    // INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
    // Important: Run "make" to regenerate code after modifying this file
}
type Password struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`
    Spec   PasswordSpec   `json:"spec,omitempty"`
    Status PasswordStatus `json:"status,omitempty"`
}
```
- `TypeMeta`: API version and Kind (all Kubernetes objects have)
- `ObjectMeta`: name, namespace, labels, ... (all Kubernetes object have)
- `Spec`: Desired State
- `Status`: Actual State
- `+kubebuilder:object:root` comment is called a marker. -> telling controller-tools (our code and YAML generator) extra information.
    - `+kubebuilder:object:root`: tell the object generator that this type represents a Kind. → the object generator generates an implementation of the runtime.Object interface (Kinds must implement)
- add the Go types to the API group
    ```go
    func init() {
      SchemeBuilder.Register(&Password{}, &PasswordList{})
    }
    ```

When you create a new resource with `kubebuilder create api`, it automatically adds a field `Foo` in spec of the new resource. You can see it in the CRD.

```bash
kubectl get crd passwords.secret.example.com -o jsonpath='{.spec.versions[].schema.openAPIV3Schema.properties.spec}' | jq
{
    "description": "PasswordSpec defines the desired state of Password",
    "properties": {
    "foo": {
        "description": "Foo is an example field of Password. Edit password_types.go to remove/update",
        "type": "string"
    }
    },
    "type": "object"
}
```

Let's remove `Foo` field from `api/v1alpha1/password_types.go` and run `make manifests` to update the CRD yaml files `config/crd/bases/secret.example.com_passwords.yaml`.

We also need to update the CRD registered in `api-server` as `Foo` is already removed:

```bash
make install
```
Now you can confirm the field `Foo` is removed.

```bash
kubectl get crd passwords.secret.example.com -o jsonpath='{.spec.versions[].schema.openAPIV3Schema.properties.spec}' | jq
{
    "description": "PasswordSpec defines the desired state of Password",
    "type": "object"
}
```

Commit
```
git commit -am "[API] Remove Foo field from custom resource Password"
```

**Point**: When updating API resource:
![](development-api.drawio.svg)
1. update in `api/<version>/<custom_resource>_types.go`
1. `make install`
    1. `make manifests`: Generate CRD `config/crd/bases/<custom_resource>.<domain>_<custom_resource>.yaml`
    1. `$(KUSTOMIZE) build config/crd | kubectl apply -f -`: Apply crd yaml file.

Other API files:
- `groupversion_info.go`:
- `zz_generated.deepcopy.go`:

## About Controller
In [controller-runtime](https://pkg.go.dev/sigs.k8s.io/controller-runtime), the logic that implements the reconciling for a specific kind is called a **Reconciler**.

<details><summary>what's controller-runtime?</summary>

We studied in [operator development method](../../06-operator-development-method/05-methods/)

![](../../06-operator-development-method/05-methods/comparison.drawio.svg)

</details>


A reconciler takes the name of an object, and returns whether or not we need to try again. (err -> reconcile again later, no error -> reconciliation completed.)

[Reconciler](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile#Reconciler):

```go
type Reconciler interface {
    // Reconcile performs a full reconciliation for the object referred to by the Request.
    // The Controller will requeue the Request to be processed again if an error is non-nil or
    // Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
    Reconcile(context.Context, Request) (Result, error)
}
```

[Request](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile#Request):

```go
type Request struct {
    // NamespacedName is the name and namespace of the object to reconcile.
    types.NamespacedName
}
```

[Result](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile#Result):

```go
type Result struct {
    // Requeue tells the Controller to requeue the reconcile key.  Defaults to false.
    Requeue bool

    // RequeueAfter if greater than 0, tells the Controller to requeue the reconcile key after the Duration.
    // Implies that Requeue is true, there is no need to set Requeue to true at the same time as RequeueAfter.
    RequeueAfter time.Duration
}
```

- `PasswordReconciler` with `client.Client`.
- RBAC markers for autogeneration of rbac yaml.
- Request just has a name.
- Register `PasswordReconciler` to Manager.

## 5. [Controller] Fetch Password object

Add the following lines to `Reconcile function`
```go
// Fetch Password object
var password secretv1alpha1.Password
if err := r.Get(ctx, req.NamespacedName, &password); err != nil {
    logger.Error(err, "Fetch Password object - failed")
    return ctrl.Result{}, client.IgnoreNotFound(err)
}

logger.Info("Fetch Password object - succeeded", "password", password.Name, "createdAt", password.CreationTimestamp)
```

run!

```
make run
```

```
kubectl apply -f config/samples
```

```
1.651098875458412e+09   INFO    controller.password     Fetch Password object - succeeded   {"reconciler group": "secret.example.com", "reconciler kind": "Password", "name": "password-sample", "namespace": "default", "password": "password-sample", "createdAt": "2022-04-28 07:34:35 +0900 JST"}
```

```
kubectl delete -f config/samples
```

```
1.651102065284576e+09   ERROR   controller.password     Fetch Password object - failed      {"reconciler group": "secret.example.com", "reconciler kind": "Password", "name": "password-sample", "namespace": "default", "error": "Password.secret.example.com \"password-sample\" not found"}
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).Reconcile
        /Users/nakamasato/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.11.0/pkg/internal/controller/controller.go:114
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).reconcileHandler
        /Users/nakamasato/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.11.0/pkg/internal/controller/controller.go:311
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).processNextWorkItem
        /Users/nakamasato/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.11.0/pkg/internal/controller/controller.go:266
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).Start.func2.2
        /Users/nakamasato/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.11.0/pkg/internal/controller/controller.go:227
```

Commit

```
git commit -am "[Controller] Fetch Password object"
```

## 6. [Controller] Create Secret object if not exists

Logic:
1. Try to fetch a `Secret` with the same name as `Password` object.
1. Return if already exists. Otherwise, create a `Secret`.

Secret: https://pkg.go.dev/k8s.io/api/core/v1#Secret

```go
import (
    // existing packages...
    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/errors"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
```

```go
    // Create Secret object if not exists
    var secret corev1.Secret
    if err := r.Get(ctx, req.NamespacedName, &secret); err != nil {
        if errors.IsNotFound(err) {
            // Create Secret
            logger.Info("Create Secret object if not exists - create secret")
            secret := newSecretFromPassword(&password)
            err = r.Create(ctx, secret)
            if err != nil {
                logger.Error(err, "Create Secret object if not exists - failed to create Secret")
                return ctrl.Result{}, err
            }
            logger.Info("Create Secret object if not exists - Secret successfully created")
        } else {
            logger.Error(err, "Create Secret object if not exists - failed to fetch Secret")
            return ctrl.Result{}, err
        }
    }

    logger.Info("Create Secret object if not exists - completed")
```

```go
func newSecretFromPassword(password *secretv1alpha1.Password) *corev1.Secret {
    secret := &corev1.Secret{
        ObjectMeta: metav1.ObjectMeta{
            Name:      password.Name,
            Namespace: password.Namespace,
        },
        Data: map[string][]byte{
            "password": []byte("123456789"), // password=123456789
        },
    }
    return secret
}
```

Run!

```
make run
```

```
kubectl apply -f config/samples
```

```
1.6511075442946272e+09  INFO    controller.password     Reconcile is called.    {"reconciler group": "secret.example.com", "reconciler kind": "Password", "name": "password-sample", "namespace": "default"}
1.651107544294673e+09   INFO    controller.password     Fetch Password object - succeeded    {"reconciler group": "secret.example.com", "reconciler kind": "Password", "name": "password-sample", "namespace": "default", "password": "password-sample", "createdAt": "2022-04-28 09:59:04 +0900 JST"}
1.6511075444954848e+09  INFO    controller.password     Create Secret object if not exists - create secret       {"reconciler group": "secret.example.com", "reconciler kind": "Password", "name": "password-sample", "namespace": "default"}
1.651107544582052e+09   INFO    controller.password     Create Secret object if not exists - Secret successfully created {"reconciler group": "secret.example.com", "reconciler kind": "Password", "name": "password-sample", "namespace": "default"}
1.651107544582082e+09   INFO    controller.password     Create Secret object if not exists - completed   {"reconciler group": "secret.example.com", "reconciler kind": "Password", "name": "password-sample", "namespace": "default"}
```

```
kubectl get secret password-sample
NAME              TYPE     DATA   AGE
password-sample   Opaque   1      47s
```

```
kubectl get secret password-sample -o jsonpath='{.data}'
{"password":"MTIzNDU2Nzg5"}%
echo -n MTIzNDU2Nzg5 | base64 --decode
123456789%
```

```
kubectl delete -f config/samples
```

```
1.651107652362049e+09   INFO    controller.password     Reconcile is called.    {"reconciler group": "secret.example.com", "reconciler kind": "Password", "name": "password-sample", "namespace": "default"}
1.651107652362087e+09   ERROR   controller.password     Fetch Password object - failed       {"reconciler group": "secret.example.com", "reconciler kind": "Password", "name": "password-sample", "namespace": "default", "error": "Password.secret.example.com \"password-sample\" not found"}
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).Reconcile
        /Users/nakamasato/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.11.0/pkg/internal/controller/controller.go:114
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).reconcileHandler
        /Users/nakamasato/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.11.0/pkg/internal/controller/controller.go:311
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).processNextWorkItem
        /Users/nakamasato/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.11.0/pkg/internal/controller/controller.go:266
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).Start.func2.2
        /Users/nakamasato/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.11.0/pkg/internal/controller/controller.go:227
```

```
kubectl get password
No resources found in default namespace.
```

```
kubectl get secret password-sample
NAME              TYPE     DATA   AGE
password-sample   Opaque   1      2m25s
```

manually delete secret at this point.

```
kubectl delete secret password-sample
```

Let's run the operator on kubernetes cluster!

```
IMG=password-operator:v1
make docker-build IMG=$IMG
kind load docker-image $IMG
make deploy IMG=$IMG
```

Check pod

```
kubectl get po -n password-operator-system
NAME                                                    READY   STATUS    RESTARTS   AGE
password-operator-controller-manager-796b9d99b6-x7qnn   2/2     Running   0          24s
```

Check logs

```
kubectl logs password-operator-controller-manager-796b9d99b6-x7qnn -n password-operator-system -f
```

```
kubectl apply -f config/samples
```

```
E0428 01:52:38.175924       1 reflector.go:138] pkg/mod/k8s.io/client-go@v0.23.0/tools/cache/reflector.go:167: Failed to watch *v1.Secret: failed to list *v1.Secret: secrets is forbidden: User "system:serviceaccount:password-operator-system:password-operator-controller-manager" cannot list resource "secrets" in API group "" at the cluster scope
```

The permission granted to the service account used by the controller-manager is not enough!

**Important point!** This controller needs to manipulate Secret! Let's grant the permissions by the following RBAC markers!

```go
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;
```

Need to run `make manifests` to update `ClusterRole` for the controller.

```
make manifests
```

You can see the following lines are added to `config/rbac/role.yaml`.

```yaml
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - create
  - get
  - list
  - watch
```

```
make deploy IMG=$IMG
```

```
kubectl get secret
NAME                  TYPE                                  DATA   AGE
default-token-nppdh   kubernetes.io/service-account-token   3      26h
password-sample       Opaque                                1      11s
```

Now you can see a new secret is created by the operator!

```
make undeploy
```

```
kubectl delete secret password-sample
```

Next: Clean up the orphaned secret!

Commit

```
git add . && git commit -m "[Controller] Create Secret object if not exists"
```

## 7. [Controller] Clean up Secret when Password is deleted

Logic:
1. If Password is deleted, the corresponding Secret is also deleted.

Add the following lines to `Reconcile` function just after `secret := newSecretFromPassword(&password)`

```go
    err := ctrl.SetControllerReference(&password, secret, r.Scheme) // Set owner of this Secret
    if err != nil {
        logger.Error(err, "Create Secret object if not exists - failed to set SetControllerReference")
        return ctrl.Result{}, err
    }
```

[SetControllerReference](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/controller/controllerutil#SetControllerReference)

> SetControllerReference sets owner as a Controller OwnerReference on controlled. This is used for **garbage collection of the controlled object** and for reconciling the owner object on changes to controlled (with a Watch + EnqueueRequestForOwner)

```
make install run
```

```
kubectl apply -f config/samples
```

```
kubectl get secret
NAME                  TYPE                                  DATA   AGE
default-token-nppdh   kubernetes.io/service-account-token   3      26h
password-sample       Opaque                                1      24s
```

```
kubectl delete -f config/samples
```

Check Secret

```
kubectl get secret
NAME                  TYPE                                  DATA   AGE
default-token-nppdh   kubernetes.io/service-account-token   3      26h
```

Secret `password-sample` is deleted!

Commit
```
git commit -am "[Controller] Clean up Secret when Password is deleted"
```

## 8. [Controller] Generate random password

Use for password generation: https://github.com/sethvargo/go-password

Import `github.com/sethvargo/go-password/password`

```go
import (
    passwordGenerator "github.com/sethvargo/go-password/password"
    ...
)
```

remove:
```diff
- secret := newSecretFromPassword(&password)
- err := ctrl.SetControllerReference(&password, secret, r.Scheme) // Set owner of this Secret
```

new:
```go
    passwordStr, err := passwordGenerator.Generate(64, 10, 10, false, false)
    if err != nil {
        logger.Error(err, "Create Secret object if not exists - failed to generate password")
        return ctrl.Result{}, err
    }
    secret := newSecretFromPassword(&password, passwordStr)
    err = ctrl.SetControllerReference(&password, secret, r.Scheme) // Set owner of this Secret
```

Update `newSecretFromPassword` to pass `passwordStr` as an argument:

```go
func newSecretFromPassword(password *secretv1alpha1.Password, passwordStr string) *corev1.Secret {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      password.Name,
			Namespace: password.Namespace,
		},
		Data: map[string][]byte{
			"password": []byte(passwordStr),
		},
	}
	return secret
}
```

```
make run
```

```
kubectl apply -f config/samples
```

Check password

```
kubectl get secret password-sample -o jsonpath='{.data.password}' |
 base64 --decode
yh0<B-?qfOkolF#EKa>D5Ig924rZQxiU_dneAX86S1YsLR@TjvW}u\3mM7]NHGVz%
```

Delete & create a new one

```
kubectl delete -f config/samples
password.secret.example.com "password-sample" deleted
kubectl apply -f config/samples
password.secret.example.com/password-sample created
kubectl get secret password-sample -o jsonpath='{.data.password}' | base64 --decode
noY$Xa9KI3At(J+bwvLdqi4hDB/CT~ZxGfpR[7elWrS5Ocz=VMym)u#2F1_60jN8%
```

Confirmed password is randomly generated.

Commit

```
git commit -am "[Controller] Generate random password"
```

## 9. [API&Controller] Make password configurable with CRD fields

When generating a password in this line,

```go
passwordStr, err := password.Generate(64, 10, 10, false, false)
```

`64, 10, 10, false, false` are hard-coded in the Reconcile function.

Each argument represents:
1. password length
1. the number of digits
1. the number of symbols
1. allow upper and lower case letters if true
1. disallow repeat characters if true

Let's enable to configure them with our custom resource `Password`.

`api/v1alpha1/password_types.go`:

```go
type PasswordSpec struct {
    //+kubebuilder:validation:Minimum=8
    //+kubebuilder:default:=20
    //+kubebuilder:validation:Required
    Length int `json:"length"`

    //+kubebuilder:validation:Minimum=0
    //+kubebuilder:default:=10
    //+kubebuilder:validation:Optional
    Digit int `json:"digit"`

    //+kubebuilder:validation:Minimum=0
    //+kubebuilder:default:=10
    //+kubebuilder:validation:Optional
    Symbol int `json:"symbol"`

    //+kubebuilder:default:=false
    //+kubebuilder:validation:Optional
    CaseSensitive  bool `json:"caseSensitive"`
    //+kubebuilder:default:=false
    //+kubebuilder:validation:Optional
    DisallowRepeat bool `json:"disallowRepeat"`
}
```

Update `config/crd/bases/secret.example.com_passwords.yaml` and apply CRD:

```
make install
```

Apply `Password`:

```yaml
apiVersion: secret.example.com/v1alpha1
kind: Password
metadata:
  name: password-sample
spec:
  length: 20
```

```
kubectl apply -f config/samples/
```

<details>

```
kubectl get -f config/samples/ -o yaml
apiVersion: secret.example.com/v1alpha1
kind: Password
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"secret.example.com/v1alpha1","kind":"Password","metadata":{"annotations":{},"name":"password-sample","namespace":"default"},"spec":{"length":20}}
  creationTimestamp: "2022-05-02T00:20:46Z"
  generation: 2
  name: password-sample
  namespace: default
  resourceVersion: "61061"
  uid: 28eab3e8-bc1a-4a36-9fa9-26684bb40892
spec:
  caseSensitive: false
  digit: 10
  disallowRepeat: false
  length: 20
  symbol: 10
```

</details>

You can see the spec:
```yaml
spec:
  caseSensitive: false
  digit: 10
  disallowRepeat: false
  length: 20
  symbol: 10
```

Use these values in the controller (`controllers/password_controller.go`).

```go
            passwordStr, err := passwordGenerator.Generate(
                password.Spec.Length,
                password.Spec.Digit,
                password.Spec.Symbol,
                password.Spec.CaseSensitive,
                password.Spec.DisallowRepeat,
            )
```

```
make run
```

Recreate the custom resource `Password`.
```
kubectl delete -f config/samples/ && kubectl apply -f config/samples/
```

Check the length of the password of the generated Secret:

```
kubectl get secret password-sample -o jsonpath='{.data.password}' | base64 --decode
,<!463[1$58#90_7>2)~
kubectl get secret password-sample -o jsonpath='{.data.password}' | base64 --decode | wc -m
      20
```

Change the length to 30 in `config/samples/secret_v1alpha1_password.yaml` and recreate the custom resource `Password`:

```
kubectl delete -f config/samples/ && kubectl apply -f config/samples/
```

```
kubectl get secret password-sample -o jsonpath='{.data.password}' | base64 --decode
?f*%76X0:/41Y.3V8$2=r>q9m5I{ax%
kubectl get secret password-sample -o jsonpath='{.data.password}' | base64 --decode | wc -m
      30
```

Change the length to 10 and recreate the custom resource `Password`:

```
kubectl delete -f config/samples/ && kubectl apply -f config/samples/
```

You'll see the following `Error` log: `number of digits and symbols must be less than total length`

```
1.651452110121384e+09   ERROR   controller.password     Create Secret object if not exists - failed to generate password {"reconciler group": "secret.example.com", "reconciler kind": "Password", "name": "password-sample", "namespace": "default", "error": "number of digits and symbols must be less than total length"}
```

The reconcilation loop failed to generate password. `Secret` was not successfully generated.

```
kubectl get secret password-sample -o jsonpath='{.data.password}'
Error from server (NotFound): secrets "password-sample" not found
```

Commit

```
git commit -am "[API&Controller] Make password configurable with CRD fields"
```

Next:
- [ ] Update `Password`'s status to tell if custom resource is successfully updated.
- [ ] Add validation for `digit`, `symbol`, and `length`: `number of digits and symbols must be less than total length`

## 10. [API&Controller] Add Password Status

Update Go types:

Create new type `PasswordState`

```go
type PasswordState string

const (
	PasswordInSync  PasswordState = "InSync"
	PasswordFailed  PasswordState = "Failed"
)
```

Add `State` to `PasswordStatus`

```go
// PasswordStatus defines the observed state of Password
type PasswordStatus struct {

    // Information about if Password is in-sync.
    State PasswordState `json:"state,omitempty"` // in-sync, failed
}
```

Update CRD yaml files:

```
make manifests
```

Add the following lines at the end of `Reconcile` function in the controller:

```go
func (r *PasswordReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    // ...
    password.Status.State = secretv1alpha1.PasswordInSync
    if err := r.Status().Update(ctx, &password); err != nil {
        logger.Error(err, "Failed to update Password status")
        return ctrl.Result{}, err
    }
    return ctrl.Result{}, nil
}
```

Run (Apply new CRD and run the controller):

```
make install run
```

```yaml
apiVersion: secret.example.com/v1alpha1
kind: Password
metadata:
  name: password-sample
spec:
  length: 20
```

```
kubectl delete -f config/samples/ && kubectl apply -f config/samples/
```

```
kubectl get password password-sample -o jsonpath='{.status}'
{"state":"InSync"}
```

In the same way, let's add status for `Failed`. Add the following lines to where to return error in reconcile function.

```go
password.Status.State = secretv1alpha1.PasswordFailed
if err := r.Status().Update(ctx, &password); err != nil {
    logger.Error(err, "Failed to update Password status")
    return ctrl.Result{}, err
}
```

Run

```
make run
```

```yaml
apiVersion: secret.example.com/v1alpha1
kind: Password
metadata:
  name: password-sample
spec:
  length: 10 # this will cause 'number of digits and symbols must be less than total length' error
```

```
kubectl delete -f config/samples/ && kubectl apply -f config/samples/
```

```
kubectl get password password-sample -o jsonpath='{.status}'
{"state":"Failed"}
```

Commit
```
git commit -am "[API&Controller] Add Password Status"
```

Homework: Add `Reason` to `PasswordStatus` and store the failure reason for `Failed` state.
1. Add new field to `PasswordStatus` (Go types).
1. Regenerate CRD manifests.
1. Add a logic to update it (controller).
1. Run controller and apply custom resource with failure condition.

## 11. [API] Add AdditionalPrinterColumns

We cannot see `State` with `kubectl get` now. Let's make it visible!

```
kubectl get password password-sample
NAME              AGE
password-sample   6m34s
```

Add a marker to `api/v1alpha1/password_types.go`:

```diff
 //+kubebuilder:object:root=true
 //+kubebuilder:subresource:status
+//+kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.state`

 // Password is the Schema for the passwords API
 type Password struct {
 	metav1.TypeMeta   `json:",inline"`
 	metav1.ObjectMeta `json:"metadata,omitempty"`

 	Spec   PasswordSpec   `json:"spec,omitempty"`
 	Status PasswordStatus `json:"status,omitempty"`
 }
```

Update the CRD yaml file by the command:
```
make manifests
```

With this marker, controller-gen will update CRD by adding `additionalPrinterColumns`:

```diff
+++ b/config/crd/bases/secret.example.com_passwords.yaml
@@ -15,7 +15,11 @@ spec:
     singular: password
   scope: Namespaced
   versions:
-  - name: v1alpha1
+  - additionalPrinterColumns:
+    - jsonPath: .status.state
+      name: State
+      type: string
+    name: v1alpha1
```

```
make install run
```

```
kubectl delete -f config/samples/ && kubectl apply -f config/samples/
```

```
kubectl get password password-sample
NAME              STATE
password-sample   InSync
```

Why `AGE` column is gone after setting `additionalPrinterColumns`? -> `AGE` is added if there's no `additionalPrinterColumn` specified: [apiserver/helpers.go#L81-L88](https://github.com/kubernetes/kubernetes/blob/e7a2ce75e5df96ba6ea51d904bf2735397b3e203/staging/src/k8s.io/apiextensions-apiserver/pkg/apiserver/helpers.go#L81-L88). If you want to get `AGE` back, you can add the following line:

```go
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
```

```diff
-  - name: v1alpha1
+  - additionalPrinterColumns:
+    - jsonPath: .metadata.creationTimestamp
+      name: Age
+      type: date
+    - jsonPath: .status.state
+      name: State
+      type: string
+    name: v1alpha1
```

```
NAME              AGE   STATE
password-sample   18m   InSync
```

Commit

```
git commit -am "[API] Add AdditionalPrinterColumns"
```

## 11. [kubebuilder] Create validating admission webhook

In this section, we'll implement a validation for `digit`, `symbol`, and `length`: `number of digits and symbols must be less than total length`

### 11.1. Admission Webhook in general

Let's start with **Admission Webhook**.

> Admission webhooks are HTTP callbacks that receive admission requests and do something with them.

For more details, you can read the official documentation:

1. [validating admission webhook](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#validatingadmissionwebhook):
1. [mutating admission webhook](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#mutatingadmissionwebhook):

Admission Controller: a piece of code that intercepts requests to the Kubernetes API server prior to persistence of the object, but after the request is authenticated and authorized.

![](https://raw.githubusercontent.com/nakamasato/kubernetes-training/main/contents/kubernetes-operator/admission-webhook.drawio.svg)

### 11.2. Admission Webhook with kubebuilder

`kubebuilder` makes it much simpler and easier to create a webhook by automating the following steps:
1. Create a webhook with `kubebuilder` command.
1. Add the webhook server to the manager.
1. Create handlers for the webhook.
1. Register each handler with a path in your server.

All you need to do is to implement the [Defaulter](https://github.com/kubernetes-sigs/controller-runtime/blob/3cb67228604239d3cd764b41719565bb4a46add6/pkg/webhook/admission/defaulter.go#L27-L31) and the [Validator](https://github.com/kubernetes-sigs/controller-runtime/blob/3cb67228604239d3cd764b41719565bb4a46add6/pkg/webhook/admission/validator.go#L29-L35) interface.

Create a webhook with `kubebuilder` command:

<details><summary>kubebuilder</summary>

```
kubebuilder create webhook --help
Scaffold a webhook for an API resource. You can choose to scaffold defaulting,
validating and/or conversion webhooks.

Usage:
  kubebuilder create webhook [flags]

Examples:
  # Create defaulting and validating webhooks for Group: ship, Version: v1beta1
  # and Kind: Frigate
  kubebuilder create webhook --group ship --version v1beta1 --kind Frigate --defaulting --programmatic-validation

  # Create conversion webhook for Group: ship, Version: v1beta1
  # and Kind: Frigate
  kubebuilder create webhook --group ship --version v1beta1 --kind Frigate --conversion


Flags:
      --conversion                if set, scaffold the conversion webhook
      --defaulting                if set, scaffold the defaulting webhook
      --force                     attempt to create resource even if it already exists
      --group string              resource Group
  -h, --help                      help for webhook
      --kind string               resource Kind
      --plural string             resource irregular plural form
      --programmatic-validation   if set, scaffold the validating webhook
      --version string            resource Version

Global Flags:
      --plugins strings   plugin keys to be used for this subcommand execution
```

</details>

### 11.3. Admission Webhook in our case
To validate our custom resource object, we'll use **Validating Admission Webhook**.

```
kubebuilder create webhook --group secret --version v1alpha1 --kind Password --programmatic-validation
```

The command above automatically adds the following lines to `main.go`, which adds the webhook server to the manager.

```diff
+       if err = (&secretv1alpha1.Password{}).SetupWebhookWithManager(mgr); err != nil {
+               setupLog.Error(err, "unable to create webhook", "webhook", "Password")
+               os.Exit(1)
+       }
```

Generated files:
1. `api`: implementation of webhook
    <details>

    1. `api/v1alpha1/password_webhook.go`: main logic of validating and defaulting.
    1. `api/v1alpha1/webhook_suite_test.go`: test for webhook

    </details>
1. `config`:
    1. certmanager for certificate used by webhook
        <details>

        1. `config/certmanager/certificate.yaml`
        1. `config/certmanager/kustomization.yaml`
        1. `config/certmanager/kustomizeconfig.yaml`

        </details>
    1. expose port and create service for webhook
        <details>

        1. `config/default/manager_webhook_patch.yaml`: Patch `controller-manager` to expose a port for webhook and mount certificate.
        1. `config/default/webhookcainjection_patch.yaml`: Patch `MutatingWebhookConfiguration` and `ValidatingWebhookConfiguration` to add annotations.

        </details>
    1. definition of webhook

        <details>

        1. `config/webhook/kustomization.yaml`
        1. `config/webhook/kustomizeconfig.yaml`: reference service from webhook config
        1. `config/webhook/manifests.yaml`: `MutatingWebhookConfiguration` and `ValidatingWebhookConfiguration`
        1. `config/webhook/service.yaml`: service for webhook

        </details>

1. `main.go`: Register the webhook server to the manager.

Let's start writing a validation logic in `api/v1alpha1/password_webhook.go`.

There are already four functions:
1. `SetupWebhookWithManager`: This is used in `main.go` to register the webhook.
    `main.go`:
    ```go
	if err = (&secretv1alpha1.Password{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "Password")
		os.Exit(1)
	}
    ```
1. `ValidateCreate`: Validation logic for `CREATE`
1. `ValidateUpdate`: Validation logic for `UPDATE`
1. `ValidateDelete`: Validation logic for `DELETE`

Commit:

```
git add . && git commit -am "[kubebuilder] Create validating admission webhook"
```

## 13. [API] Implement Validating Admission Webhook

Implement a common validate function `validatePassword` and use it in `ValidateCreate` and `ValidateUpdate`. We just leave `ValidateDelete` as it is as we don't need to validate on deletion.

```go
// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Password) ValidateCreate() error {
	passwordlog.Info("validate create", "name", r.Name)

	return r.validatePassword()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Password) ValidateUpdate(old runtime.Object) error {
	passwordlog.Info("validate update", "name", r.Name)

	return r.validatePassword()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Password) ValidateDelete() error {
	passwordlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

var ErrSumOfDigitAndSymbolMustBeLessThanLength = errors.New("Number of digits and symbols must be less than total length")

func (r *Password) validatePassword() error {
	if r.Spec.Digit+r.Spec.Symbol >= r.Spec.Length {
		return ErrSumOfDigitAndSymbolMustBeLessThanLength
	}
	return nil
}
```


Run! (We need to run the controller with cert manager because webhook requires TLS)

Install [Cert Manager](https://github.com/cert-manager/cert-manager)

```
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.8.0/cert-manager.yaml
```

Uncomment `[WEBHOOK]` and `[CERTMANAGER]` sections in `config/default/kustomization.yaml` and `config/crd/kustomization.yaml`

Comment out `mutating`:

<details><summary>config/webhook/kustomizeconfig.yaml</summary>

```yaml
# the following config is for teaching kustomize where to look at when substituting vars.
# It requires kustomize v2.1.0 or newer to work properly.
nameReference:
- kind: Service
  version: v1
  fieldSpecs:
  # - kind: MutatingWebhookConfiguration
  #   group: admissionregistration.k8s.io
  #   path: webhooks/clientConfig/service/name
  - kind: ValidatingWebhookConfiguration
    group: admissionregistration.k8s.io
    path: webhooks/clientConfig/service/name

namespace:
# - kind: MutatingWebhookConfiguration
#   group: admissionregistration.k8s.io
#   path: webhooks/clientConfig/service/namespace
#   create: true
- kind: ValidatingWebhookConfiguration
  group: admissionregistration.k8s.io
  path: webhooks/clientConfig/service/namespace
  create: true

varReference:
- path: metadata/annotations
```

</details>

<details><summary>config/default/webhookcainjection_patch.yaml</summary>

```yaml
# This patch add annotation to admission webhook config and
# the variables $(CERTIFICATE_NAMESPACE) and $(CERTIFICATE_NAME) will be substituted by kustomize.
# apiVersion: admissionregistration.k8s.io/v1
# kind: MutatingWebhookConfiguration
# metadata:
#   name: mutating-webhook-configuration
#   annotations:
#     cert-manager.io/inject-ca-from: $(CERTIFICATE_NAMESPACE)/$(CERTIFICATE_NAME)
# ---
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  name: validating-webhook-configuration
  annotations:
    cert-manager.io/inject-ca-from: $(CERTIFICATE_NAMESPACE)/$(CERTIFICATE_NAME)
```

</details>

```
make install
```

```
IMG=password-operator:webhook
make docker-build IMG=$IMG
kind load docker-image $IMG
make deploy IMG=$IMG
```

```
kubectl get po -n password-operator-system
NAME                                                    READY   STATUS    RESTARTS   AGE
password-operator-controller-manager-5cf44d87cf-drxfq   2/2     Running   0          2m33s
```

Apply invalid `Password`:

```yaml
apiVersion: secret.example.com/v1alpha1
kind: Password
metadata:
  name: password-sample
spec:
  length: 10
```

```
kubectl apply -f config/samples/secret_v1alpha1_password.yaml
Error from server (Number of digits and symbols must be less than total length): error when creating "config/samples/secret_v1alpha1_password.yaml": admission webhook "vpassword.kb.io" denied the request: Number of digits and symbols must be less than total length
```

**Validating Admission Webhook** works!

Clean up

```
make undeploy
kubectl delete -f https://github.com/cert-manager/cert-manager/releases/download/v1.8.0/cert-manager.yaml
```

Commit: `git add . && git commit -am "[API] Implement validating admission webhook"`

## Wrap Up

1. `[kubebuilder]` Init project
1. `[kubebuilder]` Create API Password (Controller & Resource)
1. Implement controller
    1. `[Controller]` Add log in Reconcile function
    1. `[API]` Remove Foo field from custom resource Password
    1. `[Controller]` Fetch Password object
    1. `[Controller]` Create Secret object if not exists
    1. `[Controller]` Clean up Secret when Password is deleted
    1. `[Controller]` Generate random password
1. Design API
    1. `[API&Controller]` Make password configurable with CRD fields
    1. `[API&Controller]` Add Password Status
    1. `[API]` Add AdditionalPrinterColumns
1. Webhook
    1. `[kubebuilder]` Create validating admission webhook
    1. `[API]` Implement validating admission webhook

## Versions

Checked version combinations:

|Docker|kind|kubernetes|kubebuilder|
|---|-----|---|---|
|[4.7.0 (77141)](https://docs.docker.com/desktop/mac/release-notes/#docker-desktop-471)|[v0.12.0](https://github.com/kubernetes-sigs/kind/releases/tag/v0.12.0)|v1.23.4|[v3.4.0](https://github.com/kubernetes-sigs/kubebuilder/releases/tag/v3.4.0)|
|[4.7.0 (77141)](https://docs.docker.com/desktop/mac/release-notes/#docker-desktop-471)|[v0.12.0](https://github.com/kubernetes-sigs/kind/releases/tag/v0.12.0)|v1.23.4|[v3.4.1](https://github.com/kubernetes-sigs/kubebuilder/releases/tag/v3.4.1)|

## References

1. Groups, Kinds, Versions https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md & https://book.kubebuilder.io/cronjob-tutorial/gvks.html
