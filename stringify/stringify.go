// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
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
