package context_test

import (
	"sync"
	"testing"
	"time"

	"github.com/posener/context"
	"github.com/stretchr/testify/assert"
)

const shortTime = 100 * time.Millisecond

func TestContextPropagation(t *testing.T) {
	t.Parallel()

	context.Init()

	assertNoDeadline(t)

	ctx, cancel := context.WithTimeout(context.Get(), shortTime)
	defer cancel()

	assertNoDeadline(t)

	func() {
		assertNoDeadline(t)
		func() { assertNoDeadline(t) }()
	}()

	context.RunCtx(ctx, func() {
		assertWithDeadline(t)
		func() { assertWithDeadline(t) }()
	})

	var wg sync.WaitGroup
	wg.Add(2)

	context.Go(func() {
		assertNoDeadline(t)
		func() {
			assertNoDeadline(t)
			wg.Done()
		}()
	})

	context.GoCtx(ctx, func() {
		assertWithDeadline(t)
		context.Go(func() {
			assertWithDeadline(t)
			wg.Done()
		})
	})

	wg.Wait()
}

func assertWithDeadline(t *testing.T) {
	t.Helper()
	if _, ok := context.Get().Deadline(); !ok {
		t.Error("deadline did not propagated")
	}
}

func assertNoDeadline(t *testing.T) {
	t.Helper()
	if _, ok := context.Get().Deadline(); ok {
		t.Error("no deadline was defined")
	}
}

func TestPanic(t *testing.T) {
	t.Parallel()
	context.Init()

	var wg sync.WaitGroup
	wg.Add(2)

	t.Run("Using context.Get inside non-context goroutine", func(t *testing.T) {
		go func() {
			assert.Panics(t, func() { context.Get() })
			wg.Done()
		}()
	})

	t.Run("Using context.Go inside non-context goroutine", func(t *testing.T) {
		go func() {
			assert.Panics(t, func() { context.Go(func() {}) })
			wg.Done()
		}()
	})

	wg.Wait()
}
