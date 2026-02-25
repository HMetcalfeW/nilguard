# nilguard

`nilguard` is a Go static analyzer that enforces a defensive coding policy:
**any pointer used in a function must be nil-checked somewhere in that function**.

## Why This Exists

`go vet` and Staticcheck catch *definite* or *inconsistent* nil dereferences. They do **not**
enforce the absence of a nil guard as a policy. `nilguard` fills that gap.

## Installation

### Go Install

```bash
go install github.com/HMetcalfe/nilguard/cmd/nilguard@latest
```

### Homebrew (macOS/Linux)

```bash
brew tap HMetcalfeW/tap
brew install nilguard
```

### From Source

```bash
git clone https://github.com/HMetcalfe/nilguard
cd nilguard
make build
```

## Usage

### Standalone CLI

```bash
nilguard ./...
```

### Via go vet

```bash
go install github.com/HMetcalfe/nilguard/cmd/nilguard-vet@latest
go vet -vettool=$(which nilguard-vet) ./...
```

### golangci-lint Plugin

Build the plugin and reference it in `.golangci.yml`:

```bash
make plugin   # builds bin/nilguard.so (Linux only)
```

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

## What It Detects

nilguard flags pointer **uses** that lack a nil-check anywhere in the same function:

| Use Pattern | Example |
|---|---|
| Star dereference | `*p` |
| Selector on pointer | `p.Field` |
| Method call on pointer | `p.Method()` |

### Qualifying Nil-Checks

- `if p != nil { ... }`
- `if p == nil { return }` (or `panic`, `break`, `continue`, `goto`)

A single qualifying check anywhere in the function satisfies all uses of that pointer.

### Suppression

Add `//nolint:nilguard` to suppress a specific line:

```go
_ = p.X //nolint:nilguard
```

## Known Limitations

- **No alias tracking** — `q := p; q.Method()` is not traced back to `p`
- **No cross-function analysis** — constructors and factory functions are not treated specially
- **No flow-sensitive dominance** — a nil-check anywhere in the function satisfies all uses
- **Nested function literals** — analyzed independently; a check in the outer function does not satisfy uses in a closure
- **matchLabels only** — only simple `p != nil` / `p == nil` conditions; compound expressions like `p != nil && p.X > 5` not yet recognized
- **golangci-lint plugin** — requires `-buildmode=plugin`, which only works on Linux

## Development

Requires Go 1.25+.

```bash
make test       # run tests
make lint       # golangci-lint + staticcheck
make build      # build CLI + vettool
make plugin     # build golangci-lint plugin (Linux only)
make clean      # remove bin/
```

## Releasing

Pushing a version tag triggers the GitHub Actions release pipeline:

```bash
git tag v0.1.0
git push origin v0.1.0
```

## License

Apache 2.0. See [LICENSE](LICENSE) for details.
