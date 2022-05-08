# Cheatsheet - kubebuilder

|Command|Explanation|
|---|---|
|`kubebuilder init --domain <domain> --repo <repo>` |Initialize a project|
|`kubebuilder create api --group <group> --version <version> --kind <kind> --controller --resource`|Create new API resource|
|`kubebuilder create webhook --group <group> secret --version <version> --kind <kind> <option>`|Create new webhook. option: `--conversion`, `defaulting`, `--programmatic-validation`|
