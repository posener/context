package runtime

import (
	"bytes"
	rt "runtime"
	"strconv"
	"sync"
)

const bufSize = 64

var bufPool = sync.Pool{New: func() interface{} { return make([]byte, bufSize) }}

// GID returns the current goroutine ID
func GID() uint64 {
	b := bufPool.Get().([]byte)
	defer bufPool.Put(b[:bufSize])
	b = b[:rt.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
