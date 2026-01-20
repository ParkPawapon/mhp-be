package utils

import "strings"

func MaskCitizenID(id string) string {
	if id == "" {
		return ""
	}
	if len(id) <= 4 {
		return strings.Repeat("*", len(id))
	}
	maskLen := len(id) - 4
	return strings.Repeat("*", maskLen) + id[len(id)-4:]
}
