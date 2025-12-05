package generator

import (
	"strings"
	"unicode"
)

func toCamelCase(s string) string {
	// snake_case to camelCase
	// e.g. login_req -> loginReq

	var result strings.Builder
	nextUpper := false
	for i, r := range s {
		if r == '_' {
			nextUpper = true
			continue
		}
		if i == 0 {
			result.WriteRune(unicode.ToLower(r))
		} else if nextUpper {
			result.WriteRune(unicode.ToUpper(r))
			nextUpper = false
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func toPascalCase(s string) string {
	// snake_case to PascalCase
	// e.g. packet -> Packet, login_req -> LoginReq

	var result strings.Builder
	nextUpper := true
	for _, r := range s {
		if r == '_' {
			nextUpper = true
			continue
		}
		if nextUpper {
			result.WriteRune(unicode.ToUpper(r))
			nextUpper = false
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
