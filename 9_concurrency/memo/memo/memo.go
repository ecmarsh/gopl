// Package memo is an example of a concurrent non-blocking cache
// it is concurrency-safe and avoids the contention associated
// with designs based on a single lock for the entire cache.
package memo

// Func is the type of the function to memoize
type Func func(key string) (interface{}, error)

// result of a memoized function call
type result struct {
	value interface{}
	err   error
}

// A request is a message requesting that the Func be applied to key.
type request struct {
	key      string
	response chan<- result // the client wants a single result
}

// Memo consists of a channel through which the caller of
// Get communicates with the monitor goroutine.
// The channel only carries a single value.
type Memo struct {
	requests chan request
}

// New returns a memoization of f. Clients must subsequently call Close.
func New(f Func) *Memo {
	m := &Memo{requests: make(chan request)}
	go m.server(f)
	return m
}

// Get communicates with the monitor goroutine.
// The caller of Get sends the monitor goroutine both the
// key (the argument to the memoized function), and another
// channel response, over which the result should be sent back
// when it becomes available.
func (m *Memo) Get(key string) (interface{}, error) {
	response := make(chan result)
	m.requests <- request{key, response}
	res := <-response
	return res.value, res.err
}

// Close closes the memo's requests channel.
func (m *Memo) Close() {
	close(m.requests)
}

// The cache variable is defined to server now.
// The monitor reads requests in a loop until the request
// channel is closed by the Close method.
// For each request, it consults the cache, creating and
// inserting a new entry if none was found.
func (m *Memo) server(f Func) {
	cache := make(map[string]*entry)
	for req := range m.requests {
		e := cache[req.key]
		if e == nil {
			// This is the first request for this key.
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key) // call f(key)
		}
		go e.deliver(req.response)
	}
}

// Each entry contains the memoized result of a call to the
// function f and contains a channel called ready. Just
// after the entry's result has been set, this channel will
// be closed, to broadcast to any other gourtines that it is
// now safe for them to read the result from the entry.
type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

// NOTE: the call and deliver methods must be called in their own goroutines
// to ensure that the monitor goroutine does not stop processing new reqs.

// The first request for a given key becomes responsible
// for calling the function f on that key, storing the
// result in the entry, and broadcasting the readiness of the entry
// by closing the ready channel.
func (e *entry) call(f Func, key string) {
	// Evaluate the function.
	e.res.value, e.res.err = f(key)
	// Broadcast the ready condition.
	close(e.ready)
}

// A subsequent request for the same key finds the existing
// entry in the map, waits for the result to become ready, and sends
// the result through the response channel to the client goroutine
// that called Get.
func (e *entry) deliver(response chan<- result) {
	// Wait for the ready condition.
	<-e.ready
	// Send the result to the client.
	response <- e.res
}
