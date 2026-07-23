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
	"encoding/json"
	"fmt"
	"testing"
)

// ExampleToLogger demonstrates flat struct conversion to slog-ready
// key-value pairs, verifying that keys appear in sorted order.
func ExampleToLogger() {
	type Person struct {
		Name string
		Age  int
	}

	result, err := ToLogger(Person{Name: "Alice", Age: 30})
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	for i := 0; i < len(result); i += 2 {
		fmt.Printf("%s: %v\n", result[i], result[i+1])
	}

	// Output:
	// Age: 30
	// Name: Alice
}

// ExampleToLogger_nested demonstrates dot-notation flattening for
// nested structs, proving that hierarchy is preserved in key names.
func ExampleToLogger_nested() {
	type Address struct {
		City    string
		Country string
	}

	type Contact struct {
		Name    string
		Address Address
	}

	result, err := ToLogger(Contact{
		Name:    "Bob",
		Address: Address{City: "London", Country: "UK"},
	})
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	for i := 0; i < len(result); i += 2 {
		fmt.Printf("%s: %v\n", result[i], result[i+1])
	}

	// Output:
	// Address.City: London
	// Address.Country: UK
	// Name: Bob
}

// BenchmarkToLogger measures allocation and throughput for ToLogger
// across varying nesting depths, establishing a performance baseline.
func BenchmarkToLogger(b *testing.B) {
	b.ReportAllocs()

	benchCases := map[string]any{
		"flat_struct": Settings{
			Verbose: true, Timeout: 30, Name: "bench",
		},
		"nested_struct": UserInfo{
			Name: "q1mi", Age: 18,
			Profile: Profile{Hobby: "reading"},
		},
		"deep_nesting": Company{
			CEO: Person{
				Home: Location{
					Coords: GPS{Lat: 51.5, Lon: -0.1},
				},
			},
		},
		"multiple_nested": Order{
			Customer: Customer{Name: "Ada"},
			Address:  Address{City: "London", Country: "UK"},
			Total:    99,
		},
	}

	for name, input := range benchCases {
		b.Run(name, func(b *testing.B) {
			b.ResetTimer()

			for range b.N {
				_, _ = ToLogger(input) // lint:allow_unhandled
			}
		})
	}
}

// BenchmarkToLoggerParallel confirms ToLogger is goroutine-safe and
// measures contention under parallel workloads.
func BenchmarkToLoggerParallel(b *testing.B) { // lint:no_dupe
	b.ReportAllocs()

	type benchCase struct {
		input any
		name  string
	}

	cases := []benchCase{
		{name: "flat_struct", input: Settings{
			Verbose: true, Timeout: 30, Name: "bench",
		}},
		{name: "nested_struct", input: UserInfo{
			Name: "q1mi", Age: 18,
			Profile: Profile{Hobby: "reading"},
		}},
		{name: "deep_nesting", input: Company{
			CEO: Person{
				Home: Location{
					Coords: GPS{Lat: 51.5, Lon: -0.1},
				},
			},
		}},
		{name: "multiple_nested", input: Order{
			Customer: Customer{Name: "Ada"},
			Address:  Address{City: "London", Country: "UK"},
			Total:    99,
		}},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, _ = ToLogger(tc.input) // lint:allow_unhandled
				}
			})
		})
	}
}

// FuzzToLogger feeds randomly-generated JSON payloads through
// ToLogger to surface panics or invariant violations (e.g., odd-
// length output) on unexpected input shapes.
func FuzzToLogger(f *testing.F) {
	// Seed the corpus with JSON-encoded structs of varying shapes.
	seeds := []any{
		Settings{Verbose: true, Timeout: 30, Name: "test"},
		Settings{Verbose: false, Timeout: 0, Name: ""},
		Minimal{Value: "hello"},
		Numeric{X: 1, Y: 9.8, Z: -3},
		Mixed{A: "abc", B: 42, C: true, D: 2.71},
	}

	for _, s := range seeds {
		data, marshalErr := json.Marshal(s)
		if marshalErr != nil {
			continue
		}

		f.Add(string(data))
	}

	f.Fuzz(
		func(t *testing.T, jsonInput string) {
			// Decode into a Settings struct and run ToLogger.
			var s Settings

			decodeErr := json.Unmarshal([]byte(jsonInput), &s)
			if decodeErr != nil {
				return
			}

			result, err := ToLogger(s)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result == nil {
				t.Fatal("expected non-nil slice")
			}

			if len(result)%2 != 0 {
				t.Fatalf("result length must be even, got %d", len(result))
			}
		},
	)
}
