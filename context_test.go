package context

import (
	"sync"
	"testing"
	"time"
)

func TestSetTimeout(t *testing.T) {
	t.Parallel()
	Init()

	ctx := Get()
	cancel := SetTimeout(100 * time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	GoCtx(ctx, func() {
		if _, ok := Deadline(); ok {
			t.Error("no deadline was defined")
		}
		select {
		case <-Done():
			t.Error("no deadline was defined")
		case <-time.After(100 * time.Millisecond):
		}
		wg.Done()
	})

	Go(func() {
		if _, ok := Deadline(); !ok {
			t.Error("deadline did not propagated")
		}
		select {
		case <-Done():
		case <-time.After(500 * time.Millisecond):
			t.Error("deadline did not propagated")
		}
		wg.Done()
	})

	wg.Wait()
}

func TestNewWithTimeout(t *testing.T) {
	t.Parallel()
	Init()
	ctx, cancel := NewWithTimeout(100 * time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	Go(func() {
		if _, ok := Deadline(); ok {
			t.Error("no deadline was defined")
		}
		select {
		case <-Done():
			t.Error("no deadline was defined")
		case <-time.After(100 * time.Millisecond):
		}
		wg.Done()
	})

	GoCtx(ctx, func() {
		if _, ok := Deadline(); !ok {
			t.Error("deadline did not propagated")
		}
		select {
		case <-Done():
		case <-time.After(500 * time.Millisecond):
			t.Error("deadline did not propagated")
		}
		wg.Done()
	})

	wg.Wait()
}
