// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadConfig(t *testing.T) {
	testCases := []struct {
		name                   string
		mockedGodotenvLoad     func(filenames ...string) (err error)
		mockedEnvconfigProcess func(prefix string, spec interface{}) error
		expectedError          error
	}{
		{
			name: "happy path",
			mockedGodotenvLoad: func(filenames ...string) (err error) {
				return nil
			},
			mockedEnvconfigProcess: func(prefix string, spec interface{}) error {
				return nil
			},
		},
		{
			name: "error loading env vars",
			mockedGodotenvLoad: func(filenames ...string) (err error) {
				return errors.New("random error")
			},
			expectedError: errors.New("loading env vars: random error"),
		},
		{
			name: "error processing env vars",
			mockedGodotenvLoad: func(filenames ...string) (err error) {
				return nil
			},
			mockedEnvconfigProcess: func(prefix string, spec interface{}) error {
				return errors.New("random error")
			},
			expectedError: errors.New("processing env vars: random error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			godotenvLoad = tc.mockedGodotenvLoad
			envconfigProcess = tc.mockedEnvconfigProcess
			config, err := Read(".env")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf(`expected no error to occur, got "%v"`, err)
				}
				require.Nil(t, config)
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf(`expected error to occur, got nil`)
				}
				require.NotNil(t, config)
			}
		})
	}
}
