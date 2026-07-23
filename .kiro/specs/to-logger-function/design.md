# Design Document: ToLogger Function

## Overview

`ToLogger` converts a Go struct into a flat `[]any` slice of alternating key-value pairs, suitable for direct use with `slog` structured logging. Unlike `ToMap` which produces nested `map[string]any`, `ToLogger` recursively flattens nested structs using dot notation for keys (e.g., `"Profile.Hobby"`), sorts all keys alphabetically, and returns the result as a flat slice ready for variadic slog calls.

This function complements the existing `ToMap` function in the `structutils` package, reusing the same input validation pattern but producing logger-optimized output.

## Architecture

The implementation follows the same architectural pattern as `ToMap`:

```text
ToLogger(in any)
  ├── Input validation (reflect.ValueOf → Kind check)
  ├── structToLogger(v reflect.Value, prefix string)
  │     ├── Iterate exported fields
  │     ├── Dereference pointers
  │     ├── Recurse on struct fields (with prefix)
  │     └── Emit primitive fields as {key, value} pairs
  ├── Sort collected pairs by key
  └── Flatten to []any
```

The function lives in `structs.go` alongside `ToMap`, sharing the same package and file.

## Components and interfaces

### Public API

```go
// ToLogger converts a struct to a flat []any slice of alternating
// key-value pairs suitable for slog structured logging. Nested
// struct fields are flattened using dot notation.
func ToLogger(in any) ([]any, error)
```

### Internal helper

```go
// pair holds a single key-value entry during collection before
// sorting.
type pair struct {
    key string
    val any
}

// structToLogger recursively collects exported fields from a struct
// value into pairs, using prefix for dot-notation nesting.
func structToLogger(v reflect.Value, prefix string) []pair
```

### Dependencies

* `reflect` — struct field iteration and type introspection
* `fmt` — error formatting
* `slices` — `SortFunc` for deterministic ordering
* `strings` — `Compare` for sort comparator

## Data models

### Input

Any Go value. Valid inputs are:

* A struct value (any struct type)
* A pointer to a struct (dereferenced before processing)

All other types produce an error.

### Output

A `[]any` slice with the following structure:

```text
[key₁, val₁, key₂, val₂, ..., keyₙ, valₙ]
```

Where:

* Each `keyᵢ` is a `string` (dot-notation field path)
* Each `valᵢ` is the field's value as `any`
* Keys are sorted alphabetically
* Length is always `2 * (number of exported leaf fields)`

### Internal pair type

```go
type pair struct {
    key string  // dot-notation field path
    val any     // field value (primitive or nil for nil pointers)
}
```

Used only during collection; not exposed in the public API.

## Algorithm

```go
func ToLogger(in any) ([]any, error) {
    v := reflect.ValueOf(in)
    if v.Kind() == reflect.Pointer {
        v = v.Elem()
    }
    if v.Kind() != reflect.Struct {
        return nil, fmt.Errorf(
            "ToLogger only accepts struct or struct pointer; received %T", v)
    }

    pairs := structToLogger(v, "")

    slices.SortFunc(pairs, func(a, b pair) int {
        return strings.Compare(a.key, b.key)
    })

    result := make([]any, 0, len(pairs)*2)
    for _, p := range pairs {
        result = append(result, p.key, p.val)
    }
    return result, nil
}

func structToLogger(v reflect.Value, prefix string) []pair {
    var pairs []pair
    t := v.Type()

    for i := range v.NumField() {
        field := v.Field(i)
        fieldType := t.Field(i)

        if !fieldType.IsExported() {
            continue
        }

        key := prefix + fieldType.Name

        // Dereference pointers
        for field.Kind() == reflect.Pointer {
            if field.IsNil() {
                pairs = append(pairs, pair{key: key, val: nil})
                break
            }
            field = field.Elem()
        }

        // If pointer was nil, already stored
        if field.Kind() == reflect.Invalid {
            continue
        }
        if _, stored := findKey(pairs, key); stored {
            continue
        }

        if field.Kind() == reflect.Struct {
            pairs = append(pairs, structToLogger(field, key+".")...)
            continue
        }

        pairs = append(pairs, pair{key: key, val: field.Interface()})
    }

    return pairs
}
```

Note: The "findKey" check mirrors the nil-pointer-already-stored pattern from `ToMap`. In practice, the simplest approach is to track whether the nil break was hit using a local boolean, avoiding a linear scan.

## Correctness Properties

_A property is a characteristic or behavior that should hold true across all valid executions of a system — essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees._

### Property 1: Pointer equivalence

_For any_ struct value `s`, `ToLogger(s)` and `ToLogger(&s)` SHALL produce identical output slices.

**Validates: Requirements 1.2**

### Property 2: Non-struct rejection

_For any_ non-struct, non-pointer-to-struct value, `ToLogger` SHALL return a nil slice and a non-nil error.

**Validates: Requirements 1.3**

### Property 3: Even length invariant

_For any_ valid struct input, the output slice length SHALL be even.

**Validates: Requirements 5.1**

### Property 4: Keys are strings at even indices

_For any_ valid struct input, every element at an even index in the output slice SHALL be of type `string`.

**Validates: Requirements 5.2**

### Property 5: Output keys are sorted

_For any_ valid struct input, the string keys extracted from even indices SHALL be in non-decreasing alphabetical order.

**Validates: Requirements 4.1, 4.2**

### Property 6: Deterministic output

_For any_ struct value, calling `ToLogger` twice on the same input SHALL produce identical output slices.

**Validates: Requirements 4.3**

### Property 7: Nested fields use dot notation

_For any_ struct containing nested struct fields, all keys corresponding to nested fields SHALL contain at least one dot character, and the prefix before the dot SHALL match the parent field name.

**Validates: Requirements 2.1, 2.2**

### Property 8: Nil pointer produces key with nil value

_For any_ struct containing a nil pointer-to-struct field, the output SHALL contain the field's dot-notation key paired with a nil value.

**Validates: Requirements 2.4**

### Property 9: Only exported fields appear in output

_For any_ struct with a mix of exported and unexported fields, no key in the output SHALL correspond to an unexported field name.

**Validates: Requirements 3.1, 3.2**

## Error handling

| Condition | Behavior |
|-----------|----------|
| Non-struct, non-pointer input | Return `nil, error` with message: `"ToLogger only accepts struct or struct pointer; received %T"` |
| Nil pointer-to-struct field | Emit key with `nil` value (no error) |
| Struct with no exported fields | Return empty `[]any{}`, no error |
| Nil pointer argument (e.g., `(*T)(nil)`) | The `reflect.ValueOf` + `Elem()` produces zero Value; treat as invalid kind → return error |

The error message pattern mirrors `ToMap` but references `ToLogger` for clarity.

## Testing strategy

### Unit tests (ExampleToLogger)

* Demonstrate typical usage with a flat struct
* Demonstrate nested struct with dot-notation output
* Verify output is suitable for `slog.Info("msg", result...)`

### Benchmark tests

* `BenchmarkToLogger` — flat struct, nested struct, deep nesting, multiple nested
* `BenchmarkToLoggerParallel` — same cases with `b.RunParallel`

### Fuzz test

* `FuzzToLogger` — JSON-seeded fuzzing targeting struct decoding paths

### Property-based tests (pgregory.net/rapid)

Each correctness property above maps to a dedicated property test:

* Minimum 100 iterations per property test (rapid default)
* Each test tagged with: **Feature: to-logger-function, Property N: [title]**
* Tests use diverse struct types (flat, nested, deeply nested, pointer fields, mixed exported/unexported)
* Test file: `structs_tologger_property_test.go`
