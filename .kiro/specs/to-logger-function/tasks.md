# Implementation Plan: ToLogger Function

## Overview

Implement the `ToLogger` function in the `structutils` package. The function converts a struct into a flat `[]any` slice of alternating key-value pairs with dot-notation keys for nested fields, sorted alphabetically. All new code goes in `structs.go`; all new tests go in new test files.

## Tasks

- [x] 1. Implement ToLogger function and internal helper
  - [x] 1.1 Add the `pair` struct type and `structToLogger` recursive helper to `structs.go`
    - Define `type pair struct { key string; val any }`
    - Implement `func structToLogger(v reflect.Value, prefix string) []pair`
    - Handle: exported field iteration, pointer dereferencing, nil pointer emission, struct recursion with dot prefix, primitive field emission
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 3.1, 3.2_
  - [x] 1.2 Add the public `ToLogger` function to `structs.go`
    - Validate input (struct or pointer-to-struct, same pattern as `ToMap`)
    - Call `structToLogger`, sort pairs with `slices.SortFunc` and `strings.Compare`, flatten to `[]any`
    - Add required imports: `slices`, `strings`
    - _Requirements: 1.1, 1.2, 1.3, 4.1, 4.2, 5.1, 5.2, 5.3_

- [x] 2. Add example, benchmark, and fuzz tests
  - [x] 2.1 Create `structs_tologger_test.go` with `ExampleToLogger`
    - Demonstrate flat struct and nested struct usage
    - Include `// Output:` comment for `go test` verification
    - _Requirements: 6.1_
  - [x] 2.2 Add `BenchmarkToLogger` and `BenchmarkToLoggerParallel` to the same file
    - Cover: flat struct, nested struct, deep nesting, multiple nested fields
    - Use `b.ReportAllocs()` and follow existing benchmark patterns
    - _Requirements: 6.2_
  - [x] 2.3 Add `FuzzToLogger` to the same file
    - Seed corpus with JSON-encoded structs of varying shapes
    - Decode into a test struct and verify no panic/unexpected error
    - _Requirements: 6.3_

- [x] 3. Checkpoint
  - Ensure all tests pass with `go test ./structutils/...`, ask the user if questions arise.

- [x] 4. Add property-based tests
  - [x] 4.1 Create `structs_tologger_property_test.go` with property test for pointer equivalence
    - **Property 1: Pointer equivalence**
    - **Validates: Requirements 1.2**
    - _Requirements: 6.4, 6.5_
  - [x]* 4.2 Write property test for non-struct rejection
    - **Property 2: Non-struct rejection**
    - **Validates: Requirements 1.3**
    - _Requirements: 6.4_
  - [x]* 4.3 Write property test for even length invariant
    - **Property 3: Even length invariant**
    - **Validates: Requirements 5.1**
    - _Requirements: 6.4_
  - [x]* 4.4 Write property test for keys-are-strings-at-even-indices
    - **Property 4: Keys are strings at even indices**
    - **Validates: Requirements 5.2**
    - _Requirements: 6.4_
  - [x]* 4.5 Write property test for output keys sorted
    - **Property 5: Output keys are sorted**
    - **Validates: Requirements 4.1, 4.2**
    - _Requirements: 6.4_
  - [x]* 4.6 Write property test for deterministic output
    - **Property 6: Deterministic output**
    - **Validates: Requirements 4.3**
    - _Requirements: 6.4_
  - [x]* 4.7 Write property test for nested dot notation
    - **Property 7: Nested fields use dot notation**
    - **Validates: Requirements 2.1, 2.2**
    - _Requirements: 6.4_
  - [x]* 4.8 Write property test for nil pointer key emission
    - **Property 8: Nil pointer produces key with nil value**
    - **Validates: Requirements 2.4**
    - _Requirements: 6.4_
  - [x]* 4.9 Write property test for exported-only fields
    - **Property 9: Only exported fields appear in output**
    - **Validates: Requirements 3.1, 3.2**
    - _Requirements: 6.4_

- [x] 5. Final checkpoint
  - Run `go test ./structutils/...` and `golangci-lint run --fix ./...` to verify zero test failures and zero lint issues. Ask the user if questions arise.

## Notes

* Tasks marked with `*` are optional and can be skipped for faster MVP
* Do NOT modify existing test files (`structs_test.go`, `structs_preservation_test.go`, `structs_bugcondition_test.go`)
* Use `pgregory.net/rapid` for property tests (already a dependency)
* Use `github.com/go-openapi/testify/v2` for assertions
* Use `slices.SortFunc` with `strings.Compare` for sorting
* Follow the same copyright header format as existing files
* All code must pass `golangci-lint run --fix ./...`
