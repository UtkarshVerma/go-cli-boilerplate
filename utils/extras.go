package utils

import "strconv"

// TypeOf returns the datatype of passed value.
func TypeOf(value string) string {
	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		return "int64"
	}
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return "float64"
	}
	if _, err := strconv.ParseBool(value); err == nil {
		return "bool"
	}
	return "string"
}
