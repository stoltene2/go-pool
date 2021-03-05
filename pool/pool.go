package pool

import (
	"io"
	"sync"
	"errors"
)

type Pool struct {
	m         sync.Mutex
	resources chan io.Closer
	factory   func() (io.Closer, error)
	closed    bool
}

var ErrPoolClosed = errors.New("Pool has been closed")

// Create a new resource pool.
// Takes a factory function, fn
// Takes the size of the pool. Throws an error if it is less than 1
// Returns pointer to pool or an error
func New(fn func() (io.Closer, error), size uint) (*Pool, error) {
	return &Pool{
		resources: make(chan io.Closer, size),
		factory:   fn,
	}, nil
}

// Acquire a resource from the existing pool or create a new one
func (p *Pool) Acquire() (io.Closer, error) {
	return nil, nil
}

// Release a resource back into a pool or close it
func (p *Pool) Release(r io.Closer) {

}

// Close the pool down
func (p *Pool) Close() {

}
