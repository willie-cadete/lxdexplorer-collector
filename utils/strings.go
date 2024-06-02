package utils

import "strings"

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	// strings.Title is deprecated in Go 1.18
	return strings.ToUpper(s[:1]) + s[1:]
}

func Lowercase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s)
}
