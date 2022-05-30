# Run Operator

## Overview

We'll run [Prometheus](https://prometheus.io/) with and without [Prometheus Operator](https://github.com/prometheus-operator/prometheus-operator).

The Prometheus Operator serves to make running Prometheus on top of Kubernetes as easy as possible, while preserving Kubernetes-native configuration options.

## 1. Run Prometheus without operator

> Prometheus can reload its configuration at runtime. If the new configuration is not well-formed, the changes will not be applied. A configuration reload is triggered by sending a SIGHUP to the Prometheus process or sending a HTTP POST request to the /-/reload endpoint (when the --web.enable-lifecycle flag is enabled). This will also reload any configured rule files.

## 2. Run Prometheus with operator

1. Install operator ([bundle.yaml](https://github.com/prometheus-operator/prometheus-operator/blob/main/bundle.yaml))
    ```
    kubectl create -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/master/bundle.yaml
    ```

    <details><summary>Resources in bundle.yaml</summary>

    1. Custom Resource Definitions:
        1. `AlertmanagerConfig`
        1. `Alertmanager`
        1. `PodMonitor`
        1. `Probe`
        1. `Prometheus`
        1. `PrometheusRule`
        1. `ServiceMonitor`
        1. `ThanosRuler`
    1. `ClusterRoleBinding`: `prometheus-operator`
    1. `ClusterRole`: `prometheus-operator`
    1. `Deployment`: `prometheus-operator`
    1. `ServiceAccount`: `prometheus-operator`
    1. `Service`: `prometheus-operator`

    </details>

## 2. Run Grafana with operator

### 1. Install operator

```
kubectl apply -k github.com/grafana-operator/grafana-operator/deploy/manifests/
```

### 2. Create Grafana

Be sure to deploy in the same namespace as the operator.

```
kubectl apply -f https://raw.githubusercontent.com/grafana-operator/grafana-operator/master/deploy/examples/Grafana.yaml -n grafana-operator-system
```

Check status
```
kubectl get grafana example-grafana -n grafana-operator-system -o jsonpath='{.status}'
{"message":"success","phase":"reconciling","previousServiceName":"grafana-service"}
```

port forward

```
kubectl port-forward -n grafana-operator-system svc/grafana-service 3000
```

### 3. Dashboard

simple-dashboard (with raw json in `spec.json`)

```
kubectl apply -f https://raw.githubusercontent.com/grafana-operator/grafana-operator/master/deploy/examples/dashboards/SimpleDashboard.yaml -n grafana-operator-system
```

grafana-dashboard-from-grafana (with `spec.grafanaCom.id` and `spec.grafanaCom.revision`)

```
kubectl apply -f https://raw.githubusercontent.com/grafana-operator/grafana-operator/master/deploy/examples/dashboards/DashboardFromGrafana.yaml -n grafana-operator-system
```

keycloak-dashboard

```
kubectl apply -f https://raw.githubusercontent.com/grafana-operator/grafana-operator/master/deploy/examples/dashboards/KeycloakDashboard.yaml -n grafana-operator-system
```

## More operators

- https://github.com/argoproj/argo-cd
- https://github.com/strimzi/strimzi-kafka-operator
- https://github.com/mysql/mysql-operator
- https://github.com/zalando/postgres-operator
- https://github.com/aws-controllers-k8s/community
- https://github.com/grafana-operator/grafana-operator
