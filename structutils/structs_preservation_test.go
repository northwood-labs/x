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

	"github.com/go-openapi/testify/v2/assert"
	"pgregory.net/rapid"
)

// **Validates: Requirements 3.1, 3.2, 3.3**.
type (
	// Settings has bool, int, and string fields.
	Settings struct {
		Name    string
		Timeout int
		Verbose bool
	}

	// Metrics has int, float64, and string fields.
	Metrics struct {
		Label string
		Count int
		Rate  float64
	}

	// Mixed has string, int, bool, and float64 fields.
	Mixed struct {
		A string
		B int
		C bool
		D float64
	}

	// Minimal has a single string field.
	Minimal struct {
		Value string
	}

	// Numeric has only numeric primitive fields.
	Numeric struct {
		X int
		Y float64
		Z int
	}
)

// TestPreservation_PrimitiveOnlyStructs verifies that primitive-only structs
// produce the expected flat map on unfixed code.
func TestPreservation_PrimitiveOnlyStructs(t *testing.T) {
	tests := map[string]struct {
		Input    any
		Expected map[string]any
	}{
		"Profile": {
			Input:    Profile{Hobby: "reading"},
			Expected: map[string]any{"Hobby": "reading"},
		},
		"Settings": {
			Input: Settings{
				Verbose: true, Timeout: 30, Name: "test",
			},
			Expected: map[string]any{
				"Verbose": true, "Timeout": 30, "Name": "test",
			},
		},
		"Metrics": {
			Input: Metrics{Count: 5, Rate: 3.14, Label: "ops"},
			Expected: map[string]any{
				"Count": 5, "Rate": 3.14, "Label": "ops",
			},
		},
		"Mixed": {
			Input: Mixed{A: "hello", B: 42, C: false, D: 2.71},
			Expected: map[string]any{
				"A": "hello", "B": 42, "C": false, "D": 2.71,
			},
		},
		"Minimal": {
			Input:    Minimal{Value: "only"},
			Expected: map[string]any{"Value": "only"},
		},
		"Numeric": {
			Input: Numeric{X: 1, Y: 9.8, Z: -3},
			Expected: map[string]any{
				"X": 1, "Y": 9.8, "Z": -3,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := ToMap(tc.Input)
			assert.NoError(t, err)
			assert.Equal(t, tc.Expected, result)
		})
	}
}

// TestPreservation_PointerDereference verifies that pointer-to-struct inputs
// produce the same result as by-value for diverse types.
func TestPreservation_PointerDereference(t *testing.T) {
	t.Run("Profile", func(t *testing.T) {
		byVal, err1 := ToMap(Profile{Hobby: "y"})
		assert.NoError(t, err1)

		byPtr, err2 := ToMap(&Profile{Hobby: "y"})
		assert.NoError(t, err2)
		assert.Equal(t, byVal, byPtr)
	})

	t.Run("Settings", func(t *testing.T) {
		s := Settings{Verbose: true, Timeout: 60, Name: "srv"}

		byVal, err1 := ToMap(s)
		assert.NoError(t, err1)

		byPtr, err2 := ToMap(&s)
		assert.NoError(t, err2)
		assert.Equal(t, byVal, byPtr)
	})

	t.Run("Metrics", func(t *testing.T) {
		m := Metrics{Count: 100, Rate: 0.5, Label: "rps"}

		byVal, err1 := ToMap(m)
		assert.NoError(t, err1)

		byPtr, err2 := ToMap(&m)
		assert.NoError(t, err2)
		assert.Equal(t, byVal, byPtr)
	})

	t.Run("Mixed", func(t *testing.T) {
		mx := Mixed{A: "z", B: -1, C: true, D: 0.0}

		byVal, err1 := ToMap(mx)
		assert.NoError(t, err1)

		byPtr, err2 := ToMap(&mx)
		assert.NoError(t, err2)
		assert.Equal(t, byVal, byPtr)
	})
}

// TestPreservation_NonStructError verifies the error path for
// non-struct inputs.
func TestPreservation_NonStructError(t *testing.T) {
	tests := map[string]struct {
		Input any
	}{
		"int":    {Input: 42},
		"string": {Input: "hello"},
		"slice":  {Input: []int{1, 2, 3}},
		"bool":   {Input: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := ToMap(tc.Input)
			assert.Nil(t, result)
			assert.Error(t, err)
			assert.True(
				t,
				strings.Contains(
					err.Error(),
					"ToFlattenedMap only accepts struct"+
						" or struct pointer; received",
				),
				"error message format mismatch: %q",
				err.Error(),
			)
		})
	}
}

// assertFlatMapMatches verifies that the ToMap result has the
// expected number of keys and each key matches its expected value.
func assertFlatMapMatches(
	t *rapid.T,
	result, expected map[string]any,
) {
	if len(result) != len(expected) {
		t.Fatalf(
			"expected %d keys, got %d: %v",
			len(expected), len(result), result,
		)
	}

	for k, want := range expected {
		if result[k] != want {
			t.Fatalf(
				"%s: expected %v, got %v",
				k, want, result[k],
			)
		}
	}
}

