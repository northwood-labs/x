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

// Adapted from https://www.sobyte.net/post/2021-06/several-ways-to-convert-struct-to-mapstringinterface/

package structutils

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// pair holds a single key-value entry collected during recursive
// traversal. Collecting pairs first and sorting afterward produces
// deterministic output without requiring the traversal itself to
// visit fields in alphabetical order.
type pair struct {
	val any
	key string
}

// ToMap converts a struct (or pointer-to-struct) into a nested
// map[string]any, recursing into struct-typed fields. This is the
// building block for serializing arbitrary structs to JSON-like
// formats without manually writing a conversion for each type.
// Accepting a pointer transparently covers the common case where
// callers already hold a *T.
func ToMap(in any) (map[string]any, error) {
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Pointer { // Structure Pointer.
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf( // lint:allow_errorf
			"ToFlattenedMap only accepts struct or struct pointer; received %T",
			v,
		)
	}

	return structToMap(v), nil
}

// ToLogger converts a struct into a flat []any slice of alternating
// key-value pairs suitable for slog structured logging. Nested struct
// fields are flattened using dot notation (e.g., "Address.City") so
// that deeply structured data can be logged without losing hierarchy
// information. Keys are sorted alphabetically to guarantee
// deterministic log output — important for log-based alerting and
// diff-friendly test assertions.
func ToLogger(in any) ([]any, error) {
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf( // lint:allow_errorf
			"ToLogger only accepts struct or struct pointer; received %T",
			in,
		)
	}

	pairs := structToLogger(v, "")

	// Sort pairs by key to produce deterministic output regardless of
	// struct field declaration order or reflection iteration quirks.
	slices.SortFunc(pairs, func(a, b pair) int {
		return strings.Compare(a.key, b.key)
	})

	// Pre-allocate the result slice at exactly 2× the pair count to
	// avoid incremental growth allocations.
	result := make([]any, 0, len(pairs)*2)
	for _, p := range pairs {
		result = append(result, p.key, p.val)
	}

	return result, nil
}

// structToMap recursively walks a struct's exported fields and
// converts them into a map. Struct-typed fields become nested maps,
// preserving the hierarchy. Pointer fields are dereferenced (possibly
// multiple levels) so that *T and T produce equivalent output.
func structToMap(v reflect.Value) map[string]any {
	out := make(map[string]any)
	t := v.Type()

	for i := range v.NumField() {
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

		// After nil-pointer handling, check whether nil was stored
		// via the break above (field remains Pointer kind if nil).
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

// structToLogger recursively collects all exported leaf fields as
// key-value pairs. Struct-typed fields are not emitted directly;
// instead the function recurses, prepending the parent field name
// and a dot to produce fully-qualified keys like "Address.City".
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
		nilStored := false

		// Dereference pointers (possibly multiple levels).
		for field.Kind() == reflect.Pointer {
			if field.IsNil() {
				pairs = append(pairs, pair{key: key, val: nil})
				nilStored = true

				break
			}

			field = field.Elem()
		}

		if nilStored {
			continue
		}

		// Skip invalid values (from nil pointer deref edge cases).
		if field.Kind() == reflect.Invalid {
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
