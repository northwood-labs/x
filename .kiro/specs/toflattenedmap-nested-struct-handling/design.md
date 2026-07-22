# ToFlattenedMap nested struct handling bugfix design

## Overview

The `ToFlattenedMap` function in `structutils/structs.go` incorrectly handles nested/embedded structs. It stores raw struct values in the output map instead of recursively converting them to `map[string]any`, and it flattens inner struct fields into the top-level map. The fix replaces the iterative queue-based approach with a recursive, reflection-based function that properly converts nested structs into `map[string]any` at each level, and removes the unused `spew` dependency.

The function is fully generic via reflection and must correctly handle any arbitrary struct type at any nesting depth. The types used in the test file (`UserInfo`, `Profile`) are merely illustrative examples. The implementation must support all valid Go struct compositions: multiple nested struct fields at the same level, deeply nested hierarchies, pointer-to-struct fields (including nil pointers), mixed primitive and struct fields, and unexported fields. No type-specific logic is permitted — the recursive helper must operate solely on `reflect.Value` and `reflect.Kind`.

## Glossary

* **Bug_Condition (C)**: The condition that triggers the bug — when the input struct contains at least one field whose type is itself a struct (embedded or named).
* **Property (P)**: The desired behavior — nested struct fields are recursively converted to `map[string]any` and stored under their field name, without flattening inner fields into the parent map.
* **Preservation**: Existing behavior for structs with only primitive fields, pointer dereferencing, and error handling for non-struct inputs must remain unchanged.
* **ToFlattenedMap**: The function in `structutils/structs.go` that converts a struct to a `map[string]any`.
* **Recursive conversion**: The approach of calling the conversion function on each nested struct field and storing the returned map as the value.

## Bug details

### Bug condition

The bug manifests when a struct contains an embedded or named struct field. The `ToFlattenedMap` function adds the nested struct to a processing queue (which flattens its fields into the top-level map) AND falls through to the "General Fields" section which stores the raw struct value. This produces two defects: raw struct values in the map, and unwanted flattening of inner fields.

**Formal Specification:**

```text
FUNCTION isBugCondition(input)
  INPUT: input of type any (argument to ToFlattenedMap)
  OUTPUT: boolean

  v := reflect.ValueOf(input)
  IF v.Kind() == reflect.Pointer THEN
    v = v.Elem()
  END IF

  IF v.Kind() != reflect.Struct THEN
    RETURN false
  END IF

  FOR i := 0 TO v.NumField() - 1 DO
    field := v.Field(i)
    IF field.Kind() == reflect.Pointer THEN
      field = field.Elem()
    END IF
    IF field.Kind() == reflect.Struct THEN
      RETURN true
    END IF
  END FOR

  RETURN false
END FUNCTION
```

### Examples

* `UserInfo{Name: "q1mi", Age: 18, Profile: Profile{Hobby: "Two Color Ball"}}` — Expected: `map[string]any{"Name": "q1mi", "Age": 18, "Profile": map[string]any{"Hobby": "Two Color Ball"}}`. Actual: `map[string]any{"Age": 18, "Hobby": "Two Color Ball", "Name": "q1mi", "Profile": structutils.Profile{Hobby: "Two Color Ball"}}`.
* A struct with a pointer to another struct — Expected: the pointed-to struct is recursively converted. Actual: the raw struct value is stored and its fields are flattened.
* A struct with only primitive fields (e.g., `Profile{Hobby: "x"}`) — Expected and Actual are both correct: `map[string]any{"Hobby": "x"}`. This case is unaffected.
* A struct with multiple levels of nesting (e.g., `A{B: B{C: C{Val: 1}}}`) — Expected: `map[string]any{"B": map[string]any{"C": map[string]any{"Val": 1}}}`. Actual: all fields flattened to top level with raw struct values.
* **Multiple nested fields at the same level**: `Order{Customer: Customer{Name: "Ada"}, Address: Address{City: "London", Country: "UK"}, Total: 99}` — Expected: `map[string]any{"Customer": map[string]any{"Name": "Ada"}, "Address": map[string]any{"City": "London", "Country": "UK"}, "Total": 99}`.
* **Deeply nested (3+ levels)**: `Company{CEO: Person{Home: Location{Coords: GPS{Lat: 51.5, Lon: -0.1}}}}` — Expected: `map[string]any{"CEO": map[string]any{"Home": map[string]any{"Coords": map[string]any{"Lat": 51.5, "Lon": -0.1}}}}`.
* **Mixed types with nil pointer-to-struct**: `Node{Value: 42, Left: nil, Right: &Node{Value: 7, Left: nil, Right: nil}}` (where `Left`/`Right` are `*Node`) — Expected: `map[string]any{"Value": 42, "Left": nil, "Right": map[string]any{"Value": 7, "Left": nil, "Right": nil}}`.
* **Struct with unexported fields**: Only exported fields appear in the output map. Unexported fields are skipped because `reflect.Value.Interface()` panics on them.
* **Wide struct with many nested fields**: `Config{DB: DBConfig{Host: "localhost", Port: 5432}, Cache: CacheConfig{TTL: 60}, Auth: AuthConfig{Secret: "s"}}` — Expected: each nested struct becomes its own `map[string]any` under its field name, with no cross-contamination between sibling maps.

