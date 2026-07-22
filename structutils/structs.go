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
	"fmt"
	"reflect"
)

// ToMap converts a struct to a map[string]any.
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
