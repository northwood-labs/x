# Bugfix Requirements Document

## Introduction

The `ToFlattenedMap` function in `structutils/structs.go` incorrectly handles nested and embedded structs. When a struct contains an embedded struct field, the function stores the raw struct value in the output map rather than recursively converting it to `map[string]any`. Additionally, the function unconditionally flattens embedded struct fields into the top-level map, causing both the raw struct and its inner fields to appear at the top level. The fix should recursively convert nested structs into `map[string]any` values and avoid flattening inner fields into the parent map.

The function is reflection-based and generic — it must work correctly with any arbitrary struct type at any nesting depth, not only the specific types used in the test file (e.g., `UserInfo`, `Profile`). Those test types are merely examples; the implementation must handle all valid Go struct compositions.

## Bug analysis

### Current behavior (defect)

1.1 WHEN a struct contains an embedded struct field (e.g., `Profile` embedded in `UserInfo`) THEN the system stores the raw struct value (e.g., `structutils.Profile{Hobby:"Two Color Ball"}`) as the map value instead of a recursively converted `map[string]any`.

1.2 WHEN a struct contains an embedded struct field THEN the system flattens the inner struct's fields (e.g., `Hobby`) into the top-level map alongside the parent struct's own fields.

1.3 WHEN the loop processes a struct-typed field THEN the system unconditionally executes the "General Fields" assignment, adding both the raw struct value and its individual fields to the output map.

### Expected behavior (correct)

2.1 WHEN a struct contains an embedded struct field (e.g., `Profile` embedded in `UserInfo`) THEN the system SHALL recursively convert that embedded struct to a `map[string]any` and store the resulting map as the value under the field's name.

2.2 WHEN a struct contains an embedded struct field THEN the system SHALL NOT flatten the inner struct's fields into the top-level map; they SHALL remain nested inside their parent key as a `map[string]any`.

2.3 WHEN the loop processes a struct-typed field THEN the system SHALL skip the "General Fields" assignment for that field, since the recursive conversion handles it.

2.4 WHEN any arbitrary struct type is passed (at any nesting depth, with any combination of primitive fields, nested struct fields, and pointer-to-struct fields) THEN the system SHALL correctly and recursively convert all nested structs to `map[string]any` — the function is fully generic via reflection and is not limited to specific struct types used in tests.

### Unchanged behavior (regression prevention)

3.1 WHEN a struct contains only primitive fields (e.g., `string`, `int`, `bool`) THEN the system SHALL CONTINUE TO store each field's value directly in the output map keyed by the field name.

3.2 WHEN the input is a pointer to a struct THEN the system SHALL CONTINUE TO dereference the pointer and produce the same map output as if the struct were passed by value.

3.3 WHEN the input is not a struct or struct pointer THEN the system SHALL CONTINUE TO return an error indicating that only struct or struct pointer inputs are accepted.

---

### Bug condition (formal)

```pascal
FUNCTION isBugCondition(X)
  INPUT: X of type any (input to ToFlattenedMap)
  OUTPUT: boolean

  // Returns true when the struct contains at least one field whose
  // type is itself a struct (embedded or named).
  // This applies to ANY arbitrary struct type — not limited to specific
  // types like UserInfo or Profile. The function uses reflection and must
  // handle all valid Go struct compositions at any nesting depth.
  RETURN X is a struct AND X has at least one field of kind reflect.Struct
END FUNCTION
```

### Property specification (fix checking)

```pascal
// Property: Fix Checking - Nested struct fields are recursively converted
// This property holds for ANY arbitrary struct type, not just test-specific
// types. The function is reflection-based and must handle:
// - Structs nested N levels deep
// - Structs with mixed field types (primitives + nested structs + pointers)
// - Structs with multiple nested struct fields at the same level
// - Any user-defined struct composition
FOR ALL X of any struct type WHERE isBugCondition(X) DO
  result <- ToFlattenedMap'(X)
  FOR EACH field F in X WHERE F.Kind == reflect.Struct DO
    ASSERT result[F.Name] is of type map[string]any
    ASSERT result[F.Name] == ToFlattenedMap'(F.Value)
    ASSERT no field from F appears as a top-level key in result
  END FOR
END FOR
```

### Preservation checking

```pascal
// Property: Preservation Checking - Non-struct fields are unchanged
// Applies to any arbitrary struct type without nested struct fields.
FOR ALL X of any struct type WHERE NOT isBugCondition(X) DO
  ASSERT ToFlattenedMap(X) = ToFlattenedMap'(X)
END FOR
```

This ensures that for all inputs without nested struct fields (regardless of struct type), the fixed code behaves identically to the original.
