// Copyright (c) 2022 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
package task

import (
	"context"
	"sync"
)

// Worker must be implemented by types that want to use
// the run pool.
type Worker interface {
	Work(ctx context.Context)
}

// Task provides a pool of goroutines that can execute any Worker
// tasks that are submitted.
type Task struct {
	ctx  context.Context
	work chan Worker
	wg   sync.WaitGroup
}

// New creates a new work pool.
func New(ctx context.Context, maxGoroutines int) *Task {
	t := Task{

		// Using an unbuffered channel because we want the
		// guarantee of knowing the work being submitted is
		// actually being worked on after the call to Run returns.
		work: make(chan Worker),
		ctx:  ctx,
	}

	// The goroutines are the pool. So we could add code
	// to change the size of the pool later on.

	t.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func() {
			for w := range t.work {
				w.Work(ctx)
			}
			t.wg.Done()
		}()
	}

	return &t
}

// Shutdown waits for all the goroutines to shutdown.
func (t *Task) Shutdown() {
	close(t.work)
	t.wg.Wait()
}

// Do submits work to the pool.
func (t *Task) Do(w Worker) {
	t.work <- w
}
