# teleport-ui

Go CLI wrapping `tsh` with a fuzzy finder for selecting Teleport databases / kube clusters.

## Workflow rules

Every change must, before being considered done:

1. **Format** — `golangci-lint fmt ./...` (runs gofmt + goimports per `.golangci.yml`).
2. **Lint** — `golangci-lint run ./...`. Fix all findings; do not silence with `//nolint` unless justified.
3. **Test** — add or update tests for changed behavior, then `go test ./...`. All packages green.

Run all three before reporting work complete.
