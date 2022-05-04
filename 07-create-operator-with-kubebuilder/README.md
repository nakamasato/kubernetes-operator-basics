# Create Operator with Kubebuilder

## Version

[3.3.0](https://github.com/kubernetes-sigs/kubebuilder/releases/tag/v3.3.0)

## Contents

1. [Install kubebuilder](01-install)
1. [QuickStart](02-quick-start)
1. [PasswordOperator](03-password-operator)

## Checklist

- [ ] How to create a project with kubebuilder?
    - [ ] What's inside `api/<versions>/` directory?
    - [ ] What's inside `controllers/` directory?
    - [ ] What's inside `config/` directory?
    - [ ] What's a `Manager` in `main.go`?
- [ ] How to add an API (e.g. `Password`)?
    - Resource
        - [ ] Where will the new type definition be stored?
        - [ ] What is `Password`, `PasswordSpec`, and `PasswordStatus` structs?
    - Controller
        - [ ] What's the input (arguments) and output (returned value) of `Reconcile` function of `Reconciler` interface?
- [ ] What's markers? (e.g. `+kubebuilder:object:root`)
    - If controller needs to manipulate Secret, which marker needs to be added?
- [ ] OwnerReferences
