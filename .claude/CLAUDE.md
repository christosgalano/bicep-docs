# bicep-docs

CLI tool that generates Markdown documentation from Azure Bicep template files. It invokes the Azure CLI to compile Bicep to ARM JSON, parses the compiled template, and writes structured Markdown output.

## Directory layout

```
cmd/bicep-docs/     # main() — thin entrypoint, delegates to cli.Execute()
internal/
  cli/              # Cobra command definitions, flag parsing, orchestration
  markdown/         # Markdown generation from parsed types
  template/         # Bicep→ARM build (az cli) and ARM JSON parsing
  types/            # Domain models: Parameter, Module, Resource, Metadata, etc.
test/               # Integration-level fixtures (.bicep + compiled ARM JSON)
examples/           # Example inputs and expected outputs
```

Dependencies flow one way: `cli → markdown/template → types`. No circular deps.

## Build and run

```bash
task build        # compile to ./bin/bicep-docs
task test         # run all tests (gotestsum)
task lint         # golangci-lint run (37 linters, must pass before commit)
task security     # gosec scan
task coverage     # per-package coverage reports
task benchmark    # template parse + doc generation benchmarks
```

Direct Go commands also work (`go build ./cmd/bicep-docs`, `go test ./...`).

## Non-obvious conventions

- **No context.Context** — all operations are synchronous; do not add context parameters.
- **Concurrency** — parallel file processing uses `errgroup.Group` from `golang.org/x/sync`. No raw goroutines.
- **JSON** — import json-iterator/go as `json "github.com/json-iterator/go"`. Use `json.Decoder` for streaming large ARM templates.
- **Regexes** — compile once at package level with `regexp.MustCompile`; never inside functions.
- **String building** — use `strings.Builder` with `WriteString`; never `+` concatenation in loops.
- **Logging** — no structured logging library. `fmt.Printf` for verbose output, `fmt.Fprintln(os.Stderr, ...)` for errors.
- **Section enum** — `types.Section` is a string type; new sections need `ParseSectionFromString` updates and a constant.
- **Linting** — 200-char line limit, gocyclo max 12, funlen max 50 statements. Add `//nolint:rulename // reason` only when unavoidable.
