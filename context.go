package context

import (
	stdctx "context"
	"sync"

	"github.com/posener/context/runtime"
)

type (
	Context    = stdctx.Context
	CancelFunc = stdctx.CancelFunc
)

var (
	WithCancel   = stdctx.WithCancel
	WithTimeout  = stdctx.WithTimeout
	WithDeadline = stdctx.WithDeadline

	Background = stdctx.Background
	TODO       = stdctx.TODO

	DeadlineExceeded = stdctx.DeadlineExceeded
	Canceled         = stdctx.Canceled
)

var (
	// storage is used instead of goroutine local storage to
	// store goroutine(ID) to Context mapping.
	storage map[uint64]Context
	// mutex for locking the storage map.
	mu sync.RWMutex
)

func init() {
	storage = make(map[uint64]Context)
}

// load simulates fetching of context from goroutine local storage
// It gets the context from `storage` map according to the current
// goroutine ID.
// If the goroutine ID is not in the map, it panic. This case
// may occur when a user did not use the `context.Go` to invoke a
// goroutine.
// Note: real goroutine local storage won't need the implemented locking
// exists in this implementation, since the storage won't be accessible from
// different goroutines.
func load() Context {
	id := runtime.GID()
	mu.RLock()
	defer mu.RUnlock()
	ctx := storage[id]
	if ctx == nil {
		panic("goroutine ran without using context.Go or context.GoCtx")
	}
	return ctx
}

// store simulates storing of context in the goroutine local storage.
// It gets the context to store, and returns an GID for later usage, if needed.
// Note: real goroutine local storage won't need the implemented locking
// exists in this implementation, since the storage won't be accessible from
// different goroutines.
func store(ctx Context) uint64 {
	id := runtime.GID()
	mu.Lock()
	defer mu.Unlock()
	storage[id] = ctx
	return id
}

// remove removes the context from the thread local storage according to a
// given GID.
// Note: real goroutine local storage won't need the implemented locking
// exists in this implementation, since the storage won't be accessible from
// different goroutines.
func remove(id uint64) {
	mu.Lock()
	defer mu.Unlock()
	delete(storage, id)
}

// Init creates the first background context in a program.
// it should be called once, in the beginning of the main
// function or in init() function.
// It returns the created context.
// All following goroutine invocations should be replaced
// by context.Go or context.GoCtx.
//
// Note:
// 		This function won't be needed in the real implementation.
func Init() Context {
	ctx := Background()
	store(ctx)
	return ctx
}

// Get gets the context of the current goroutine
// It may panic if the current go routine did not ran with
// context.Go or context.GoCtx.
//
// Note:
// 		This function won't panic in the real implementation.
func Get() Context {
	return load()
}

// Set updates the context of the current goroutine.
func Set(ctx Context) {
	store(ctx)
}

// Go invokes f in a new goroutine and takes care of propagating
// the current context to the created goroutine.
// It may panic if the current goroutine was not invoked with
// context.Go or context.GoCtx.
//
// Note:
// 		In the real implementation, this should be the behavior
// 		of the `go` keyword. It will also won't panic.
func Go(f func()) {
	GoCtx(load(), f)
}

// GoCtx invokes f in a new goroutine with the given context.
//
// Note:
// 		In the real implementation, accepting the context argument
//		should be incorporated into the behavior of the `go` keyword.
func GoCtx(ctx Context, f func()) {
	go func() {
		id := store(ctx)
		defer remove(id)
		f()
	}()
}
