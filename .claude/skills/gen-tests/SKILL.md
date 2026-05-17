---
name: gen-tests
description: Generate table-driven Go tests for a function or method in this codebase. Use when the user asks to write or add tests, or says "test this function".
allowed-tools: Read Bash(grep *) Bash(go test *)
---

# Generate table-driven tests

## Steps

1. Read the function signature and its package to understand input/output types.
2. Identify the exported error sentinels or error strings it can return.
3. Build a `tests` slice covering: happy path, boundary inputs, and each error branch.
4. Write the test function in the same `_test.go` file as other tests in that package, or create `<file>_test.go` if none exists.

## Exact pattern to follow

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name           string
        // --- inputs ---
        input          InputType
        // --- expected ---
        expectedResult OutputType
        expectedError  error
    }{
        {name: "valid_input",   input: ..., expectedResult: ..., expectedError: nil},
        {name: "invalid_input", input: ..., expectedResult: ..., expectedError: ErrSomething},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := FunctionName(tt.input)
            if !errors.Is(err, tt.expectedError) {
                t.Errorf("error: got %v, want %v", err, tt.expectedError)
            }
            if !reflect.DeepEqual(got, tt.expectedResult) {
                t.Errorf("result: got %v, want %v", got, tt.expectedResult)
            }
        })
    }
}
```

## Constraints

- No testify, gomock, or any third-party assertion library.
- Use `errors.Is` for error comparisons.
- Use `reflect.DeepEqual` for struct/slice comparisons; `==` for scalars and strings.
- Test case names: `snake_case` strings.
- If the function reads files, place fixture files under `testdata/` next to the test file.
- Run `go test ./internal/<pkg>/...` after writing to confirm they pass.
