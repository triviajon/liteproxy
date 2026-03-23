package util

// Bind1 takes a binary function and a value for the first argument, returning a unary function that only requires the second argument.
func Bind1[A, B, C any](fn func(A, B) C, a A) func(B) C {
	return func(b B) C {
		return fn(a, b)
	}
}
