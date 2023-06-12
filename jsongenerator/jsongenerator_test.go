// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tiagomelo/realtime-data-kafka/transaction"
)

const fileName = "data.txt"

func cleanup() error {
	if err := os.Remove(fileName); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func TestMain(m *testing.M) {
	if err := cleanup(); err != nil {
		fmt.Printf("error when deleting %s: %v\n", fileName, err)
		os.Exit(1)
	}
	exitCode := m.Run()
	if err := cleanup(); err != nil {
		fmt.Printf("error when deleting %s: %v\n", fileName, err)
		os.Exit(1)
	}
	os.Exit(exitCode)
}

func TestRun(t *testing.T) {
	totalLines := 10
	args := []string{
		"--llmin",
		"10000",
		"--llmax",
		"30000",
		"--ulmin",
		"100",
		"--ulmax",
		"3000",
		fmt.Sprintf("-t=%d", totalLines),
		"-p=0.7",
		fmt.Sprintf("-f=%s", fileName),
	}
	run(args)
	f, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("error when opening file %s: %v", fileName, err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	count := 0
	totalSuspicious := 0
	for scanner.Scan() {
		count++
		line := scanner.Text()
		transaction, err := transaction.New(line)
		if err != nil {
			t.Fatalf("error reading transaction from file %s: %v", fileName, err)
		}
		if transaction.IsSuspicious() {
			totalSuspicious++
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("reading file %s: %v", fileName, err)
	}
	require.Equal(t, totalLines, count)
	require.Equal(t, totalSuspicious, 3)
}
