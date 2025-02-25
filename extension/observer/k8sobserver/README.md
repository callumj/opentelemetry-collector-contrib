# Kubernetes Observer

<!-- status autogenerated section -->
| Status        |           |
| ------------- |-----------|
| Stability     | [alpha]  |
| Distributions | [contrib], [splunk] |

[alpha]: https://github.com/open-telemetry/opentelemetry-collector#alpha
[contrib]: https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib
[splunk]: https://github.com/signalfx/splunk-otel-collector
<!-- end autogenerated section -->

The `k8s_observer` is a [Receiver Creator](../../../receiver/receivercreator/README.md)-compatible "watch observer" that will detect and report
Kubernetes pod, port, and node endpoints via the Kubernetes API.

## Example Config

```yaml
extensions:
  k8s_observer:
    auth_type: serviceAccount
    node: ${env:K8S_NODE_NAME}
    observe_pods: true
    observe_nodes: true

receivers:
  receiver_creator:
    watch_observers: [k8s_observer]
    receivers:
      redis:
        rule: type == "port" && pod.name matches "redis"
        config:
          password: '`pod.labels["SECRET"]`'
      kubeletstats:
        rule: type == "k8s.node"
        config:
          auth_type: serviceAccount
          collection_interval: 10s
          endpoint: "`endpoint`:`kubelet_endpoint_port`"
          extra_metadata_labels:
            - container.id
          metric_groups:
            - container
            - pod
            - node
```

The `node` field can be set to the node name to limit discovered endpoints. For example, its name value can be obtained using the downward API inside a Collector pod spec as follows:

```yaml
env:
  - name: K8S_NODE_NAME
    valueFrom:
      fieldRef:
        fieldPath: spec.nodeName
```

This spec-determined value would then be available via the `${env:K8S_NODE_NAME}` usage in the observer configuration.

## Config

All fields are optional.

| Name | Type | Default | Docs |
| ---- | ---- | ------- | ---- |
| auth_type | string | `serviceAccount` | How to authenticate to the K8s API server.  This can be one of `none` (for no auth), `serviceAccount` (to use the standard service account token provided to the agent pod), or `kubeConfig` to use credentials from `~/.kube/config`. |
| node | string | <no value> | The node name to limit the discovery of pod, port, and node endpoints. Providing no value (the default) results in discovering endpoints for all available nodes. |
| observe_pods | bool | `true` | Whether to report observer pod and port endpoints. If `true` and `node` is specified it will only discover pod and port endpoints whose `spec.nodeName` matches the provided node name. If `true` and `node` isn't specified, it will discover all available pod and port endpoints. Please note that Collector connectivity to pods from other nodes is dependent on your cluster configuration and isn't guaranteed. | 
| observe_nodes | bool | `false` | Whether to report observer k8s.node endpoints. If `true` and `node` is specified it will only discover node endpoints whose `metadata.name` matches the provided node name. If `true` and `node` isn't specified, it will discover all available node endpoints. Please note that Collector connectivity to nodes is dependent on your cluster configuration and isn't guaranteed.| 
