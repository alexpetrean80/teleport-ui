# teleport-ui

A terminal UI for [Teleport](https://goteleport.com/) that wraps `tsh` with a fuzzy finder for quickly connecting to databases and Kubernetes clusters.

## Prerequisites

- [tsh](https://goteleport.com/docs/connect-your-client/tsh/) installed and authenticated
- Go 1.25+

## Installation

```bash
go install github.com/alexptr80/teleport-ui@latest
```

Or build from source:

```bash
git clone https://github.com/alexptr80/teleport-ui.git
cd teleport-ui
go build -o teleport-ui .
```

## Usage

```
teleport-ui <command> [flags] [filter args...]
```

### Commands

| Command | Description |
|---------|-------------|
| `db`   | List and connect to databases |
| `kube` | List and login to Kubernetes clusters |

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--clear-cache` | `-c` | Clear cached `tsh` data |
| `--help` | `-h` | Show help message |

### Examples

Connect to a database:

```bash
teleport-ui db
```

Login to a Kubernetes cluster:

```bash
teleport-ui kube
```

Search for a specific database:

```bash
teleport-ui db --search=prod
```

Clear cache and refresh (e.g. after a new resource is added):

```bash
teleport-ui db -c
```

Clear all cached data:

```bash
teleport-ui -c
```

## How it works

1. Runs `tsh db ls` or `tsh kube ls` to fetch available resources.
2. Presents a fuzzy finder in the terminal to search and select a resource.
3. For databases, prompts for a database user via a second fuzzy finder, then runs `tsh db connect`.
4. For Kubernetes, runs `tsh kube login` for the selected cluster.

### Caching

Results from `tsh` are cached locally at `~/.cache/teleport-ui/` (respects `XDG_CACHE_HOME`) using Go's gob encoding for fast serialization. Cached data is used on subsequent runs when no filter args are provided, making startup near-instant.

Cache is bypassed when filter args like `--search` are passed, since filtering is handled server-side by `tsh`.

Use `--clear-cache` / `-c` to force a fresh fetch when resources have changed.

### Fuzzy finder controls

| Key | Action |
|-----|--------|
| Type | Filter results |
| `Up` / `Ctrl+K` | Move up |
| `Down` / `Ctrl+J` | Move down |
| `Enter` | Select |
| `Esc` / `Ctrl+C` | Cancel |
