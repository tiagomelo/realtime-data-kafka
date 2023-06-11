package stringify

import "fmt"

// VariadicToStringArray converts variadic arguments to a string array.
func VariadicToStringArray(v ...any) []string {
	result := make([]string, len(v))
	for i, arg := range v {
		switch v := arg.(type) {
		case string:
			result[i] = v
		default:
			result[i] = fmt.Sprint(v)
		}
	}
	return result
}
