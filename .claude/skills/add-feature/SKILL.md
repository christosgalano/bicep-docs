---
name: add-feature
description: Add support for a new Bicep language feature to bicep-docs. Follows the fixture → build → inspect ARM JSON → update internal code pipeline. Use when the user wants to implement a new capability identified by the bicep-researcher agent.
allowed-tools: Read Write Edit Bash
---

# Add a new Bicep language feature to bicep-docs

The canonical implementation pipeline:

1. Write a minimal Bicep fixture that exercises the new feature
2. Compile it with `az bicep build` to produce real ARM JSON
3. Inspect the ARM JSON to understand the exact schema
4. Update `internal/types/types.go` — add fields to the right struct(s)
5. Update `internal/types/unmarshal.go` — add JSON unmarshaling for new fields
6. Update `internal/markdown/markdown.go` — add rendering (new column, row, or section)
7. Update `internal/markdown/create.go` if a new section toggle is needed
8. Add test cases to the relevant `*_test.go` files
9. Verify: `go test ./...` then `task lint`

---

## Step 1 — Write the Bicep fixture

Create `test/<feature-name>.bicep` with the **minimum** Bicep that exercises the new feature. Keep it alongside the existing `test/any.bicep` pattern — one focused fixture per capability.

```bicep
// test/secure-params.bicep  (example for @secure() feature)
@secure()
param adminPassword string

@secure()
param config object
```

Also add a companion fixture in `internal/template/testdata/<feature-name>.bicep` if the feature needs to be covered by the template package unit tests (check whether existing fixtures like `extended.bicep` can be extended instead of creating a new one).

---

## Step 2 — Compile to ARM JSON

```bash
az bicep build --file test/<feature-name>.bicep --outfile test/<feature-name>.json
```

If `az` is not available, try `bicep build test/<feature-name>.bicep --outfile test/<feature-name>.json`.

Read the generated JSON immediately — it is the ground truth for what bicep-docs must parse.

---

## Step 3 — Inspect the ARM JSON

Focus on the sections that bicep-docs parses:

| Bicep construct | ARM JSON location |
|---|---|
| Parameters | `$.parameters.<name>` |
| Outputs | `$.outputs.<name>` |
| Variables | `$.variables` |
| User-defined types (UDDTs) | `$.definitions.<name>` |
| User-defined functions (UDFs) | `$.functions[*].members.<name>` |
| Resources | Not in ARM JSON — parsed from `.bicep` source |
| Modules | Not in ARM JSON — parsed from `.bicep` source |

Identify: which new JSON key(s) appear? What is their type (bool, string, object, array)? Are they present on existing objects or do they introduce a new top-level section?

---

## Step 4 — Update `internal/types/types.go`

Add fields to the relevant struct(s). Follow existing patterns:

- Use `json:"keyName"` struct tags matching the ARM JSON key exactly.
- Use `omitempty` for optional fields.
- Use pointer types (`*bool`, `*string`) when the zero value (`false`, `""`) would be indistinguishable from "field absent" — especially important for bool flags.
- For nested objects, add a new named struct rather than an anonymous one.

```go
// Example: adding Secure to Parameter and Output
type Parameter struct {
    // ... existing fields ...
    Secure bool `json:"-"` // derived from type string; not a direct JSON field
}
```

If the feature introduces a new top-level ARM JSON section (rare), add it to `Template` struct and add a corresponding `Section` constant and `ParseSectionFromString` case.

---

## Step 5 — Update `internal/types/unmarshal.go`

Add or extend the `UnmarshalJSON` method for the affected type. The existing pattern uses a type alias to avoid infinite recursion:

```go
func (p *Parameter) UnmarshalJSON(data []byte) error {
    type Alias Parameter
    aux := &struct{ *Alias }{Alias: (*Alias)(p)}
    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }

    // Handle type/ref first (existing)
    tr, err := unmarshalTypeOrRef(data)
    if err != nil {
        return err
    }
    p.Type = tr

    // NEW: derive Secure from the resolved type string
    p.Secure = p.Type == "secureString" || p.Type == "secureObject"

    return nil
}
```

For fields that appear as nested JSON objects (e.g., `discriminator`), unmarshal a minimal helper struct:

```go
var raw struct {
    Discriminator *struct {
        PropertyName string            `json:"propertyName"`
        Mapping      map[string]any    `json:"mapping"`
    } `json:"discriminator"`
}
if err := json.Unmarshal(data, &raw); err != nil {
    return err
}
if raw.Discriminator != nil {
    u.Discriminator = &Discriminator{
        PropertyName: raw.Discriminator.PropertyName,
        Mapping:      raw.Discriminator.Mapping,
    }
}
```

---

## Step 6 — Update `internal/markdown/markdown.go`

Find the function that generates the affected section (e.g., `generateParametersSection`, `generateUserDefinedDataTypesSection`). Follow the `MarkdownTable` pattern:

```go
// To add a "Secure" column to Parameters:
// 1. Add "Secure" to the headers slice
// 2. Add the value to each row slice in the same position

headers := []string{"Name", "Status", "Type", "Secure", "Description", "Default"}
// ...
row := []string{
    p.Name,
    p.GetStatus().String(),
    extractType(p.Type),
    boolToYesNo(p.Secure),  // new field
    p.Description,
    formatDefault(p.DefaultValue),
}
```

The `--show-all-decorators` flag pattern (see `showAllDecorators` parameter in `generateParametersSection`) controls whether sparse columns are shown. Consider whether the new column should be gated behind that flag or always shown.

For a completely new section, follow the pattern of `generateUserDefinedDataTypesSection`: create a `generate<Feature>Section` function, add it to `create.go`'s `GenerateDocumentation` dispatch, and add a `Section` constant.

---

## Step 7 — Add test cases

### Unit test for the new field (`internal/types/*_test.go`)

Add a table-driven test case to the existing test that covers the struct being changed. Use a real ARM JSON snippet (copy from the compiled fixture) as the input:

```go
{
    name: "secure_string_param",
    input: []byte(`{"type":"secureString","metadata":{"description":"admin password"}}`),
    expectedResult: types.Parameter{
        Type:        "secureString",
        Secure:      true,
        Description: "admin password",
    },
    expectedError: nil,
},
```

### Markdown rendering test (`internal/markdown/*_test.go`)

Add a case to the relevant table-driven test in `create_test.go` or `markdown_test.go` that exercises the new column/section with a minimal `types.Template`.

### Template/parse test (`internal/template/*_test.go`)

If the feature affects parsing of the `.bicep` source (not just ARM JSON), add the fixture to the testdata directory and add a test case in `parse_test.go` or `build_test.go`.

---

## Step 8 — Verify

```bash
go test ./...
task lint
```

Fix any lint issues before considering the feature done. Pay attention to:
- `funlen` violations if a table or switch grew too long — may need `//nolint:funlen // reason`
- `mnd` (magic numbers) for any new numeric constants — define them as named constants
- `godot` — all new exported comments must end with a period
