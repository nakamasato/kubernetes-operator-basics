# Environment Setup

## 1. Docker

Install from https://docs.docker.com/desktop/mac/install/

## 2. Homebrew

```
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

## 3. kind

```
brew install kind
```

Ref: https://kind.sigs.k8s.io/docs/user/quick-start/#installation

## 4. kubectl

```
brew install kubectl
```

Ref: https://kubernetes.io/docs/tasks/tools/install-kubectl-macos/#install-with-homebrew-on-macos

## 5. Golang

Install from https://go.dev/doc/install

## Check versions

```
docker version
brew --version
kind --version
kubectl version --client --short
go version
```
