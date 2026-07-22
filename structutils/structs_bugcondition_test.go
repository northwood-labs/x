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
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"pgregory.net/rapid"
)

// Validates: Requirements 1.1, 1.2, 1.3, 2.1, 2.2, 2.3, 2.4.
// Diverse struct types beyond UserInfo/Profile.
type (
	// Customer represents a customer with a name.
	Customer struct {
		Name string
	}

	// Address represents a postal address.
	Address struct {
		City    string
		Country string
	}

	// Order has multiple nested struct fields at the same level.
	Order struct {
		Customer Customer
		Address  Address
		Total    int
	}

	// GPS holds geographic coordinates.
	GPS struct {
		Lat float64
		Lon float64
	}

	// Location wraps GPS coordinates.
	Location struct {
		Coords GPS
	}

	// Person has a home location.
	Person struct {
		Home Location
	}

	// Company has a CEO who is a Person (4 levels deep).
	Company struct {
		CEO Person
	}

	// Inner is a simple struct for multi-level nesting.
	Inner struct {
		Val int
	}

	// Middle wraps Inner.
	Middle struct {
		C Inner
	}

	// Outer wraps Middle (3 levels: Outer -> Middle -> Inner).
	Outer struct {
		B Middle
	}

	// Node is a self-referencing struct via pointers.
	Node struct {
		Left  *Node
		Right *Node
		Value int
	}

	// WithPointer has a pointer-to-struct field.
	WithPointer struct {
		Details *Customer
		Name    string
	}

	// HasUnexported has both exported and unexported fields.
	HasUnexported struct {
		Public string
		hidden string // lint:allow_format
		Inner  Customer
	}

	// Profile is a test struct.
	Profile struct {
		Hobby string `json:"hobby"`
	}

	// UserInfo is a test struct.
	UserInfo struct {
		Profile `json:"profile"`

		Name string `json:"name"`
		Age  int    `json:"age"`
	}
)

// TestBugCondition_EmbeddedStruct tests the primary bug case from the
// design document.
func TestBugCondition_EmbeddedStruct(t *testing.T) {
	input := UserInfo{
		Name:    "q1mi",
		Age:     18,
		Profile: Profile{Hobby: "Two Color Ball"},
	}

	result, err := ToMap(input)
	assert.NoError(t, err)

	// The Profile field must be a map[string]any, not a raw struct.
	profileVal, ok := result["Profile"]
	assert.True(t, ok, "Profile key must exist in result")

	profileMap, isMap := profileVal.(map[string]any)
	assert.True(t, isMap,
		"Profile value must be map[string]any, got %T", profileVal)

	if isMap {
		assert.Equal(t, "Two Color Ball", profileMap["Hobby"])
	}

	// Inner field "Hobby" must NOT leak to the top-level map.
	_, hasHobby := result["Hobby"]
	assert.False(t, hasHobby,
		"inner field Hobby must not appear at top level")

	// Verify expected full result.
	expected := map[string]any{
		"Name":    "q1mi",
		"Age":     18,
		"Profile": map[string]any{"Hobby": "Two Color Ball"},
	}
	assert.Equal(t, expected, result)
}

// TestBugCondition_MultiLevelNesting tests 3-level nesting:
// Outer -> Middle -> Inner.
func TestBugCondition_MultiLevelNesting(t *testing.T) {
	input := Outer{B: Middle{C: Inner{Val: 1}}}

	result, err := ToMap(input)
	assert.NoError(t, err)

	expected := map[string]any{
		"B": map[string]any{
			"C": map[string]any{
				"Val": 1,
			},
		},
	}
	assert.Equal(t, expected, result)
}

// TestBugCondition_MultipleNestedFieldsSameLevel tests structs with
// multiple nested struct fields at the same level.
func TestBugCondition_MultipleNestedFieldsSameLevel(t *testing.T) {
	input := Order{
		Customer: Customer{Name: "Ada"},
		Address:  Address{City: "London", Country: "UK"},
		Total:    99,
	}

	result, err := ToMap(input)
	assert.NoError(t, err)

	expected := map[string]any{
		"Customer": map[string]any{"Name": "Ada"},
		"Address":  map[string]any{"City": "London", "Country": "UK"},
		"Total":    99,
	}
	assert.Equal(t, expected, result)

	// No inner fields should leak to the top level.
	_, hasName := result["Name"]
	assert.False(t, hasName,
		"inner field Name must not leak to top level")

	_, hasCity := result["City"]
	assert.False(t, hasCity,
		"inner field City must not leak to top level")
}

