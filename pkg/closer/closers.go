package closer

import (
	"container/list"
	"context"
	"errors"
	"sync"
)

// Closer defines close signature.
type Closer func(context.Context) error

// Closers manages closers in correct order.
type Closers struct {
	list *list.List
	mu   sync.Mutex
}

// New returns a Closers instance.
func New() *Closers {
	return &Closers{
		list: list.New(),
		mu:   sync.Mutex{},
	}
}

// Add a closer to the list.
func (c *Closers) Add(closer Closer) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.list.PushFront(closer) // Close in reversed order, just like defer does.
}

// Close the closers in reversed order.
func (c *Closers) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var errs error

	for el := c.list.Front(); el != nil; el = el.Next() {
		closer, _ := el.Value.(Closer)
		if err := closer(ctx); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	c.list.Init() // Clear the list.

	return errs
}
