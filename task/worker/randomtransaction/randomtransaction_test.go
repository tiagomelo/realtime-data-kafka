// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
package randomtransaction

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tiagomelo/realtime-data-kafka/stringify"
)

func TestWork(t *testing.T) {
	testCases := []struct {
		name                string
		mockPrintToLog      func(log *log.Logger, v ...any)
		mockOpenFile        func(name string, flag int, perm fs.FileMode) (*os.File, error)
		mockJsonMarshal     func(v any) ([]byte, error)
		mockFileWriteString func(file *os.File, s string) (n int, err error)
	}{
		{
			name:           "happy path",
			mockPrintToLog: func(log *log.Logger, v ...any) {},
			mockOpenFile: func(name string, flag int, perm fs.FileMode) (*os.File, error) {
				return new(os.File), nil
			},
			mockJsonMarshal: func(v any) ([]byte, error) {
				return []byte(""), nil
			},
			mockFileWriteString: func(file *os.File, s string) (n int, err error) {
				return 0, nil
			},
		},
		{
			name: "error when opening file",
			mockOpenFile: func(name string, flag int, perm fs.FileMode) (*os.File, error) {
				return nil, errors.New("random error")
			},
			mockPrintToLog: func(log *log.Logger, v ...any) {
				expectedMsg := []string{"[error opening file: random error]"}
				c := stringify.VariadicToStringArray(v)
				require.Equal(t, expectedMsg, c)
			},
		},
		{
			name: "error when marshaling",
			mockOpenFile: func(name string, flag int, perm fs.FileMode) (*os.File, error) {
				return new(os.File), nil
			},
			mockJsonMarshal: func(v any) ([]byte, error) {
				return nil, errors.New("random error")
			},
			mockPrintToLog: func(log *log.Logger, v ...any) {
				expectedMsg := []string{"[error marshalling json: random error]"}
				c := stringify.VariadicToStringArray(v)
				require.Equal(t, expectedMsg, c)
			},
		},
		{
			name: "error when writing to file",
			mockOpenFile: func(name string, flag int, perm fs.FileMode) (*os.File, error) {
				return new(os.File), nil
			},
			mockJsonMarshal: func(v any) ([]byte, error) {
				return []byte(""), nil
			},
			mockFileWriteString: func(file *os.File, s string) (n int, err error) {
				return 0, errors.New("random error")
			},
			mockPrintToLog: func(log *log.Logger, v ...any) {
				expectedMsg := []string{"[error writing to file: random error]"}
				c := stringify.VariadicToStringArray(v)
				require.Equal(t, expectedMsg, c)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			printToLog = tc.mockPrintToLog
			openFile = tc.mockOpenFile
			jsonMarshal = tc.mockJsonMarshal
			fileWriteString = tc.mockFileWriteString
			worker := &Worker{
				MinAmount: 10,
				MaxAmount: 100,
				FilePath:  "filepath",
			}
			worker.Work(context.TODO())
		})
	}
}
