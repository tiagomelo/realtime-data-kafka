// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
package screen

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/stretchr/testify/require"
	"github.com/tiagomelo/realtime-data-kafka/stats"
)

var (
	originalPtermDefaultAreaStart    = ptermDefaultAreaStart
	originalPtermDefaultCenterSprint = ptermDefaultCenterSprint
	originalStopAreaPrinter          = stopAreaPrinter
	originalLayoutSprint             = layoutSprint
	originalUpdateAreaPrinter        = updateAreaPrinter
)

func TestNewKafkaConsumerScreen(t *testing.T) {
	testCases := []struct {
		name                      string
		mockPtermDefaultAreaStart func(text ...interface{}) (*pterm.AreaPrinter, error)
		expectedError             error
	}{
		{
			name: "happy path",
			mockPtermDefaultAreaStart: func(text ...interface{}) (*pterm.AreaPrinter, error) {
				return &pterm.AreaPrinter{}, nil
			},
		},
		{
			name: "error",
			mockPtermDefaultAreaStart: func(text ...interface{}) (*pterm.AreaPrinter, error) {
				return nil, errors.New("random error")
			},
			expectedError: errors.New("starting printer: random error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ptermDefaultAreaStart = tc.mockPtermDefaultAreaStart
			screen, err := NewKafkaConsumerScreen(&stats.KafkaConsumerStats{})
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf(`expected "%v" error, got nil`, err)
				}
				require.NotNil(t, screen)
			}
		})
	}
	ptermDefaultAreaStart = originalPtermDefaultAreaStart
	ptermDefaultCenterSprint = originalPtermDefaultCenterSprint
	stopAreaPrinter = originalStopAreaPrinter
	layoutSprint = originalLayoutSprint
	updateAreaPrinter = originalUpdateAreaPrinter
}

func TestKafkaConsumerUpdateContent(t *testing.T) {
	testCases := []struct {
		name                         string
		finalUpdate                  bool
		mockPtermDefaultCenterSprint func(a ...interface{}) string
		mockLayoutSprint             func(layout *pterm.CenterPrinter, a ...interface{}) string
		mockUpdateAreaPrinter        func(areaPrinter *pterm.AreaPrinter, text ...interface{})
		mockStopAreaPrinter          func(areaPrinter *pterm.AreaPrinter) error
		expectedError                error
	}{
		{
			name: "happy path when it is not the final update",
			mockPtermDefaultCenterSprint: func(a ...interface{}) string {
				return "some string"
			},
			mockLayoutSprint: func(layout *pterm.CenterPrinter, a ...interface{}) string {
				return "some string"
			},
			mockUpdateAreaPrinter: func(areaPrinter *pterm.AreaPrinter, text ...interface{}) {},
		},
		{
			name: "happy path when it is the final update",
			mockPtermDefaultCenterSprint: func(a ...interface{}) string {
				return "some string"
			},
			mockLayoutSprint: func(layout *pterm.CenterPrinter, a ...interface{}) string {
				return "some string"
			},
			mockUpdateAreaPrinter: func(areaPrinter *pterm.AreaPrinter, text ...interface{}) {},
			mockStopAreaPrinter: func(areaPrinter *pterm.AreaPrinter) error {
				return nil
			},
			finalUpdate: true,
		},
		{
			name: "error when it is the final update",
			mockPtermDefaultCenterSprint: func(a ...interface{}) string {
				return "some string"
			},
			mockLayoutSprint: func(layout *pterm.CenterPrinter, a ...interface{}) string {
				return "some string"
			},
			mockUpdateAreaPrinter: func(areaPrinter *pterm.AreaPrinter, text ...interface{}) {},
			mockStopAreaPrinter: func(areaPrinter *pterm.AreaPrinter) error {
				return errors.New("random error")
			},
			finalUpdate:   true,
			expectedError: errors.New("stopping printer: random error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ptermDefaultCenterSprint = tc.mockPtermDefaultCenterSprint
			layoutSprint = tc.mockLayoutSprint
			updateAreaPrinter = tc.mockUpdateAreaPrinter
			stopAreaPrinter = tc.mockStopAreaPrinter
			screen, err := NewKafkaConsumerScreen(&stats.KafkaConsumerStats{})
			require.Nil(t, err)
			err = screen.UpdateContent(tc.finalUpdate)
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf(`expected "%v" error, got nil`, err)
				}
				require.NotNil(t, screen)
			}
		})
	}
	ptermDefaultAreaStart = originalPtermDefaultAreaStart
	ptermDefaultCenterSprint = originalPtermDefaultCenterSprint
	stopAreaPrinter = originalStopAreaPrinter
	layoutSprint = originalLayoutSprint
	updateAreaPrinter = originalUpdateAreaPrinter
}

