// Copyright 2026, Northwood Labs, LLC <license@northwood-labs.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package structutils

import (
	"reflect"
	"strings"
	"testing"

	"pgregory.net/rapid"
)

// assertLoggerPointerEquivalence checks that ToLogger produces
// identical results for a struct passed by value and by pointer.
func assertLoggerPointerEquivalence(t *rapid.T, input any) {
	byVal, err1 := ToLogger(input)
	if err1 != nil {
		t.Fatalf("by-value error: %v", err1)
	}

	// Create pointer via reflect for proper *T type.
	ptr := reflect.New(reflect.TypeOf(input))
	ptr.Elem().Set(reflect.ValueOf(input))

	byPtr, err2 := ToLogger(ptr.Interface())
	if err2 != nil {
		t.Fatalf("by-pointer error: %v", err2)
	}

	if len(byVal) != len(byPtr) {
		t.Fatalf(
			"length mismatch: val=%d ptr=%d",
			len(byVal), len(byPtr),
		)
	}

	for i := range byVal {
		if byVal[i] != byPtr[i] {
			t.Fatalf(
				"index %d: val=%v ptr=%v",
				i, byVal[i], byPtr[i],
			)
		}
	}
}

// TestToLogger_Property_PointerEquivalence verifies that ToLogger
// produces the same result for a struct passed by value and by
// pointer.
//
// **Validates: Requirements 1.2**.
func TestToLogger_Property_PointerEquivalence(t *testing.T) {
	t.Run("Settings", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			input := Settings{
				Verbose: rapid.Bool().Draw(t, "verbose"),
				Timeout: rapid.Int().Draw(t, "timeout"),
				Name:    rapid.String().Draw(t, "name"),
			}

			assertLoggerPointerEquivalence(t, input)
		})
	})

	t.Run("Order", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			input := Order{
				Customer: Customer{
					Name: rapid.String().Draw(t, "customerName"),
				},
				Address: Address{
					City:    rapid.String().Draw(t, "city"),
					Country: rapid.String().Draw(t, "country"),
				},
				Total: rapid.Int().Draw(t, "total"),
			}

			assertLoggerPointerEquivalence(t, input)
		})
	})
}

// TestToLogger_Property_NonStructReturnsError verifies that
// non-struct inputs always return a nil slice and a non-nil error.
//
// **Validates: Requirements 1.3**.
func TestToLogger_Property_NonStructReturnsError(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			input := rapid.Int().Draw(t, "input")

			result, err := ToLogger(input)
			if err == nil {
				t.Fatal("expected error for int input")
			}

			if result != nil {
				t.Fatalf("expected nil result, got %v", result)
			}
		})
	})

	t.Run("string", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			input := rapid.String().Draw(t, "input")

			result, err := ToLogger(input)
			if err == nil {
				t.Fatal("expected error for string input")
			}

			if result != nil {
				t.Fatalf("expected nil result, got %v", result)
			}
		})
	})

	t.Run("bool", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			input := rapid.Bool().Draw(t, "input")

			result, err := ToLogger(input)
			if err == nil {
				t.Fatal("expected error for bool input")
			}

			if result != nil {
				t.Fatalf("expected nil result, got %v", result)
			}
		})
	})
}

// TestToLogger_Property_EvenLengthInvariant verifies that the
// output slice length is always even for any valid struct input.
//
// **Validates: Requirements 5.1**.
func TestToLogger_Property_EvenLengthInvariant(t *testing.T) {
	t.Run("Settings", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			input := Settings{
				Verbose: rapid.Bool().Draw(t, "verbose"),
				Timeout: rapid.Int().Draw(t, "timeout"),
				Name:    rapid.String().Draw(t, "name"),
			}

			result, err := ToLogger(input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result)%2 != 0 {
				t.Fatalf("length must be even, got %d", len(result))
			}
		})
	})

	t.Run("Order", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			input := Order{
				Customer: Customer{
					Name: rapid.String().Draw(t, "customerName"),
				},
				Address: Address{
					City:    rapid.String().Draw(t, "city"),
					Country: rapid.String().Draw(t, "country"),
				},
				Total: rapid.Int().Draw(t, "total"),
			}

			result, err := ToLogger(input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result)%2 != 0 {
				t.Fatalf("length must be even, got %d", len(result))
			}
		})
	})
}

// TestToLogger_Property_KeysAreStrings verifies that every element
// at an even index in the output is a string.
//
// **Validates: Requirements 5.2**.
func TestToLogger_Property_KeysAreStrings(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		input := Order{
			Customer: Customer{
				Name: rapid.String().Draw(t, "customerName"),
			},
			Address: Address{
				City:    rapid.String().Draw(t, "city"),
				Country: rapid.String().Draw(t, "country"),
			},
			Total: rapid.Int().Draw(t, "total"),
		}

		result, err := ToLogger(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		for i := 0; i < len(result); i += 2 {
			if _, ok := result[i].(string); !ok {
				t.Fatalf(
					"index %d: expected string key, got %T",
					i, result[i],
				)
			}
		}
	})
}

