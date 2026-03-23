package util

import "fmt"

// Bind1 takes a binary function and a value for the first argument, returning a unary function that only requires the second argument.
// Requires that fn is not nil.
// Returns a unary function, otherwise an error describing the constraint violation.
func Bind1[A, B, C any](fn func(A, B) C, a A) (func(B) C, error) {
	if fn == nil {
		return nil, fmt.Errorf("fn must not be nil")
	}
	return func(b B) C {
		return fn(a, b)
	}, nil
}
