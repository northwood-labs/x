package maps

// MapToLogger converts a map[string]any to a slice of alternating keys and values
// compatible with slog's key-value logging format.
func MapToLogger(m map[string]any) []any {
	if len(m) == 0 {
		return nil
	}

	result := make([]any, 0, len(m)*2)
	for k, v := range m {
		result = append(result, k, v)
	}

	return result
}
