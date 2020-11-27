package utils

import "strconv"

// TypeOf returns the datatype of passed value.
func TypeOf(value string) string {
	if _, err := strconv.ParseBool(value); err == nil {
		return "bool"
	}
	if _, err := strconv.Atoi(value); err == nil {
		return "int"
	}
	return "string"
}