## Expected behavior

### Preservation requirements

**Unchanged Behaviors:**

* Structs containing only primitive fields (string, int, bool, etc.) must continue to produce a flat map keyed by field name.
* Pointer-to-struct inputs must continue to be dereferenced and produce the same result as passing the struct by value.
* Non-struct inputs must continue to return an error with the message format `"ToFlattenedMap only accepts struct or struct pointer; received %T"`.
* The function signature `func ToFlattenedMap(in any) (map[string]any, error)` must remain unchanged.

**Scope:**

All inputs that do NOT contain nested struct fields should be completely unaffected by this fix. This includes:

* Structs with only primitive fields
* Pointer-to-struct inputs (dereferencing behavior)
* Non-struct inputs (error path)

## Hypothesized root cause

Based on the code analysis, the root causes are:

1. **Missing `continue` after struct detection**: When a struct field is encountered (line 63), it is added to the queue but execution falls through to the "General Fields" section (line 67-69), which stores the raw `vi.Interface()` value. A `continue` statement or an `else` branch is needed to skip the general assignment for struct fields.

2. **Queue-based iteration flattens all levels**: The queue approach processes all nested structs in a flat loop, writing their fields directly into the single `out` map. This design fundamentally cannot produce nested `map[string]any` values — it can only flatten everything to one level.

3. **Pointer-to-struct also falls through**: The pointer handling block (lines 50-61) adds pointer-to-struct values to the queue but does not skip the subsequent struct check or the general fields assignment, potentially processing the same field multiple times.

4. **Unused debug dependency**: The `spew.Dump(ti)` call on line 58 is a debugging artifact that should be removed, along with the `github.com/davecgh/go-spew` import.

## Correctness properties

Property 1: Bug Condition - Nested structs are recursively converted to maps

_For any_ input struct that contains at least one field of kind `reflect.Struct` (isBugCondition returns true), the fixed `ToFlattenedMap` function SHALL store a `map[string]any` value (produced by recursively converting the nested struct) under the field's name, and SHALL NOT include any of the nested struct's fields as top-level keys in the parent map.

**Validates: Requirements 2.1, 2.2, 2.3.**

Property 2: Preservation - Primitive-only structs produce identical output

_For any_ input struct that contains only primitive fields (isBugCondition returns false), the fixed `ToFlattenedMap` function SHALL produce the same `map[string]any` result as the original function, preserving all field names and values.

**Validates: Requirements 3.1, 3.2, 3.3.**

## Fix implementation

### Changes required

Assuming our root cause analysis is correct:

**File**: `structutils/structs.go`

**Function**: `ToFlattenedMap`

**Specific Changes**:

1. **Replace queue-based iteration with recursion**: Convert the function body to use a recursive helper (or make `ToFlattenedMap` itself recursive). The helper accepts a `reflect.Value` of kind `Struct` and returns `map[string]any`.

2. **Handle struct fields recursively**: When a field's kind is `reflect.Struct`, recursively call the conversion function on that field and store the returned `map[string]any` as the value — do not flatten its fields into the parent map.

3. **Handle pointer-to-struct fields**: When a field is a pointer, dereference it. If the dereferenced value is a struct, recursively convert it. If it is a primitive, store it directly. Support multi-level pointer indirection (`**Struct`) via a dereference loop.

