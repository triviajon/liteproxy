package rewritepipeline

import "io"

// Rewrite takes the original body and returns a modified version.
type Rewriter interface {
	// CanHandle returns true if this rewriter should process this mime type.
	CanHandle(contentType string) bool
	// Rewrite processes the stream.
	Rewrite(input io.Reader, contentType string) (io.ReadCloser, error)
}
