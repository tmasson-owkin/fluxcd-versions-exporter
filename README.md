# Fluxcd Versions Exporter

To run it:

    go build
    ./flux_version_exporter [flags]

## Exported Metrics

| Metric | Description | Labels |
| ------ | ------- | ------ |
| fluxcd_up | Was the last fluxcd query successful | |
| fluxcd_version_applied | List the current installed apps version | "namespace", "name", "ready", "message", "revision", "suspended" |

## Flags
    ./flux_version_exporter --help

| Flag | Description | Default |
| ---- | ----------- | ------- |
| log.level | Logging level | `info` |
| web.listen-address | Address to listen on for telemetry | `:8080` |
| web.telemetry-path | Path under which to expose metrics | `/metrics` |
