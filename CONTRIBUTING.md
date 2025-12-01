# Contributing to Resgate

## Dev Environment Setup

**Requirements:** Go 1.21+, Git, make

1. Fork the repository on GitHub: [github.com/mdryaan/resgate](https://github.com/mdryaan/resgate)

2. Clone your fork:

```bash
git clone https://github.com/your_username/resgate.git
cd resgate
go mod tidy
make build
./resgate --help
```

3. Add the upstream remote:

```bash
git remote add upstream https://github.com/mdryaan/resgate.git
```

Verify your build with a quick workflow:

```bash
./resgate pool create --name pool1 --cpu 16 --memory 32768 --gpu 4
./resgate tenant add --name dev --priority 3
./resgate reserve --tenant dev --pool pool1 --cpu 2 --memory 4096
./resgate status
```

---

## Adding a New Command

1. Create `cmd/<name>.go` in the `cmd` package
2. Define a `*cobra.Command` and register it in `init()` via `rootCmd.AddCommand`
3. Access the engine through the package-level `engine` variable — it is initialized in `root.go`'s `PersistentPreRunE`
4. Use `pkg/output` for all terminal output — never `fmt.Println` directly in commands

```go
package cmd

import (
    "github.com/mdryaan/resgate/pkg/output"
    "github.com/spf13/cobra"
)

var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Short description",
    RunE: func(cmd *cobra.Command, args []string) error {
        if someFlag == "" {
            output.Fatal("--someflag is required")
        }
        // use engine.*
        return nil
    },
}

func init() {
    myCmd.Flags().StringVar(&someFlag, "someflag", "", "description")
    rootCmd.AddCommand(myCmd)
}
```

---

## Adding a New Exporter

1. Create `pkg/exporter/<format>_exporter.go`
2. Implement the `Exporter` interface:

```go
type Exporter interface {
    Export(r *Report) ([]byte, error)
}
```

3. Register it in `pkg/exporter/exporter.go` inside `New()`:

```go
case "xml":
    return &XMLExporter{}, nil
```

---

## Extending the Reservation Engine

The engine is split across files in `pkg/reservation/`:

- `engine.go` — constructor and delegation methods
- `lifecycle.go` — `Reserve` and `Unreserve`
- `preemptor.go` — `Preempt`
- `expiry.go` — TTL sweep logic
- `conflict.go` — duplicate reservation detection

Add new behavior as methods on `*Engine`. If the method mutates store state, acquire `e.store.Lock()` at the top. If it only reads, use `e.store.RLock()`. Never nest locks.

---

## PR Guidelines

- One logical change per PR
- Run `go mod tidy` and `make build` before submitting
- PR title follows conventional commits: `feat(scope): description`
- Include command output in the PR description for new commands
- No comments in code — names must be self-explanatory

---

## Code Style Rules

- Zero comments — names must speak for themselves
- No `fmt.Println` in commands — use `pkg/output` functions
- Acquire the narrowest lock scope needed — prefer `RLock` for reads
- Never nest `store.Lock` inside another `store.Lock` — deadlock guaranteed
- All domain types live in `internal/models/` — never define them inside `pkg/`
- Return errors upward; call `output.Fatal` only at the command boundary
- Strong types everywhere — no `interface{}`, no `map[string]interface{}`
