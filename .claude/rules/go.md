---
globs: "**/*.go"
---

# Go idioms for this codebase

## Error handling

Wrap errors with context at every call boundary:
```go
fmt.Errorf("failed to parse template: %w", err)
```
Use `errors.New` only for leaf errors that have no cause to wrap. Never swallow or ignore errors.

## Packages

- All implementation lives under `internal/`. There is no `pkg/` public API.
- Package names are lowercase, single words, matching their directory name.
- Exported identifiers use PascalCase; unexported use camelCase.
- Do not create new top-level packages — work within the existing four (`cli`, `markdown`, `template`, `types`).

## Receivers

- Use a **pointer receiver** when the method reads complex struct state (avoids copy overhead and is consistent with the rest of the type's methods).
- Use a **value receiver** only for lightweight conversions like `String()` on small types.
- All methods on a given type must use the same receiver kind.

## Concurrency

Use `errgroup.Group` from `golang.org/x/sync/errgroup` for parallel work. Return the first error; do not spawn raw goroutines.

```go
var g errgroup.Group
g.Go(func() error { return processFile(path) })
return g.Wait()
```

## Regexes

Compile once at package level; never inside a function body:
```go
var moduleRegex = regexp.MustCompile(`^module\s+(\S+)\s+'(\S+)'`)
```

## String building

Use `strings.Builder` for any multi-part string assembly:
```go
var sb strings.Builder
sb.WriteString("## Parameters\n")
```

## JSON

Import json-iterator as the standard JSON package:
```go
json "github.com/json-iterator/go"
```
Use `json.Decoder` for streaming; use `json.Unmarshal` for small in-memory blobs.

## Context

Do not add `context.Context` parameters. This tool has no cancellation or timeout requirements.

## Linting guard rails

Before adding any `//nolint` directive, exhaust reasonable refactors first. When a nolint is unavoidable, always include the specific rule name and a short reason:
```go
//nolint:funlen // this function maps all section types; splitting adds no clarity
```
