# netpol-fqdns

List the egress FQDNs allowed by the `CiliumNetworkPolicy` resources in a Kubernetes cluster.

Cilium lets policies allow egress by DNS name via `toFQDNs`, but those names are
scattered across every policy in every namespace. This CLI collects them into a single
list you can read, diff, or feed into something else.

## Requirements

- Go 1.27.1 (or [mise](https://mise.jdx.dev), which pins it for you)
- A kubeconfig for a cluster running Cilium, with permission to list
  `ciliumnetworkpolicies.cilium.io` cluster-wide

## Install

```sh
go install github.com/henrywhitakercommify/netpol-fqdns@latest
```

Or run it from a checkout with the [mise](https://mise.jdx.dev) task, which takes the
same flags:

```sh
mise run list --format table
mise run list --format json --kubeconfig ~/.kube/staging.yaml
```

## Usage

```
netpol-fqdns [flags]
```

| Flag           | Default          | Description                          |
| -------------- | ---------------- | ------------------------------------ |
| `--kubeconfig` | `~/.kube/config` | Path to the kubeconfig file. `~` is expanded. |
| `--format`     | `table`          | Output format: `table`, `json` or `csv`. |

The cluster and user are taken from the current context of the kubeconfig. Every
namespace is queried, so there is no namespace flag.

## Output formats

### `table`

Human-readable, one row per FQDN, colourised:

```
Name             Namespace  FQDN
api-egress       prod       api.stripe.com
api-egress       prod       *.githubusercontent.com
telemetry        monitoring otel.honeycomb.io
```

### `csv`

The same one-row-per-FQDN shape, for spreadsheets and `diff`:

```csv
name,namespace,fqdn
api-egress,prod,api.stripe.com
api-egress,prod,*.githubusercontent.com
telemetry,monitoring,otel.honeycomb.io
```

### `json`

One object per policy, with the FQDNs nested — the only format that keeps the
grouping:

```json
[
  {
    "name": "api-egress",
    "namespace": "prod",
    "fqdns": ["api.stripe.com", "*.githubusercontent.com"]
  }
]
```

## Development

```sh
mise run list --format json                      # run against your current context
mise run list --kubeconfig ~/.kube/staging.yaml  # ...or another cluster
go build ./...
go vet ./...
```

The code is in two packages:

- `internal/collector` — builds the Cilium clientset from the kubeconfig and lists
  every `CiliumNetworkPolicy`, paging through the API.
- `internal/format` — renders `[]collector.Rule` as a table, CSV or JSON. To add a
  format, add a `func Xxx([]collector.Rule) ([]byte, error)` here and a `case` to the
  switch in `main.go`.

## Status

Work in progress. The collector lists the policies but does not yet extract the
`toFQDNs` rules from them, so `Collect` returns no rules and every format currently
prints an empty result. The output above is the intended shape.
