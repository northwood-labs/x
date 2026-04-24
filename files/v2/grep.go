package files

import (
	"os"
	"regexp"
)

// GrepFile will search a file for a regular expression and return true if it is
// found. The regular expression is compiled and used to search the file
// contents.
func GrepFile(path, s string) (bool, error) {
	b, err := os.ReadFile(path) // lint:allow_possible_insecure
	if err != nil {
		return false, err
	}

	found, err := regexp.Match(s, b)
	if err != nil {
		return false, err
	}

	return found, nil
}
