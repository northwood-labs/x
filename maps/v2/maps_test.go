package maps

import "testing"

func TestMapToLoggerNilOnEmpty(t *testing.T) {
	result := MapToLogger(map[string]any{})

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
