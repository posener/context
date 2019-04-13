package context

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testContext int

func (testContext) Deadline() (deadline time.Time, ok bool) { return }
func (testContext) Done() <-chan struct{}                   { return make(<-chan struct{}) }
func (testContext) Err() error                              { return nil }
func (testContext) Value(key interface{}) interface{}       { return nil }

func TestSet(t *testing.T) {
	t.Parallel()
	ctx1 := Init()
	ctx2 := testContext(2)
	ctx3 := testContext(3)

	unset := Set(ctx2)

	var wg sync.WaitGroup
	wg.Add(4)

	Go(func() {
		assert.Equal(t, Get(), ctx2)
		wg.Done()
	})

	GoCtx(ctx2, func() {
		assert.Equal(t, Get(), ctx2)
		wg.Done()
	})

	GoCtx(ctx3, func() {
		assert.Equal(t, Get(), ctx3)
		wg.Done()
	})

	unset()

	Go(func() {
		assert.Equal(t, Get(), ctx1)
		wg.Done()
	})

	wg.Wait()
}

func TestSetNested(t *testing.T) {
	t.Parallel()
	ctx1 := Init()
	ctx2 := testContext(2)
	ctx3 := testContext(3)

	assert.Equal(t, Get(), ctx1)
	unset2 := Set(ctx2)
	assert.Equal(t, Get(), ctx2)
	unset3 := Set(ctx3)
	assert.Equal(t, Get(), ctx3)
	unset3()
	assert.Equal(t, Get(), ctx2)
	unset2()
	assert.Equal(t, Get(), ctx1)
}

func TestFunctionScope(t *testing.T) {
	t.Parallel()
	ctx1 := Init()
	ctx2 := testContext(2)

	func() {
		assert.Equal(t, Get(), ctx1)
		defer Set(ctx2)()
		assert.Equal(t, Get(), ctx2)
	}()

	assert.Equal(t, Get(), ctx1)
}

func TestPanic(t *testing.T) {
	t.Parallel()
	Init()

	t.Run("Using context.Get inside non-context goroutine", func(t *testing.T) {
		assert.Panics(t, func() { Get() })
	})

	t.Run("Using context.Go inside non-context goroutine", func(t *testing.T) {
		assert.Panics(t, func() { Go(func() {}) })
	})

	t.Run("Invoking unset twice", func(t *testing.T) {
		unset := Set(testContext(1))
		unset()
		assert.Panics(t, unset)
	})

	t.Run("Invoking unset unordered", func(t *testing.T) {
		unset1 := Set(testContext(1))
		unset2 := Set(testContext(2))
		assert.Panics(t, unset1)
		unset2()
		unset1()
	})
}