4. **Skip general assignment for struct fields**: Use `continue` or an `else` branch to ensure struct-typed fields are not also stored as raw values in the general assignment.

5. **Remove `spew` dependency**: Delete the `spew.Dump(ti)` call and remove the `github.com/davecgh/go-spew/spew` import. Run `go mod tidy` to clean up `go.mod` and `go.sum`.

6. **Handle nil pointers gracefully**: When a pointer-to-struct field is nil, store `nil` in the output map rather than panicking or omitting the key.

7. **Skip unexported fields**: Check `fieldType.IsExported()` before calling `field.Interface()` to prevent panics on arbitrary structs that contain unexported fields.

8. **Handle multiple struct fields at the same level**: The recursive helper must process each field independently — no shared state between field iterations that could cause cross-contamination between sibling nested maps.

9. **Handle deeply nested structs**: The recursion must naturally terminate for any finite nesting depth without an artificial depth limit. Each recursive call processes one level.

### Proposed implementation sketch

```go
func ToMap(in any) (map[string]any, error) {
    v := reflect.ValueOf(in)
    if v.Kind() == reflect.Pointer {
        v = v.Elem()
    }
    if v.Kind() != reflect.Struct {
        return nil, fmt.Errorf("ToFlattenedMap only accepts struct or struct pointer; received %T", v)
    }

    return structToMap(v), nil
}

func structToMap(v reflect.Value) map[string]any {
    out := make(map[string]any)
    t := v.Type()

    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        fieldType := t.Field(i)

        // Skip unexported fields — Interface() panics on them.
        if !fieldType.IsExported() {
            continue
        }

        // Dereference pointers (possibly multiple levels).
        for field.Kind() == reflect.Pointer {
            if field.IsNil() {
                out[fieldType.Name] = nil
                break
            }
            field = field.Elem()
        }

        // If the pointer was nil, we already stored nil above.
        if field.Kind() == reflect.Invalid {
            continue
        }

        // After nil-pointer handling, check for nil again via the
        // break above (field remains Pointer kind if nil was stored).
        if _, stored := out[fieldType.Name]; stored {
            continue
        }

        if field.Kind() == reflect.Struct {
            out[fieldType.Name] = structToMap(field)
            continue
        }

        out[fieldType.Name] = field.Interface()
    }

    return out
}
```

Key design decisions for arbitrary struct support:

1. **Unexported field guard**: `fieldType.IsExported()` check prevents panics on unexported fields in any user-defined struct.
2. **Multi-level pointer dereference**: The `for` loop handles `**Struct`, `***Struct`, etc.
3. **Nil pointer safety**: Nil pointers store `nil` in the map rather than panicking.
4. **No type-specific logic**: The helper operates solely on `reflect.Kind`, making it work for any struct composition.
5. **Multiple struct fields at the same level**: Each field is processed independently in the loop — no shared state between iterations.

## Testing strategy

### Validation approach

The testing strategy follows a two-phase approach: first, surface counterexamples that demonstrate the bug on unfixed code, then verify the fix works correctly and preserves existing behavior.

### Exploratory bug condition checking

**Goal**: Surface counterexamples that demonstrate the bug BEFORE implementing the fix. Confirm or refute the root cause analysis. If we refute, we will need to re-hypothesize.

**Test Plan**: Write tests that call `ToFlattenedMap` with structs containing embedded struct fields and assert the output matches the expected nested `map[string]any` structure. Run these tests on the UNFIXED code to observe failures.

**Test Cases**:

1. **Single embedded struct**: `UserInfo{Name: "q1mi", Age: 18, Profile: Profile{Hobby: "Two Color Ball"}}` — assert `Profile` key contains `map[string]any` (will fail on unfixed code)
2. **Multi-level nesting**: A struct with two levels of nesting — assert each level is a nested map (will fail on unfixed code)
3. **No top-level leakage**: Assert that inner field names (e.g., `Hobby`) do NOT appear as top-level keys (will fail on unfixed code)
4. **Pointer-to-struct field**: A struct with a `*Profile` field — assert correct recursive conversion (will fail on unfixed code)

**Expected Counterexamples**:

* `Profile` key contains raw `structutils.Profile{}` instead of `map[string]any`
* `Hobby` appears as a top-level key alongside `Profile`
* Possible causes: missing `continue` after struct detection, queue flattening design

