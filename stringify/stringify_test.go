// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
package stringify

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVariadicToStringArray(t *testing.T) {
	testCases := []struct {
		name           string
		input          []any
		expectedOutput []string
	}{
		{
			name:           "string",
			input:          []any{"a", "b"},
			expectedOutput: []string{"a", "b"},
		},
		{
			name:           "mixed",
			input:          []any{"a", 1},
			expectedOutput: []string{"a", "1"},
		},
		{
			name:           "empty",
			expectedOutput: []string{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := VariadicToStringArray(tc.input...)
			require.Equal(t, tc.expectedOutput, output)
		})
	}
}
