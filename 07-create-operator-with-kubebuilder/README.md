# Create Operator with Kubebuilder

## Version

[3.4.0](https://github.com/kubernetes-sigs/kubebuilder/releases/tag/v3.4.0) or later

You can check the tested versions in each example:

1. [QuickStart#Versions](02-quick-start/README.md#versions)
1. [PasswordOperator#Versions](03-password-operator/README.md#versions)

## Contents

1. [Install kubebuilder](01-install)
1. [QuickStart](02-quick-start)
1. [PasswordOperator](03-password-operator)

## Checklist

- [ ] How to create a project with kubebuilder?
    - [ ] What's inside `api/<versions>/` directory?
    - [ ] What's inside `internal/controller` directory (`controllers` directory for project created with kubebuilder version 3.9 or earlier)?
    - [ ] What's inside `config/` directory?
    - [ ] What's a `Manager` in `cmd/main.go` (`main.go` for project created with kubebuilder version 3.9 or earlier) ?
- [ ] How to add an API (e.g. `Password`)?
    - Resource
        - [ ] Where will the new type definition be stored?
        - [ ] What is `Password`, `PasswordSpec`, and `PasswordStatus` structs?
    - Controller
        - [ ] What's the input (arguments) and output (returned value) of `Reconcile` function of `Reconciler` interface?
- [ ] What's markers? (e.g. `+kubebuilder:object:root`)
    - If controller needs to manipulate Secret, which marker needs to be added?
- [ ] OwnerReferences
