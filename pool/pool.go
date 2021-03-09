package pool

import (
	"errors"
	"io"
	"sync"
	"log"
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
	if size < 1 {
		return nil, errors.New("Pool size is too small")
	}

	return &Pool{
		resources: make(chan io.Closer, size),
		factory:   fn,
	}, nil
}

// Acquire a resource from the existing pool or create a new one If
// there is an item available in the resource pool, use that
// one. Otherwise create a new resource.
// If the resource pool is closed we'll get ErrPoolClosed
func (p *Pool) Acquire() (io.Closer, error) {
	select {
	case r, ok := <-p.resources:
		if !ok {
			return nil, ErrPoolClosed
		}

		log.Println("Acquire:", "Shared Resource")
		return r, nil
	default:
		log.Println("Acquire:", "New Resource")
		return p.factory()
	}
}

// Release a resource back into a pool or close it.  This function is
// protected by a mutex with closing.  If there are more slots
// available in the resource pool it is inserted back in the
// queue. Otherwise the resource is closed.
func (p *Pool) Release(r io.Closer) {
	p.m.Lock()
	defer p.m.Unlock()

	if p.closed {
		return
	}

	select {
	case p.resources <- r:
		log.Println("Release:", "In queue")
	default:
		log.Println("Release:", "Closing")
		r.Close()
	}
}

// Close the pool down.  Sets the closes the pool, set the state to
// closed inside the pool struct, and closes each resource.
func (p *Pool) Close() {
	p.m.Lock()
	defer p.m.Unlock()

	if p.closed {
		return
	}

	p.closed = true
	close(p.resources)

	for r := range p.resources {
        r.Close()
	}
}
