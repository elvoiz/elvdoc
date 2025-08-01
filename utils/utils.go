package utils

import (
	"path/filepath"
	"strings"
)

// IsElvFile checks if the given file path has a .elv extension.
func IsElvFile(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".elv")
}
