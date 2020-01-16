package types

// Semaphore presents a channel that implements a semaphore pattern.
type Semaphore chan struct{}

// NewSemaphore returns a new semaphore ready to use.
func NewSemaphore(size int) Semaphore {
	return make(Semaphore, size)
}

// Acquire gets the n resources from the Semaphore.
func (s Semaphore) Acquire(n int) {
	l := struct{}{}
	for i := 0; i < n; i++ {
		s <- l
	}
}

// Release returns the specified number of resources to the Semaphore.
func (s Semaphore) Release(n int) {
	if len(s) == 0 {
		return
	}

	for i := 0; i < n; i++ {
		<-s
	}
}
