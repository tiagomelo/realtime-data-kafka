package task

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

type worker struct {
	count *int32
}

func (w *worker) Incr() {
	atomic.AddInt32(w.count, 1)
}

func (w *worker) Work(ctx context.Context) {
	w.Incr()
}

func TestWorker(t *testing.T) {
	var count int32
	ctx := context.TODO()
	maxGoroutines := runtime.GOMAXPROCS(0)
	n := 10
	task := New(ctx, maxGoroutines)
	var wg sync.WaitGroup
	wg.Add(maxGoroutines * n)

	for i := 0; i < maxGoroutines; i++ {
		for j := 0; j < n; j++ {
			w := &worker{&count}
			go func() {
				task.Do(w)
				wg.Done()
			}()
		}
	}
	wg.Wait()
	require.GreaterOrEqual(t, int32(100), count)
}