// TestToLogger_Property_KeysSorted verifies that output keys are
// in non-decreasing alphabetical order.
//
// **Validates: Requirements 4.1, 4.2**.
func TestToLogger_Property_KeysSorted(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		input := Order{
			Customer: Customer{
				Name: rapid.String().Draw(t, "customerName"),
			},
			Address: Address{
				City:    rapid.String().Draw(t, "city"),
				Country: rapid.String().Draw(t, "country"),
			},
			Total: rapid.Int().Draw(t, "total"),
		}

		result, err := ToLogger(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var prevKey string

		for i := 0; i < len(result); i += 2 {
			key, ok := result[i].(string)
			if !ok {
				t.Fatalf(
					"index %d: expected string, got %T",
					i, result[i],
				)
			}

			if prevKey > key {
				t.Fatalf(
					"keys not sorted: %q > %q",
					prevKey, key,
				)
			}

			prevKey = key
		}
	})
}

// TestToLogger_Property_Deterministic verifies that calling
// ToLogger twice on the same input produces identical results.
//
// **Validates: Requirements 4.3**.
func TestToLogger_Property_Deterministic(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		input := Order{
			Customer: Customer{
				Name: rapid.String().Draw(t, "customerName"),
			},
			Address: Address{
				City:    rapid.String().Draw(t, "city"),
				Country: rapid.String().Draw(t, "country"),
			},
			Total: rapid.Int().Draw(t, "total"),
		}

		result1, err1 := ToLogger(input)
		if err1 != nil {
			t.Fatalf("first call error: %v", err1)
		}

		result2, err2 := ToLogger(input)
		if err2 != nil {
			t.Fatalf("second call error: %v", err2)
		}

		if len(result1) != len(result2) {
			t.Fatalf(
				"length mismatch: %d vs %d",
				len(result1), len(result2),
			)
		}

		for i := range result1 {
			if result1[i] != result2[i] {
				t.Fatalf(
					"index %d: %v != %v",
					i, result1[i], result2[i],
				)
			}
		}
	})
}

// TestToLogger_Property_NestedDotNotation verifies that nested
// struct fields produce keys containing dot separators with correct
// parent prefixes.
//
// **Validates: Requirements 2.1, 2.2**.
func TestToLogger_Property_NestedDotNotation(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		input := Company{
			CEO: Person{
				Home: Location{
					Coords: GPS{
						Lat: rapid.Float64().Draw(t, "lat"),
						Lon: rapid.Float64().Draw(t, "lon"),
					},
				},
			},
		}

		result, err := ToLogger(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Expected keys for Company -> CEO -> Home -> Coords -> GPS.
		expectedKeys := []string{
			"CEO.Home.Coords.Lat",
			"CEO.Home.Coords.Lon",
		}

		if len(result) != 4 {
			t.Fatalf("expected 4 elements, got %d", len(result))
		}

		for i, expected := range expectedKeys {
			key, ok := result[i*2].(string)
			if !ok {
				t.Fatalf("index %d: expected string", i*2)
			}

			if key != expected {
				t.Fatalf(
					"key %d: expected %q, got %q",
					i, expected, key,
				)
			}

			// Each key must contain dots.
			if !strings.Contains(key, ".") {
				t.Fatalf(
					"nested key %q has no dot separator",
					key,
				)
			}
		}
	})
}

// TestToLogger_Property_NilPointerEmission verifies that nil
// pointer-to-struct fields produce a key with nil value.
//
// **Validates: Requirements 2.4**.
func TestToLogger_Property_NilPointerEmission(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		name := rapid.String().Draw(t, "name")

		input := WithPointer{
			Name:    name,
			Details: nil,
		}

		result, err := ToLogger(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should have 4 elements: "Details", nil, "Name", value
		// (sorted alphabetically).
		if len(result) != 4 {
			t.Fatalf(
				"expected 4 elements, got %d: %v",
				len(result), result,
			)
		}

		// First pair should be "Details" with nil value.
		key, ok := result[0].(string)
		if !ok {
			t.Fatalf(
				"index 0: expected string, got %T",
				result[0],
			)
		}

		if key != "Details" {
			t.Fatalf("first key: expected Details, got %q", key)
		}

		if result[1] != nil {
			t.Fatalf("Details value: expected nil, got %v", result[1])
		}
	})
}

// TestToLogger_Property_ExportedOnlyFields verifies that no
// unexported field names appear in the output keys.
//
// **Validates: Requirements 3.1, 3.2**.
func TestToLogger_Property_ExportedOnlyFields(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		input := HasUnexported{
			Public: rapid.String().Draw(t, "public"),
			Inner:  Customer{Name: rapid.String().Draw(t, "name")},
		}

		result, err := ToLogger(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		for i := 0; i < len(result); i += 2 {
			key, ok := result[i].(string)
			if !ok {
				t.Fatalf("index %d: expected string", i)
			}

			// "hidden" is the unexported field — must not appear.
			if key == "hidden" ||
				strings.Contains(key, "hidden") {
				t.Fatalf("unexported field leaked: %q", key)
			}
		}
	})
}
