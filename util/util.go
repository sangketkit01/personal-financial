package util

import "strings"

func CapitalizeWord(s string) string {
	if len(s) == 0 {
		return s
	}
	s = strings.ToLower(s)
	return strings.ToUpper(s[:1]) + s[1:]
}
