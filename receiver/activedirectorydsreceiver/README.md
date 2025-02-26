# Active Directory Domain Services Receiver

<!-- status autogenerated section -->
| Status        |           |
| ------------- |-----------|
| Stability     | [beta]: metrics   |
| Distributions | [contrib], [observiq] |

[beta]: https://github.com/open-telemetry/opentelemetry-collector#beta
[contrib]: https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib
[observiq]: https://github.com/observIQ/observiq-otel-collector
<!-- end autogenerated section -->

The `active_directory_ds` receiver scrapes metric relating to an Active Directory domain controller using the Windows Performance Counters.

## Configuration
The following settings are optional:
- `metrics` (default: see `DefaultMetricsSettings` [here](./internal/metadata/generated_metrics.go)): Allows enabling and disabling specific metrics from being collected in this receiver.
- `collection_interval` (default = `10s`): The interval at which metrics are emitted by this receiver.

Example:
```yaml
receivers:
  active_directory_ds:
    collection_interval: 10s
    metrics:
      # Disable the active_directory.ds.replication.network.io metric from being emitted
      active_directory.ds.replication.network.io: false
```

The full list of settings exposed for this receiver are documented [here](./config.go) with detailed sample configurations [here](./testdata/config.yaml).

## Metrics

Details about the metrics produced by this receiver can be found in [metadata.yaml](./metadata.yaml)
