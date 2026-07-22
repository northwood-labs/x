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

package maputils

import "testing"

func TestMapToLoggerNilOnEmpty(t *testing.T) {
	emptyMap := make(map[string]any)
	result := MapToLogger(emptyMap)

	if result != nil {
		t.Errorf("expected nil for empty map, got %v", result)
	}
}

func TestMapToLoggerLength(t *testing.T) {
	input := map[string]any{
		"env":     "production",
		"version": 2,
		"debug":   false,
	}

	result := MapToLogger(input)

	if len(result) != len(input)*2 {
		t.Errorf("expected length %d, got %d", len(input)*2, len(result))
	}
}

func TestMapToLoggerKeysAndValues(t *testing.T) {
	input := map[string]any{
		"env":     "production",
		"version": 2,
		"debug":   false,
	}

	result := MapToLogger(input)

	// Build lookup maps from the alternating key-value result.
	keys := make(map[string]bool)
	values := make(map[any]bool)

	for i := 0; i < len(result); i += 2 {
		k, ok := result[i].(string)
		if !ok {
			t.Errorf("expected string key at index %d, got %T", i, result[i])
			continue
		}

		keys[k] = true
		values[result[i+1]] = true
	}

	for wantKey, wantVal := range input {
		if !keys[wantKey] {
			t.Errorf("key %q missing from result", wantKey)
		}

		if !values[wantVal] {
			t.Errorf("value %v (for key %q) missing from result", wantVal, wantKey)
		}
	}
}
