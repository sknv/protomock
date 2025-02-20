package closer

import (
	"context"
	"fmt"
)

// PlainCloser describes a closer without using context.
type PlainCloser func() error

// CloseWithContext tries to apply the plain closer gracefully using a provided context.
func CloseWithContext(ctx context.Context, closer PlainCloser) error {
	errs := make(chan error, 1)

	go func() {
		err := closer()
		errs <- err
	}()

	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
		return fmt.Errorf("context done: %w", ctx.Err())
	}
}
