# Implementation Plan

* [x] 1. Write bug condition exploration test
  * **Property 1: Bug Condition** - Nested Struct Raw Value Bug
  * **CRITICAL**: This test MUST FAIL on unfixed code - failure confirms the bug exists
  * **DO NOT attempt to fix the test or the code when it fails**
  * **NOTE**: This test encodes the expected behavior - it will validate the fix when it passes after implementation
  * **GOAL**: Surface counterexamples that demonstrate the bug exists
  * **Scoped PBT Approach**: Scope the property to concrete failing cases: structs containing at least one embedded/named struct field (isBugCondition returns true)
  * **GENERALITY**: `UserInfo`/`Profile` are merely illustrative examples. The test MUST prove the function works for ANY arbitrary struct type. Define diverse struct types beyond the test file types to verify generality.
  * Test that `ToFlattenedMap(UserInfo{Name: "q1mi", Age: 18, Profile: Profile{Hobby: "Two Color Ball"}})` returns `map[string]any{"Name": "q1mi", "Age": 18, "Profile": map[string]any{"Hobby": "Two Color Ball"}}` (from Bug Condition in design)
  * Test that multi-level nesting `A{B: B{C: C{Val: 1}}}` returns `map[string]any{"B": map[string]any{"C": map[string]any{"Val": 1}}}` (nested maps, not flattened)
  * Test that pointer-to-struct field `*Profile` is dereferenced and recursively converted to `map[string]any`
  * Test multiple nested struct fields at the same level: e.g., `Order{Customer: Customer{Name: "Ada"}, Address: Address{City: "London", Country: "UK"}, Total: 99}` returns `map[string]any{"Customer": map[string]any{"Name": "Ada"}, "Address": map[string]any{"City": "London", "Country": "UK"}, "Total": 99}`
  * Test deeply nested structs (3+ levels): e.g., `Company{CEO: Person{Home: Location{Coords: GPS{Lat: 51.5, Lon: -0.1}}}}` returns nested maps at each level
  * Test nil pointer-to-struct field: e.g., `Node{Value: 42, Left: nil, Right: &Node{Value: 7, Left: nil, Right: nil}}` (where `Left`/`Right` are `*Node`) returns `map[string]any{"Value": 42, "Left": nil, "Right": map[string]any{"Value": 7, "Left": nil, "Right": nil}}`
  * Test struct with unexported fields: verify no panic occurs, only exported fields appear in output
  * Assert nested struct fields produce `map[string]any` values, not raw struct values
  * Assert no inner struct field name (e.g., `Hobby`) leaks to the top-level map
  * Run test on UNFIXED code
  * **EXPECTED OUTCOME**: Test FAILS (this is correct - it proves the bug exists: raw struct values stored, inner fields flattened to top level)
  * Document counterexamples found to understand root cause (missing `continue` after struct detection, queue-based flattening)
  * Mark task complete when test is written, run, and failure is documented
  * _Requirements: 1.1, 1.2, 1.3, 2.1, 2.2, 2.3, 2.4_

* [x] 2. Write preservation property tests (BEFORE implementing fix)
  * **Property 2: Preservation** - Primitive-Only Struct Behavior
  * **IMPORTANT**: Follow observation-first methodology
  * **GENERALITY**: `Profile` is merely one example. The preservation tests MUST verify behavior across diverse primitive-only struct types to prove the function is truly generic. Define multiple struct types with different primitive field combinations.
  * Observe: `ToFlattenedMap(Profile{Hobby: "x"})` returns `map[string]any{"Hobby": "x"}` on unfixed code
  * Observe: `ToFlattenedMap(&Profile{Hobby: "y"})` returns `map[string]any{"Hobby": "y"}` on unfixed code (pointer dereference)
  * Observe: `ToFlattenedMap(42)` returns error `"ToFlattenedMap only accepts struct or struct pointer; received int"` on unfixed code
  * Write property-based test: for all structs containing only primitive fields (isBugCondition returns false), the output map contains exactly the field names as keys with their corresponding values
  * Use diverse struct types beyond `Profile`: e.g., `Settings{Verbose bool; Timeout int; Name string}`, `Metrics{Count int; Rate float64; Label string}`, `Mixed{A string; B int; C bool; D float64}` — verify each produces the expected flat map
  * Write property-based test: for pointer-to-primitive-only-struct inputs, result matches the non-pointer version (test with diverse types, not just `*Profile`)
  * Write property-based test: for non-struct inputs (int, string, slice, etc.), function returns an error matching the expected format
  * Verify tests pass on UNFIXED code
  * **EXPECTED OUTCOME**: Tests PASS (this confirms baseline behavior to preserve)
  * Mark task complete when tests are written, run, and passing on unfixed code
  * _Requirements: 3.1, 3.2, 3.3_

* [x] 3. Fix for nested struct handling in ToFlattenedMap

  * [x] 3.1 Implement the fix
    * Replace queue-based iteration with recursive helper function `structToMap(v reflect.Value) map[string]any`
    * `ToFlattenedMap` validates input (pointer dereference, struct check) then delegates to `structToMap`
    * In `structToMap`: iterate fields, if field is pointer then dereference (nil pointer stores nil), if dereferenced field is struct then recurse, otherwise store `field.Interface()`
    * Must use `fieldType.IsExported()` guard to skip unexported fields — `reflect.Value.Interface()` panics on them
    * Must handle nil pointer-to-struct gracefully: store `nil` in the output map, do not panic
    * Must support multi-level pointer indirection via a dereference loop (`for field.Kind() == reflect.Pointer`)
    * **No type-specific logic allowed** — operate solely on `reflect.Kind`. The function must work for any arbitrary struct type at any nesting depth.
    * Must handle multiple nested struct fields at the same level independently (no shared state between field iterations)
    * Remove `github.com/davecgh/go-spew/spew` import from `structutils/structs.go`
    * Remove `spew.Dump(ti)` call
    * Run `go mod tidy` in `structutils/` to clean `go.mod` and `go.sum`
    * _Bug_Condition: isBugCondition(input) where input is a struct containing at least one field of kind reflect.Struct_
    * _Expected_Behavior: nested struct fields are recursively converted to map[string]any, no inner field leaks to parent map_
    * _Preservation: primitive-only structs, pointer dereference, and error path remain unchanged_
    * _Requirements: 2.1, 2.2, 2.3, 2.4, 3.1, 3.2, 3.3_

  * [x] 3.2 Verify bug condition exploration test now passes
    * **Property 1: Expected Behavior** - Nested Struct Raw Value Bug
    * **IMPORTANT**: Re-run the SAME test from task 1 - do NOT write a new test
    * The test from task 1 encodes the expected behavior
    * When this test passes, it confirms the expected behavior is satisfied
    * Run bug condition exploration test from step 1
    * **EXPECTED OUTCOME**: Test PASSES (confirms bug is fixed)
    * _Requirements: 2.1, 2.2, 2.3, 2.4_

  * [x] 3.3 Verify preservation tests still pass
    * **Property 2: Preservation** - Primitive-Only Struct Behavior
    * **IMPORTANT**: Re-run the SAME tests from task 2 - do NOT write new tests
    * Run preservation property tests from step 2
    * **EXPECTED OUTCOME**: Tests PASS (confirms no regressions)
    * Confirm all tests still pass after fix (no regressions)

* [x] 4. Checkpoint - Ensure all tests pass
  * Run `go test ./...` from `structutils/` directory
  * Verify `go build ./...` succeeds without `spew` dependency
  * Verify `go mod tidy` produces no changes (clean state)
  * Ensure all tests pass, ask the user if questions arise.