// TestPreservation_Property_PrimitiveStructsProduceFlatMap is a
// property-based test verifying that for all primitive-only structs,
// the output map contains exactly the field names as keys with their
// corresponding values.
func TestPreservation_Property_PrimitiveStructsProduceFlatMap(
	t *testing.T,
) {
	// Test with Settings struct.
	t.Run("Settings", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			input := Settings{
				Verbose: rapid.Bool().Draw(t, "verbose"),
				Timeout: rapid.Int().Draw(t, "timeout"),
				Name:    rapid.String().Draw(t, "name"),
			}

			result, err := ToMap(input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assertFlatMapMatches(t, result, map[string]any{
				"Verbose": input.Verbose,
				"Timeout": input.Timeout,
				"Name":    input.Name,
			})
		})
	})

	// Test with Metrics struct.
	t.Run("Metrics", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			input := Metrics{
				Count: rapid.Int().Draw(t, "count"),
				Rate:  rapid.Float64().Draw(t, "rate"),
				Label: rapid.String().Draw(t, "label"),
			}

			result, err := ToMap(input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assertFlatMapMatches(t, result, map[string]any{
				"Count": input.Count,
				"Rate":  input.Rate,
				"Label": input.Label,
			})
		})
	})

	// Test with Mixed struct.
	t.Run("Mixed", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			input := Mixed{
				A: rapid.String().Draw(t, "a"),
				B: rapid.Int().Draw(t, "b"),
				C: rapid.Bool().Draw(t, "c"),
				D: rapid.Float64().Draw(t, "d"),
			}

			result, err := ToMap(input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assertFlatMapMatches(t, result, map[string]any{
				"A": input.A,
				"B": input.B,
				"C": input.C,
				"D": input.D,
			})
		})
	})

	// Test with Minimal (single field).
	t.Run("Minimal", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			input := Minimal{
				Value: rapid.String().Draw(t, "value"),
			}

			result, err := ToMap(input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assertFlatMapMatches(t, result, map[string]any{
				"Value": input.Value,
			})
		})
	})

	// Test with Numeric (all numeric fields).
	t.Run("Numeric", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			input := Numeric{
				X: rapid.Int().Draw(t, "x"),
				Y: rapid.Float64().Draw(t, "y"),
				Z: rapid.Int().Draw(t, "z"),
			}

			result, err := ToMap(input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assertFlatMapMatches(t, result, map[string]any{
				"X": input.X,
				"Y": input.Y,
				"Z": input.Z,
			})
		})
	})
}

// assertPointerMatchesByValue verifies that ToMap produces the same
// result for a struct passed by value and by pointer.
func assertPointerMatchesByValue(t *rapid.T, input any) {
	byVal, err1 := ToMap(input)
	if err1 != nil {
		t.Fatalf("by-value error: %v", err1)
	}

	// Create a pointer to a copy of input via reflect so that the
	// pointer is *T (not *interface{}).
	ptr := reflect.New(reflect.TypeOf(input))
	ptr.Elem().Set(reflect.ValueOf(input))

	byPtr, err2 := ToMap(ptr.Interface())
	if err2 != nil {
		t.Fatalf("by-pointer error: %v", err2)
	}

	if len(byVal) != len(byPtr) {
		t.Fatalf(
			"length mismatch: val=%d ptr=%d",
			len(byVal), len(byPtr),
		)
	}

	for k, v := range byVal {
		if byPtr[k] != v {
			t.Fatalf(
				"key %s: val=%v ptr=%v",
				k, v, byPtr[k],
			)
		}
	}
}

// TestPreservation_Property_PointerMatchesByValue is a property-based
// test verifying that pointer-to-primitive-only-struct inputs produce
// the same result as the non-pointer version across diverse types.
func TestPreservation_Property_PointerMatchesByValue(t *testing.T) {
	t.Run("Settings", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			input := Settings{
				Verbose: rapid.Bool().Draw(t, "verbose"),
				Timeout: rapid.Int().Draw(t, "timeout"),
				Name:    rapid.String().Draw(t, "name"),
			}

			assertPointerMatchesByValue(t, input)
		})
	})

	t.Run("Metrics", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			input := Metrics{
				Count: rapid.Int().Draw(t, "count"),
				Rate:  rapid.Float64().Draw(t, "rate"),
				Label: rapid.String().Draw(t, "label"),
			}

			assertPointerMatchesByValue(t, input)
		})
	})

	t.Run("Mixed", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			input := Mixed{
				A: rapid.String().Draw(t, "a"),
				B: rapid.Int().Draw(t, "b"),
				C: rapid.Bool().Draw(t, "c"),
				D: rapid.Float64().Draw(t, "d"),
			}

			assertPointerMatchesByValue(t, input)
		})
	})

	t.Run("Numeric", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			input := Numeric{
				X: rapid.Int().Draw(t, "x"),
				Y: rapid.Float64().Draw(t, "y"),
				Z: rapid.Int().Draw(t, "z"),
			}

			assertPointerMatchesByValue(t, input)
		})
	})
}

// assertNonStructReturnsError verifies that ToMap returns nil result
// and an error matching the expected format for non-struct inputs.
func assertNonStructReturnsError(t *rapid.T, input any, label string) {
	result, err := ToMap(input)
	if err == nil {
		t.Fatalf("expected error for %s input", label)
	}

	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}

	if !strings.Contains(
		err.Error(),
		"ToFlattenedMap only accepts struct"+
			" or struct pointer; received",
	) {
		t.Fatalf("unexpected error format: %q", err.Error())
	}
}

// TestPreservation_Property_NonStructReturnsError is a property-based
// test verifying that non-struct inputs always return an error
// matching the expected format.
func TestPreservation_Property_NonStructReturnsError(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			input := rapid.Int().Draw(t, "input")
			assertNonStructReturnsError(t, input, "int")
		})
	})

	t.Run("string", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			input := rapid.String().Draw(t, "input")
			assertNonStructReturnsError(t, input, "string")
		})
	})

	t.Run("bool", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			input := rapid.Bool().Draw(t, "input")
			assertNonStructReturnsError(t, input, "bool")
		})
	})

	t.Run("float64", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			input := rapid.Float64().Draw(t, "input")
			assertNonStructReturnsError(t, input, "float64")
		})
	})
}
