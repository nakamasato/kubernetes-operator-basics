# 1. Install kubebuilder CLI

```
# download kubebuilder and install locally.
curl -L -o kubebuilder https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)
chmod +x kubebuilder && mv kubebuilder /usr/local/bin/
```

https://book.kubebuilder.io/quick-start.html#installation

```
kubebuilder version
Version: main.version{KubeBuilderVersion:"3.10.0", KubernetesVendor:"1.26.1", GitCommit:"0fa57405d4a892efceec3c5a902f634277e30732", BuildDate:"2023-04-15T08:10:35Z", GoOs:"darwin", GoArch:"arm64"}
```

If you want to specify a kubebuilder version, you can use the following command:

```
KUBEBUILDER_VERSION=v3.10.0
curl -L -o kubebuilder https://github.com/kubernetes-sigs/kubebuilder/releases/download/$KUBEBUILDER_VERSION/kubebuilder_$(go env GOOS)_$(go env GOARCH)
chmod +x kubebuilder && mv kubebuilder /usr/local/bin/
```

<details><summary>Check Commands</summary>

```
CLI tool for building Kubernetes extensions and tools.

Usage:
  kubebuilder [flags]
  kubebuilder [command]

Examples:
The first step is to initialize your project:
    kubebuilder init [--plugins=<PLUGIN KEYS> [--project-version=<PROJECT VERSION>]]

<PLUGIN KEYS> is a comma-separated list of plugin keys from the following table
and <PROJECT VERSION> a supported project version for these plugins.

                             Plugin keys | Supported project versions
-----------------------------------------+----------------------------
               base.go.kubebuilder.io/v3 |                          3
               base.go.kubebuilder.io/v4 |                          3
        declarative.go.kubebuilder.io/v1 |                       2, 3
 deploy-image.go.kubebuilder.io/v1-alpha |                          3
                    go.kubebuilder.io/v2 |                       2, 3
                    go.kubebuilder.io/v3 |                          3
                    go.kubebuilder.io/v4 |                          3
         grafana.kubebuilder.io/v1-alpha |                          3
      kustomize.common.kubebuilder.io/v1 |                          3
      kustomize.common.kubebuilder.io/v2 |                          3

For more specific help for the init command of a certain plugins and project version
configuration please run:
    kubebuilder init --help --plugins=<PLUGIN KEYS> [--project-version=<PROJECT VERSION>]

Default plugin keys: "go.kubebuilder.io/v4"
Default project version: "3"


Available Commands:
  completion  Load completions for the specified shell
  create      Scaffold a Kubernetes API or webhook
  edit        Update the project configuration
  help        Help about any command
  init        Initialize a new project
  version     Print the kubebuilder version

Flags:
  -h, --help                     help for kubebuilder
      --plugins strings          plugin keys to be used for this subcommand execution
      --project-version string   project version (default "3")

Additional help topics:
  kubebuilder alpha      Alpha-stage subcommands

Use "kubebuilder [command] --help" for more information about a command.
```

</details>