### Fix checking

**Goal**: Verify that for all inputs where the bug condition holds, the fixed function produces the expected behavior.

**Pseudocode:**

```text
FOR ALL input WHERE isBugCondition(input) DO
  result := ToFlattenedMap_fixed(input)
  FOR EACH field F in input WHERE F.Kind == reflect.Struct DO
    ASSERT result[F.Name] is of type map[string]any
    ASSERT result[F.Name] == structToMap(F.Value)
    ASSERT no key from F.fields exists at top level of result
  END FOR
END FOR
```

### Preservation checking

**Goal**: Verify that for all inputs where the bug condition does NOT hold, the fixed function produces the same result as the original function.

**Pseudocode:**

```text
FOR ALL input WHERE NOT isBugCondition(input) DO
  ASSERT ToFlattenedMap_original(input) = ToFlattenedMap_fixed(input)
END FOR
```

**Testing Approach**: Property-based testing is recommended for preservation checking because:

* It generates many test cases automatically across the input domain
* It catches edge cases that manual unit tests might miss
* It provides strong guarantees that behavior is unchanged for all non-buggy inputs

**Test Plan**: Observe behavior on UNFIXED code first for primitive-only structs, then write property-based tests capturing that behavior.

**Test Cases**:

1. **Primitive-only struct preservation**: Generate structs with various primitive field types and verify the output map matches field-by-field
2. **Pointer dereferencing preservation**: Verify that `*MyStruct` produces the same result as `MyStruct` for primitive-only structs
3. **Error path preservation**: Verify non-struct inputs (int, string, slice, etc.) continue to return the expected error
4. **Field name preservation**: Verify field names in the output map exactly match struct field names
5. **Diverse struct types**: Use structs beyond `UserInfo`/`Profile` — e.g., `Settings{Verbose bool; Timeout int; Name string}`, `Metrics{Count int; Rate float64}` — to confirm preservation holds for any arbitrary primitive-only struct

### Unit tests

* Test `UserInfo` with embedded `Profile` struct (the primary failing case from the bug report)
* Test multi-level nesting (struct within struct within struct)
* Test pointer-to-struct fields (nil and non-nil)
* Test structs with only primitive fields (regression check)
* Test non-struct inputs return appropriate error
* Test struct with multiple nested struct fields at the same level (e.g., `Order{Customer, Address, Total}`)
* Test deeply nested struct (4+ levels) to verify recursion terminates correctly
* Test struct with unexported fields — verify they are skipped without panic
* Test struct with pointer-to-pointer-to-struct (multi-level indirection)

### Property-based tests

* Generate random structs with nested struct fields and verify all struct-typed values in the output are `map[string]any` (using `pgregory.net/rapid`)
* Generate random primitive-only structs and verify the output matches a naive field-by-field extraction
* Generate random inputs and verify no inner struct field name leaks to the parent level
* **Generality testing with `rapid`**: Define a family of struct types beyond `UserInfo`/`Profile` to verify the function handles arbitrary compositions:
  * Generate structs with varying numbers of nested struct fields (1 to N siblings)
  * Generate deeply nested structs (3-5 levels) and verify recursive map structure
  * Generate structs with nil and non-nil pointer-to-struct fields
  * Generate structs with mixed field types (string, int, bool, float64, nested struct, pointer-to-struct) at each level
* **Arbitrary struct verification strategy**: Since Go's `reflect` cannot create new struct types at runtime in `rapid`, use a set of diverse pre-defined struct types that cover the structural patterns:
  * `FlatStruct` — only primitive fields (various types)
  * `SingleNested` — one nested struct field plus primitives
  * `MultiNested` — multiple nested struct fields at the same level
  * `DeeplyNested` — 3+ levels of nesting
  * `WithPointers` — pointer-to-struct fields (nil and non-nil)
  * `MixedWide` — many fields of different kinds including multiple nested structs
* Use `rapid.Custom` to generate random values for each type and assert structural properties hold across all of them
* Property: for any generated struct, every key in the output map corresponds to an exported field name, and every value that is a `map[string]any` corresponds to a field of kind `reflect.Struct`

### Integration tests

* Test round-trip: convert struct to map and verify all values are accessible via expected key paths
* Test that the function works correctly when called multiple times with different struct types
* Test that the removed `spew` dependency does not cause build failures