func TestNewKafkaProducerScreen(t *testing.T) {
	testCases := []struct {
		name                      string
		mockPtermDefaultAreaStart func(text ...interface{}) (*pterm.AreaPrinter, error)
		expectedError             error
	}{
		{
			name: "happy path",
			mockPtermDefaultAreaStart: func(text ...interface{}) (*pterm.AreaPrinter, error) {
				return &pterm.AreaPrinter{}, nil
			},
		},
		{
			name: "error",
			mockPtermDefaultAreaStart: func(text ...interface{}) (*pterm.AreaPrinter, error) {
				return nil, errors.New("random error")
			},
			expectedError: errors.New("starting printer: random error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ptermDefaultAreaStart = tc.mockPtermDefaultAreaStart
			screen, err := NewKafkaProducerScreen(&stats.KafkaProducerStats{})
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf(`expected "%v" error, got nil`, err)
				}
				require.NotNil(t, screen)
			}
		})
	}
	ptermDefaultAreaStart = originalPtermDefaultAreaStart
	ptermDefaultCenterSprint = originalPtermDefaultCenterSprint
	stopAreaPrinter = originalStopAreaPrinter
	layoutSprint = originalLayoutSprint
	updateAreaPrinter = originalUpdateAreaPrinter
}

func TestKafkaProducerUpdateContent(t *testing.T) {
	testCases := []struct {
		name                         string
		finalUpdate                  bool
		mockPtermDefaultCenterSprint func(a ...interface{}) string
		mockLayoutSprint             func(layout *pterm.CenterPrinter, a ...interface{}) string
		mockUpdateAreaPrinter        func(areaPrinter *pterm.AreaPrinter, text ...interface{})
		mockStopAreaPrinter          func(areaPrinter *pterm.AreaPrinter) error
		expectedError                error
	}{
		{
			name: "happy path when it is not the final update",
			mockPtermDefaultCenterSprint: func(a ...interface{}) string {
				return "some string"
			},
			mockLayoutSprint: func(layout *pterm.CenterPrinter, a ...interface{}) string {
				return "some string"
			},
			mockUpdateAreaPrinter: func(areaPrinter *pterm.AreaPrinter, text ...interface{}) {},
		},
		{
			name: "happy path when it is the final update",
			mockPtermDefaultCenterSprint: func(a ...interface{}) string {
				return "some string"
			},
			mockLayoutSprint: func(layout *pterm.CenterPrinter, a ...interface{}) string {
				return "some string"
			},
			mockUpdateAreaPrinter: func(areaPrinter *pterm.AreaPrinter, text ...interface{}) {},
			mockStopAreaPrinter: func(areaPrinter *pterm.AreaPrinter) error {
				return nil
			},
			finalUpdate: true,
		},
		{
			name: "error when it is the final update",
			mockPtermDefaultCenterSprint: func(a ...interface{}) string {
				return "some string"
			},
			mockLayoutSprint: func(layout *pterm.CenterPrinter, a ...interface{}) string {
				return "some string"
			},
			mockUpdateAreaPrinter: func(areaPrinter *pterm.AreaPrinter, text ...interface{}) {},
			mockStopAreaPrinter: func(areaPrinter *pterm.AreaPrinter) error {
				return errors.New("random error")
			},
			finalUpdate:   true,
			expectedError: errors.New("stopping printer: random error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ptermDefaultCenterSprint = tc.mockPtermDefaultCenterSprint
			layoutSprint = tc.mockLayoutSprint
			updateAreaPrinter = tc.mockUpdateAreaPrinter
			stopAreaPrinter = tc.mockStopAreaPrinter
			screen, err := NewKafkaProducerScreen(&stats.KafkaProducerStats{})
			require.Nil(t, err)
			err = screen.UpdateContent(tc.finalUpdate)
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf(`expected "%v" error, got nil`, err)
				}
				require.NotNil(t, screen)
			}
		})
	}
	ptermDefaultAreaStart = originalPtermDefaultAreaStart
	ptermDefaultCenterSprint = originalPtermDefaultCenterSprint
	stopAreaPrinter = originalStopAreaPrinter
	layoutSprint = originalLayoutSprint
	updateAreaPrinter = originalUpdateAreaPrinter
}
