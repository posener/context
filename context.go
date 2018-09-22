package context

import (
	. "context"
	"sync"
	"time"
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
func load() Context {
	id := goroutineID()
	ctx := storage[id]
	if ctx == nil {
		panic("goroutine ran without using context.Go or context.GoCtx")
	}
	return ctx
}

func safeLoad() Context {
	mu.RLock()
	defer mu.RUnlock()
	return load()
}

func store(ctx Context) uint64 {
	id := goroutineID()
	storage[id] = ctx
	return id
}

func safeRemove(id uint64) {
	mu.Lock()
	defer mu.Unlock()
	delete(storage, id)
}

func Init() {
	mu.Lock()
	defer mu.Unlock()
	store(Background())
}

func Get() Context {
	return safeLoad()
}

func GoCtx(ctx Context, f func()) {
	go func() {
		mu.Lock()
		id := store(ctx)
		mu.Unlock()
		defer safeRemove(id)
		f()
	}()
}

func Go(f func()) {
	GoCtx(safeLoad(), f)
}

func NewWithCancel() (ctx Context, cancel CancelFunc) {
	return WithCancel(load())

}

func NewWithDeadline(deadline time.Time) (Context, CancelFunc) {
	return WithDeadline(safeLoad(), deadline)
}

func NewWithTimeout(timeout time.Duration) (Context, CancelFunc) {
	return WithTimeout(safeLoad(), timeout)
}

func NewWithValue(key, val interface{}) Context {
	return WithValue(safeLoad(), key, val)
}

func Deadline() (deadline time.Time, ok bool) {
	return safeLoad().Deadline()
}

func Done() <-chan struct{} {
	return safeLoad().Done()
}

func Err() error {
	return safeLoad().Err()
}

func Value(key interface{}) interface{} {
	return safeLoad().Value(key)
}

func GetCancel() CancelFunc {
	mu.Lock()
	defer mu.Unlock()
	ctx, cancel := WithCancel(load())
	store(ctx)
	return cancel
}

func SetDeadline(deadline time.Time) CancelFunc {
	mu.Lock()
	defer mu.Unlock()
	ctx, cancel := WithDeadline(load(), deadline)
	store(ctx)
	return cancel
}

func SetTimeout(timeout time.Duration) CancelFunc {
	mu.Lock()
	defer mu.Unlock()
	ctx, cancel := WithTimeout(load(), timeout)
	store(ctx)
	return cancel
}

func SetValue(key, val interface{}) {
	mu.Lock()
	defer mu.Unlock()
	ctx := WithValue(load(), key, val)
	store(ctx)
}
