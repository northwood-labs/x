package grammar

// Pluralize allows you to pass an amount, and a singular + plural form of a
// word. The correct form of the word will be used. U.S. English-only.
func Pluralize(amount int, singular, plural string) string {
	if amount == 1 {
		return singular
	}

	return plural
}
