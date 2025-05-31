package idGen

import "sync/atomic"

var counter int64

var NextID = func() int64 {
	return atomic.AddInt64(&counter, 1)
}
