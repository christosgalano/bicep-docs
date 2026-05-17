---
globs: "**/*_test.go"
---

# Testing conventions

## Table-driven tests

Every test function uses a named-field struct slice with `t.Run`:

```go
func TestParseSection(t *testing.T) {
    tests := []struct {
        name           string
        input          string
        expectedResult types.Section
        expectedError  error
    }{
        {name: "valid_description", input: "description", expectedResult: types.DescriptionSection, expectedError: nil},
        {name: "invalid_section",   input: "unknown",     expectedResult: "",                       expectedError: types.ErrInvalidSection},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := types.ParseSectionFromString(tt.input)
            if !errors.Is(err, tt.expectedError) {
                t.Errorf("error: got %v, want %v", err, tt.expectedError)
            }
            if got != tt.expectedResult {
                t.Errorf("result: got %v, want %v", got, tt.expectedResult)
            }
        })
    }
}
```

## Rules

- No third-party assertion libraries (testify, gomock, etc.). Manual comparisons only.
- Use `reflect.DeepEqual` for struct or slice comparisons.
- Use `errors.Is` for error comparisons, never string matching.
- Test case names use `snake_case` strings (they appear in `go test -v` output).
- Place file-based fixtures in a `testdata/` subdirectory alongside the `_test.go` file.
- Function naming: `TestFunctionName` or `TestTypeName_MethodName`.
- Benchmark functions follow the same file-adjacent pattern: `BenchmarkX(b *testing.B)`.
