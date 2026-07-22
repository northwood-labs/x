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

// Forked from https://www.sobyte.net/post/2021-06/several-ways-to-convert-struct-to-mapstringinterface/

package structutils

import (
	"encoding/json"
	"fmt"
	"testing"
)

func ExampleToMap() {
	type Person struct {
		Name string
		Age  int
	}

	result, err := ToMap(Person{Name: "Alice", Age: 30})
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println("Name:", result["Name"])
	fmt.Println("Age:", result["Age"])
	// Output:
	// Name: Alice
	// Age: 30
}

func BenchmarkToMap(b *testing.B) {
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
				_, _ = ToMap(input) // lint:allow_unhandled
			}
		})
	}
}

func BenchmarkToMapParallel(b *testing.B) {
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
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, _ = ToMap(input) // lint:allow_unhandled
				}
			})
		})
	}
}

func FuzzToMap(f *testing.F) {
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
			// Decode into a Settings struct and run ToMap.
			var s Settings

			decodeErr := json.Unmarshal([]byte(jsonInput), &s)
			if decodeErr != nil {
				return
			}

			result, err := ToMap(s)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result == nil {
				t.Fatal("expected non-nil map")
			}
		},
	)
}
