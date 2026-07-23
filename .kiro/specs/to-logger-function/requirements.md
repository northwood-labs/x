# Requirements Document

## Introduction

Add a `ToLogger()` function to the `structutils` package that converts a struct into a flat `[]any` slice of alternating key-value pairs suitable for use with Go's `slog` structured logger. Nested struct fields are flattened using dot notation for keys, enabling ergonomic usage like `slog.Info("event", structutils.ToLogger(user)...)`.

## Glossary

* **System**: The `ToLogger` function within the `structutils` package
* **Flat_Slice**: A `[]any` slice containing alternating key and value entries (`key, value, key, value, ...`)
* **Dot_Notation**: A naming convention where nested field names are joined with periods (e.g., `"Parent.Child.Field"`)
* **Exported_Field**: A Go struct field whose name begins with an uppercase letter, making it accessible outside the declaring package
* **Unexported_Field**: A Go struct field whose name begins with a lowercase letter, inaccessible outside the declaring package

## Requirements

### Requirement 1: struct to logger slice conversion

**User Story:** As a developer, I want to convert a struct into a flat slice of key-value pairs, so that I can pass structured data directly to `slog` logging calls.

#### Acceptance criteria

1. WHEN a struct value is provided, THE System SHALL return a Flat_Slice containing alternating key and value entries for each Exported_Field
2. WHEN a pointer-to-struct value is provided, THE System SHALL dereference the pointer and return the same Flat_Slice as the by-value equivalent
3. IF a non-struct and non-pointer-to-struct value is provided, THEN THE System SHALL return a nil slice and an error describing the invalid input type

### Requirement 2: nested struct flattening

**User Story:** As a developer, I want nested struct fields flattened with dot notation keys, so that the logger output remains a single-level key-value list without nested objects.

#### Acceptance criteria

1. WHEN a struct contains a nested struct field, THE System SHALL flatten the nested fields using Dot_Notation keys (e.g., `"Profile.Hobby"`)
2. WHEN a struct contains multiple levels of nesting, THE System SHALL concatenate all ancestor field names with dots (e.g., `"CEO.Home.Coords.Lat"`)
3. WHEN a struct contains a pointer-to-struct field that is non-nil, THE System SHALL dereference the pointer and flatten the nested fields using Dot_Notation
4. WHEN a struct contains a pointer-to-struct field that is nil, THE System SHALL include the Dot_Notation key with a nil value

### Requirement 3: field filtering

**User Story:** As a developer, I want unexported fields excluded from the output, so that only publicly visible data appears in log entries.

#### Acceptance criteria

1. THE System SHALL skip all Unexported_Fields during conversion
2. WHEN a struct contains a mix of exported and unexported fields, THE System SHALL include only Exported_Fields in the Flat_Slice

### Requirement 4: deterministic output order

**User Story:** As a developer, I want the output keys to be in a deterministic order, so that log output is reproducible and comparable across runs.

#### Acceptance criteria

1. THE System SHALL sort the output keys in alphabetical order
2. WHEN nested struct fields are flattened, THE System SHALL sort using the fully-qualified Dot_Notation key
3. FOR ALL identical inputs, THE System SHALL produce identical output slices

### Requirement 5: output structure invariants

**User Story:** As a developer, I want guaranteed structural properties of the output slice, so that I can rely on it for `slog` variadic arguments without runtime errors.

#### Acceptance criteria

1. THE System SHALL return a Flat_Slice whose length is always even (one key followed by one value for each field)
2. THE System SHALL ensure every element at an even index (0, 2, 4, ...) is a string representing the field key
3. THE System SHALL return the function signature `func ToLogger(in any) ([]any, error)`

### Requirement 6: testing requirements

**User Story:** As a developer, I want comprehensive tests for the new function, so that correctness is validated and regressions are caught early.

#### Acceptance criteria

1. THE System SHALL include an example test (`ExampleToLogger`) demonstrating typical usage
2. THE System SHALL include benchmark tests covering flat structs, nested structs, deep nesting, and multiple nested fields
3. THE System SHALL include a fuzz test (`FuzzToLogger`) exercising diverse inputs
4. THE System SHALL include property-based tests using `pgregory.net/rapid` validating correctness properties
5. THE System SHALL place all new tests in new test files without modifying existing test files
