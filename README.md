# nilguard

`nilguard` is a Go analyzer that enforces a defensive coding policy:
**any pointer used in a function must be nil-checked somewhere in that function** (v1).
Nested function literals are excluded.

## Why this exists
`go vet` and Staticcheck catch *definite* or *inconsistent* nil dereferences. They do **not**
enforce the absence of a nil guard as a policy. `nilguard` fills that gap.

## Quick start

```bash
git clone https://github.com/your-org/nilguard
cd nilguard
make build           # builds CLI and vettool
./bin/nilguard ./... # run as standalone
go vet -vettool=$(pwd)/bin/nilguard-vet ./... # run via go vet
