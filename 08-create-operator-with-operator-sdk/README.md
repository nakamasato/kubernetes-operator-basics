# Create Operator with Operator SDK

Install Operator SDK CLI:

```
brew install operator-sdk
```

Or any other ways in [Installation](https://sdk.operatorframework.io/docs/installation/)

```
operator-sdk version
operator-sdk version: "v1.19.1", commit: "079d8852ce5b42aa5306a1e33f7ca725ec48d0e3", kubernetes version: "v1.23", go version: "go1.18.1", GOOS: "darwin", GOARCH: "amd64"
```

<details><summary>Check Commands</summary>

```
operator-sdk
CLI tool for building Kubernetes extensions and tools.

Usage:
  operator-sdk [flags]
  operator-sdk [command]

Examples:
The first step is to initialize your project:
    operator-sdk init [--plugins=<PLUGIN KEYS> [--project-version=<PROJECT VERSION>]]

<PLUGIN KEYS> is a comma-separated list of plugin keys from the following table
and <PROJECT VERSION> a supported project version for these plugins.

                                   Plugin keys | Supported project versions
-----------------------------------------------+----------------------------
           ansible.sdk.operatorframework.io/v1 |                          3
              declarative.go.kubebuilder.io/v1 |                       2, 3
                          go.kubebuilder.io/v2 |                       2, 3
                          go.kubebuilder.io/v3 |                          3
              helm.sdk.operatorframework.io/v1 |                          3
 hybrid.helm.sdk.operatorframework.io/v1-alpha |                          3
            kustomize.common.kubebuilder.io/v1 |                          3
           quarkus.javaoperatorsdk.io/v1-alpha |                          3

For more specific help for the init command of a certain plugins and project version
configuration please run:
    operator-sdk init --help --plugins=<PLUGIN KEYS> [--project-version=<PROJECT VERSION>]

Default plugin keys: "go.kubebuilder.io/v3"
Default project version: "3"


Available Commands:
  alpha            Alpha-stage subcommands
  bundle           Manage operator bundle metadata
  cleanup          Clean up an Operator deployed with the 'run' subcommand
  completion       Load completions for the specified shell
  create           Scaffold a Kubernetes API or webhook
  edit             Update the project configuration
  generate         Invokes a specific generator
  help             Help about any command
  init             Initialize a new project
  olm              Manage the Operator Lifecycle Manager installation in your cluster
  pkgman-to-bundle Migrates packagemanifests to bundles
  run              Run an Operator in a variety of environments
  scorecard        Runs scorecard
  version          Print the operator-sdk version

Flags:
  -h, --help                     help for operator-sdk
      --plugins strings          plugin keys to be used for this subcommand execution
      --project-version string   project version (default "3")
      --verbose                  Enable verbose logging

Use "operator-sdk [command] --help" for more information about a command.
```

</details>
