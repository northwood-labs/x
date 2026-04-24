package files

import "os"

// FileExists checks if a file exists at the given path. It returns true if the
// file exists and is not a directory, and false otherwise.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
