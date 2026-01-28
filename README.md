# nilguard

`nilguard` is a Go analyzer that enforces a defensive coding policy:
**any pointer used in a function must be nil-checked somewhere in that function** (v1).
Nested function literals are excluded.

## Why this exists
`go vet` and Staticcheck catch *definite* or *inconsistent* nil dereferences. They do **not**
enforce the absence of a nil guard as a policy. `nilguard` fills that gap.

## Quick start

```bash
git clone https://github.com/HMetcalfe/nilguard
cd nilguard
make build           # builds CLI and vettool
./bin/nilguard ./... # run as standalone
go vet -vettool=$(pwd)/bin/nilguard-vet ./... # run via go vet
make plugin          # builds golangci-lint plugin (bin/nilguard.so)
```

## golangci-lint plugin

Build the plugin and reference it in `.golangci.yml`:

```yaml
linters:
  enable:
    - nilguard

linters-settings:
  custom:
    nilguard:
      path: ./bin/nilguard.so
      description: "Flags pointer uses without a nil check"
      original-url: "https://github.com/HMetcalfe/nilguard"
```

## Development

Requires Go 1.25+. For local linting, install tools once:

```bash
make tools
```

Common tasks:

```bash
make test
make lint
make build
```

## Cache flags

nilguard defaults to using temporary cache directories to avoid permission
issues in locked-down environments. You can override these via flags:

```bash
./bin/nilguard -cache-dir /tmp/nilguard/go-build -mod-cache-dir /tmp/nilguard/go-mod ./...
```
