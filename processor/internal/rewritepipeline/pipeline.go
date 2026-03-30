package rewritepipeline

import (
	"fmt"
	"io"

	"github.com/triviajon/liteproxy/processor/internal/logging"
)

// Pipeline holds a slice of rewriters and executes them in order.
type Pipeline struct {
	rewriters []Rewriter
}

// NewPipeline creates a new Pipeline with the given rewriters.
// Precondition: at least one rewriter must be provided, no rewriter can be nil.
// Returns: error if preconditions are violated.
func NewPipeline(rewriters ...Rewriter) (*Pipeline, error) {
	if len(rewriters) == 0 {
		return nil, fmt.Errorf("at least one rewriter must be provided")
	}
	for i, r := range rewriters {
		if r == nil {
			return nil, fmt.Errorf("rewriter at index %d must not be nil", i)
		}
	}
	logging.Infof("Initialized with %d rewriter(s)", len(rewriters))
	return &Pipeline{rewriters: rewriters}, nil
}

// Process processes the input through all rewriters in the pipeline.
// Requires that input is not nil and contentType is not empty.
// Returns an io.ReadCloser with the processed content, otherwise an error from a rewriter or describing a constraint violation.
func (p *Pipeline) Process(input io.Reader, contentType string) (io.ReadCloser, error) {
	if input == nil {
		return nil, fmt.Errorf("input must not be nil")
	}
	if contentType == "" {
		return nil, fmt.Errorf("contentType must not be empty")
	}

	logging.Debugf("Processing started - content_type=%s rewriter_count=%d", contentType, len(p.rewriters))

	if len(p.rewriters) == 0 {
		logging.Debugf("No rewriters, passing input through")
		return io.NopCloser(input), nil
	}

	var currentReader io.Reader = input
	var lastCloser io.ReadCloser

	for i, r := range p.rewriters {
		logging.Debugf("Executing rewriter %d/%d", i+1, len(p.rewriters))
		output, err := r.Rewrite(currentReader, contentType)
		if err != nil {
			logging.Errorf("Rewriter %d error - error=%v", i+1, err)
			return nil, err
		}
		logging.Debugf("Rewriter %d completed successfully", i+1)
		currentReader = output
		lastCloser = output
	}

	logging.Debugf("Processing completed - content_type=%s", contentType)
	return lastCloser, nil
}
