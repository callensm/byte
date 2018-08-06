package utils

import (
	"os"
	"regexp"
)

var ignore = []string{
	"^\\.git$",
	"^node_modules$",
	"^vendor$",
	"^.*\\.db$",
	"^.*\\.lock$",
}

// IsIgnored returns whether or not the file matches
// one of the ignored file regular expressions to skip
func IsIgnored(name string) bool {
	for _, i := range ignore {
		if ok, err := regexp.MatchString(i, name); ok && err == nil {
			return true
		} else if err != nil {
			Catch(err)
		}
	}
	return false
}

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