// TestBugCondition_DeeplyNested tests 4+ levels of nesting.
func TestBugCondition_DeeplyNested(t *testing.T) {
	input := Company{
		CEO: Person{
			Home: Location{
				Coords: GPS{Lat: 51.5, Lon: -0.1},
			},
		},
	}

	result, err := ToMap(input)
	assert.NoError(t, err)

	expected := map[string]any{
		"CEO": map[string]any{
			"Home": map[string]any{
				"Coords": map[string]any{
					"Lat": 51.5,
					"Lon": -0.1,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

// TestBugCondition_PointerToStruct tests that pointer-to-struct
// fields are dereferenced and recursively converted.
func TestBugCondition_PointerToStruct(t *testing.T) {
	input := WithPointer{
		Name:    "test",
		Details: &Customer{Name: "Alice"},
	}

	result, err := ToMap(input)
	assert.NoError(t, err)

	expected := map[string]any{
		"Name":    "test",
		"Details": map[string]any{"Name": "Alice"},
	}
	assert.Equal(t, expected, result)
}

// TestBugCondition_NilPointerToStruct tests nil pointer-to-struct
// fields are stored as nil without panic.
func TestBugCondition_NilPointerToStruct(t *testing.T) {
	input := Node{
		Value: 42,
		Left:  nil,
		Right: &Node{Value: 7, Left: nil, Right: nil},
	}

	result, err := ToMap(input)
	assert.NoError(t, err)

	expected := map[string]any{
		"Value": 42,
		"Left":  nil,
		"Right": map[string]any{
			"Value": 7,
			"Left":  nil,
			"Right": nil,
		},
	}
	assert.Equal(t, expected, result)
}

// TestBugCondition_UnexportedFields verifies no panic occurs and only
// exported fields appear in output.
func TestBugCondition_UnexportedFields(t *testing.T) {
	input := HasUnexported{
		Public: "visible",
		hidden: "invisible",
		Inner:  Customer{Name: "Bob"},
	}

	result, err := ToMap(input)
	assert.NoError(t, err)

	// Only exported fields should be in the result.
	_, hasHidden := result["hidden"]
	assert.False(t, hasHidden,
		"unexported field must not appear in result")

	// Public field should exist.
	assert.Equal(t, "visible", result["Public"])

	// Inner should be a map, not a raw struct.
	innerVal, ok := result["Inner"]
	assert.True(t, ok, "Inner key must exist")

	_, isMap := innerVal.(map[string]any)
	assert.True(t, isMap,
		"Inner value must be map[string]any, got %T", innerVal)
}

// assertStructFieldIsMap verifies that a struct-typed field in the
// result map is represented as map[string]any.
func assertStructFieldIsMap(
	t *rapid.T,
	result map[string]any,
	fieldName string,
) {
	val, exists := result[fieldName]
	if !exists {
		t.Fatalf("field %s missing from result", fieldName)
	}

	_, isMap := val.(map[string]any)
	if !isMap {
		t.Fatalf(
			"field %s: expected map[string]any, got %T",
			fieldName, val,
		)
	}
}

// assertKeyAbsent verifies that a key does not exist at the top
// level of the result map.
func assertKeyAbsent(
	t *rapid.T,
	result map[string]any,
	key string,
) {
	if _, exists := result[key]; exists {
		t.Fatalf("field %s leaked to top level", key)
	}
}

// TestBugCondition_Property_NestedStructsAreAlwaysMaps is a
// property-based test that verifies for diverse struct types that all
// struct-typed fields produce map[string]any values.
func TestBugCondition_Property_NestedStructsAreAlwaysMaps(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a random Order struct.
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

		result, err := ToMap(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Every field whose type is a struct must be a map[string]any.
		v := reflect.ValueOf(input)
		tp := v.Type()

		for i := range v.NumField() {
			field := v.Field(i)
			fieldType := tp.Field(i)

			if field.Kind() != reflect.Struct {
				continue
			}

			assertStructFieldIsMap(t, result, fieldType.Name)
		}

		// No inner field should leak to the top level.
		assertKeyAbsent(t, result, "Name")
		assertKeyAbsent(t, result, "City")
		assertKeyAbsent(t, result, "Country")
	})
}

// TestBugCondition_Property_DeeplyNestedMaps is a property-based test
// verifying deeply nested structs (Company -> Person -> Location ->
// GPS) always produce nested maps.
func TestBugCondition_Property_DeeplyNestedMaps(t *testing.T) {
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

		result, err := ToMap(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// CEO must be map[string]any.
		ceoVal, ok := result["CEO"]
		if !ok {
			t.Fatal("CEO key missing")
		}

		ceoMap, isMap := ceoVal.(map[string]any)
		if !isMap {
			t.Fatalf("CEO: expected map[string]any, got %T", ceoVal)
		}

		// CEO.Home must be map[string]any.
		homeVal, ok := ceoMap["Home"]
		if !ok {
			t.Fatal("Home key missing from CEO map")
		}

		homeMap, isMap := homeVal.(map[string]any)
		if !isMap {
			t.Fatalf(
				"Home: expected map[string]any, got %T", homeVal,
			)
		}

		// CEO.Home.Coords must be map[string]any.
		coordsVal, ok := homeMap["Coords"]
		if !ok {
			t.Fatal("Coords key missing from Home map")
		}

		_, isMap = coordsVal.(map[string]any)
		if !isMap {
			t.Fatalf(
				"Coords: expected map[string]any, got %T", coordsVal,
			)
		}

		// No deep field should leak to the top level.
		for _, key := range []string{
			"Home", "Coords", "Lat", "Lon",
		} {
			if _, exists := result[key]; exists {
				t.Fatalf("field %s leaked to top level", key)
			}
		}
	})
}

// TestBugCondition_Property_PointerToStructRecursion is a
// property-based test verifying pointer-to-struct fields are
// recursively converted.
func TestBugCondition_Property_PointerToStructRecursion(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		name := rapid.String().Draw(t, "name")
		detailName := rapid.String().Draw(t, "detailName")

		input := WithPointer{
			Name:    name,
			Details: &Customer{Name: detailName},
		}

		result, err := ToMap(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Details must be map[string]any.
		detailsVal, ok := result["Details"]
		if !ok {
			t.Fatal("Details key missing")
		}

		detailsMap, isMap := detailsVal.(map[string]any)
		if !isMap {
			t.Fatalf(
				"Details: expected map[string]any, got %T",
				detailsVal,
			)
		}

		// The inner Name should be inside the Details map.
		if detailsMap["Name"] != detailName {
			t.Fatalf(
				"Details.Name: expected %q, got %v",
				detailName, detailsMap["Name"],
			)
		}
	})
}
