package utils

import (
	"strings"
	"unicode"
)

// ToCamelCase converts `in` to PascalCase
func ToCamelCase(in string) (out string) {
	for i, token := range splitTokens(in) {
		if i != 0 {
			token = strings.Title(token)
		}
		out += token
	}
	return
}

// ToPascalCase converts `in` to PascalCase
func ToPascalCase(in string) (out string) {
	for _, token := range splitTokens(in) {
		out += strings.Title(token)
	}
	return
}

// ToKebabCase converts `in` to kebab-case
func ToKebabCase(in string) string {
	return strings.Join(splitTokens(in), "-")
}

// ToSnakeCase converts `in` to kebab-case
func ToSnakeCase(in string) string {
	return strings.Join(splitTokens(in), "_")
}

func splitTokens(in string) (out []string) {
	if strings.Contains(in, "_") { // snake_case
		out = strings.Split(in, "_")
	} else if strings.Contains(in, "-") { // kebab-case
		out = strings.Split(in, "-")
	} else { // PascalCase or camelCase
		indices := []int{0}
		for i := range in {
			if i < len(in)-1 {
				if unicode.IsLower(rune(in[i])) && unicode.IsUpper(rune(in[i+1])) {
					indices = append(indices, i+1)
				}
			}
		}
		indices = append(indices, len(in))

		for i := 0; i < len(indices)-1; i++ {
			beg, end := indices[i], indices[i+1]
			out = append(out, strings.ToLower(in[beg:end]))
		}
	}

	return
}
