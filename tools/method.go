package tools

import (
	"context"
	"errors"
	"time"
)

type packedResult[T any] struct {
	Result T
	Err    error
}

func RunWithTimeout[T any](
	ctx context.Context, timeout time.Duration, fn func(context.Context) (T, error),
) (T, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	subCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resultChan := make(chan *packedResult[T], 1)
	go func() {
		defer close(resultChan)
		defer func() {
			if r := recover(); r != nil {
				resultChan <- &packedResult[T]{
					Result: *new(T), // Return zero value of T
					Err:    errors.New("panic occurred in RunWithTimeout"),
				}
			}
		}()

		result, err := fn(subCtx)
		resultChan <- &packedResult[T]{Result: result, Err: err}
	}()

	select {
	case <-subCtx.Done():
		return *new(T), subCtx.Err() // Return zero value of T on timeout
	case result := <-resultChan:
		return result.Result, result.Err
	}
}
