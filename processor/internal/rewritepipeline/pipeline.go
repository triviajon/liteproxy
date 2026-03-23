package rewritepipeline

import (
	"io"
)

// Pipeline holds a slice of rewriters and executes them in order.
type Pipeline struct {
	rewriters []Rewriter
}

func NewPipeline(rewriters ...Rewriter) *Pipeline {
	return &Pipeline{rewriters: rewriters}
}

func (p *Pipeline) Process(input io.Reader, contentType string) (io.ReadCloser, error) {
	if len(p.rewriters) == 0 {
		return io.NopCloser(input), nil
	}

	var currentReader io.Reader = input
	var lastCloser io.ReadCloser

	for _, r := range p.rewriters {
		output, err := r.Rewrite(currentReader, contentType)
		if err != nil {
			return nil, err
		}
		currentReader = output
		lastCloser = output
	}

	return lastCloser, nil
}
