package utils

import "os"

// IsDir returns whether the argued path string
// is a valid and existing directory
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsDir()
}

// IsFile returns whether the argued path string
// is a valid and existing file
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}
